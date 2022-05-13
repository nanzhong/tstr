package server

import (
	"context"
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/api/runner/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RunnerServer struct {
	pgxPool   *pgxpool.Pool
	dbQuerier db.Querier

	runner.UnimplementedRunnerServiceServer
}

func NewRunnerServer(pgxPool *pgxpool.Pool) runner.RunnerServiceServer {
	return &RunnerServer{
		pgxPool:   pgxPool,
		dbQuerier: db.New(),
	}
}

func (s *RunnerServer) RegisterRunner(ctx context.Context, req *runner.RegisterRunnerRequest) (*runner.RegisterRunnerResponse, error) {
	for _, v := range req.AcceptTestLabelSelectors {
		if _, err := regexp.Compile(v); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid accept test label selectors, must be valid RE")
		}
	}

	for _, v := range req.RejectTestLabelSelectors {
		if _, err := regexp.Compile(v); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid reject test label selectors, must be valid RE")
		}
	}

	var (
		accept pgtype.JSONB
		reject pgtype.JSONB
	)

	if err := accept.Set(req.AcceptTestLabelSelectors); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid accept test label selectors")
	}
	if err := reject.Set(req.RejectTestLabelSelectors); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reject test label selectors")
	}

	regRunner, err := s.dbQuerier.RegisterRunner(ctx, s.pgxPool, db.RegisterRunnerParams{
		Name:                     req.Name,
		AcceptTestLabelSelectors: accept,
		RejectTestLabelSelectors: reject,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to register runner in db")
		return nil, status.Error(codes.Internal, "failed to register runner")
	}

	var (
		acceptSelectors map[string]string
		rejectSelectors map[string]string
	)

	if err := regRunner.AcceptTestLabelSelectors.AssignTo(&acceptSelectors); err != nil {
		log.Error().Err(err).Msg("failed to format accept test label selectors")
		return nil, status.Error(codes.Internal, "failed to format response")
	}
	if err := regRunner.RejectTestLabelSelectors.AssignTo(&rejectSelectors); err != nil {
		log.Error().Err(err).Msg("failed to format reject test label selectors")
		return nil, status.Error(codes.Internal, "failed to format response")
	}

	return &runner.RegisterRunnerResponse{
		Runner: &common.Runner{
			Id:                       regRunner.ID.String(),
			Name:                     regRunner.Name,
			AcceptTestLabelSelectors: acceptSelectors,
			RejectTestLabelSelectors: rejectSelectors,
			RegisteredAt:             types.ToProtoTimestamp(regRunner.RegisteredAt),
			LastHeartbeatAt:          types.ToProtoTimestamp(regRunner.LastHeartbeatAt),
		},
	}, nil
}

func (s *RunnerServer) NextRun(ctx context.Context, req *runner.NextRunRequest) (*runner.NextRunResponse, error) {
	runnerID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid runner id")
	}

	dbRunner, err := s.dbQuerier.GetRunner(ctx, s.pgxPool, runnerID)
	if err != nil {
		log.Error().
			Err(err).
			Str("runner_id", req.Id).
			Msg("failed to get runner")
		return nil, status.Error(codes.Internal, "failed to find runner info")
	}

	var (
		acceptSelectors   map[string]string
		acceptSelectorsRE = make(map[string]*regexp.Regexp)
		rejectSelectors   map[string]string
		rejectSelectorsRE = make(map[string]*regexp.Regexp)

		acceptKeys []string
	)
	if err := dbRunner.AcceptTestLabelSelectors.AssignTo(&acceptSelectors); err != nil {
		log.Error().
			Err(err).
			Msg("failed to parse accept test label selectors")
		return nil, status.Error(codes.Internal, "failed to load runner info")
	}
	if err := dbRunner.RejectTestLabelSelectors.AssignTo(&rejectSelectors); err != nil {
		log.Error().
			Err(err).
			Msg("failed to parse reject test label selectors")
		return nil, status.Error(codes.Internal, "failed to load runner info")
	}

	for k, v := range acceptSelectors {
		acceptKeys = append(acceptKeys, k)
		re, err := regexp.Compile(v)
		if err != nil {
			log.Error().
				Err(err).
				Str("selector", v).
				Msg("failed to compile label selector RE")
			return nil, status.Error(codes.Internal, "failed to load runner info")
		}
		acceptSelectorsRE[k] = re
	}

	for k, v := range rejectSelectors {
		re, err := regexp.Compile(v)
		if err != nil {
			log.Error().
				Err(err).
				Str("selector", v).
				Msg("failed to compile label selector RE")
			return nil, status.Error(codes.Internal, "failed to load runner info")
		}
		rejectSelectorsRE[k] = re
	}

	// NOTE we don't care about reject keys here, because unless all reject labels have selectors that match anything (non-trivial to determine), we need to first get the tests that match the accept keys before applying filtering.
	tests, err := s.dbQuerier.ListTestsIDsMatchingLabelKeys(ctx, s.pgxPool, db.ListTestsIDsMatchingLabelKeysParams{
		IncludeLabelKeys: acceptKeys,
		FilterLabelKeys:  nil,
	})
	if err != nil {
		log.Error().
			Err(err).
			Strs("accept_keys", acceptKeys).
			Msg("failed to find tests matching label keys")
		return nil, status.Error(codes.Internal, "failed to determine matching tests for runner")
	}

	var matchingTestIDs []uuid.UUID
	for _, test := range tests {
		var labels map[string]string
		if err := test.Labels.AssignTo(&labels); err != nil {
			log.Error().
				Err(err).
				Str("test_id", test.ID.String()).
				Msg("failed to parse test labels")
			return nil, status.Error(codes.Internal, "failed to determine matching tests for runner")
		}

		matches := true
		for k, v := range labels {
			acceptRE, ok := acceptSelectorsRE[k]
			if !ok {
				continue
			}
			if !acceptRE.Match([]byte(v)) {
				matches = false
				break
			}

			rejectRE, ok := rejectSelectorsRE[k]
			if !ok {
				continue
			}
			if rejectRE.Match([]byte(v)) {
				matches = false
				break
			}
		}
		if !matches {
			continue
		}

		matchingTestIDs = append(matchingTestIDs, test.ID)
	}

	run, err := s.dbQuerier.AssignRun(ctx, s.pgxPool, db.AssignRunParams{
		RunnerID: runnerID,
		TestIds:  matchingTestIDs,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "no matching runs")
		}

		// NOTE []uuid.UUID can't be directly used as []fmt.Stringer
		matchingTestIDStrings := make([]string, len(matchingTestIDs))
		for i, s := range matchingTestIDs {
			matchingTestIDStrings[i] = s.String()
		}
		log.Error().
			Err(err).
			Stringer("runner_id", runnerID).
			Strs("test_ids", matchingTestIDStrings).
			Msg("failed to assign run to runner")
		return nil, status.Error(codes.Internal, "failed to assign run to runner")
	}

	var env map[string]string
	if err := run.Env.AssignTo(&env); err != nil {
		log.Error().
			Err(err).
			Stringer("run_id", run.ID).
			Msg("failed to parse run config env")
		return nil, status.Error(codes.Internal, "failed to format response")
	}
	return &runner.NextRunResponse{
		Run: &common.Run{
			Id:     run.ID.String(),
			TestId: run.TestID.String(),
			TestRunConfig: &common.Test_RunConfig{
				Id:             run.TestRunConfigID.String(),
				ContainerImage: run.ContainerImage,
				Command:        run.Command.String,
				Args:           run.Args,
				Env:            env,
				CreatedAt:      types.ToProtoTimestamp(run.TestRunConfigCreatedAt),
			},
			RunnerId:    run.RunnerID.UUID.String(),
			ScheduledAt: types.ToProtoTimestamp(run.ScheduledAt),
		},
	}, nil
}

func (s *RunnerServer) SubmitRun(server runner.RunnerService_SubmitRunServer) error {
	return nil
}
