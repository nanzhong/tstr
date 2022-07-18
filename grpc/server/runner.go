package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	runnerv1 "github.com/nanzhong/tstr/api/runner/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RunnerServer struct {
	pgxPool   *pgxpool.Pool
	dbQuerier db.Querier
	clock     clock.Clock

	runnerv1.UnimplementedRunnerServiceServer
}

func NewRunnerServer(pgxPool *pgxpool.Pool) runnerv1.RunnerServiceServer {
	return &RunnerServer{
		pgxPool:   pgxPool,
		dbQuerier: db.New(),
		clock:     clock.New(),
	}
}

func (s *RunnerServer) RegisterRunner(ctx context.Context, req *runnerv1.RegisterRunnerRequest) (*runnerv1.RegisterRunnerResponse, error) {
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

	return &runnerv1.RegisterRunnerResponse{
		Runner: &commonv1.Runner{
			Id:                       regRunner.ID.String(),
			Name:                     regRunner.Name,
			AcceptTestLabelSelectors: acceptSelectors,
			RejectTestLabelSelectors: rejectSelectors,
			RegisteredAt:             types.ToProtoTimestamp(regRunner.RegisteredAt),
			LastHeartbeatAt:          types.ToProtoTimestamp(regRunner.LastHeartbeatAt),
		},
	}, nil
}

func (s *RunnerServer) NextRun(ctx context.Context, req *runnerv1.NextRunRequest) (*runnerv1.NextRunResponse, error) {
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

	err = s.dbQuerier.UpdateRunnerHeartbeat(ctx, s.pgxPool, runnerID)
	if err != nil {
		log.Error().
			Err(err).
			Str("runner_id", req.Id).
			Msg("failed to update runner last_heartbeat")
		return nil, status.Error(codes.Internal, "failed to update runner heartbeat")
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
	return &runnerv1.NextRunResponse{
		Run: &commonv1.Run{
			Id:     run.ID.String(),
			TestId: run.TestID.String(),
			TestRunConfig: &commonv1.Test_RunConfig{
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

func (s *RunnerServer) SubmitRun(stream runnerv1.RunnerService_SubmitRunServer) error {
	var (
		runnerID   uuid.UUID
		runID      uuid.UUID
		result     db.RunResult
		startedAt  sql.NullTime
		finishedAt sql.NullTime
	)

	defer func() {
		if runnerID == uuid.Nil || runID == uuid.Nil {
			return
		}

		now := s.clock.Now()

		// This means the test run was not even successfully started.
		if !startedAt.Valid {
			if err := s.dbQuerier.UpdateRun(context.Background(), s.pgxPool, db.UpdateRunParams{
				ID:         runID,
				Result:     db.NullRunResult{Valid: true, RunResult: db.RunResultError},
				StartedAt:  sql.NullTime{},
				FinishedAt: sql.NullTime{},
			}); err != nil {
				log.Error().
					Err(err).
					Stringer("run_id", runID).
					Msg("failed to mark unstarted test run as error")
			}

			var pgLogs pgtype.JSONB
			if err := pgLogs.Set(db.RunLog{
				Type: commonv1.Run_Log_TSTR.String(),
				Time: now.Format(time.RFC3339Nano),
				Data: []byte("failed to start test run"),
			}); err != nil {
				log.Error().
					Err(err).
					Stringer("run_id", runID).
					Msg("failed to format start failure log for run")
				return
			}
			if err := s.dbQuerier.AppendLogsToRun(context.Background(), s.pgxPool, db.AppendLogsToRunParams{
				Logs: pgLogs,
				ID:   runID,
			}); err != nil {
				log.Error().
					Err(err).
					Stringer("run_id", runID).
					Msg("failed to append start failure log to run")
			}
			return
		}

		// This means the test run was not able to submit completion state.
		if !finishedAt.Valid {
			if err := s.dbQuerier.UpdateRun(context.Background(), s.pgxPool, db.UpdateRunParams{
				ID:         runID,
				Result:     db.NullRunResult{Valid: true, RunResult: db.RunResultError},
				StartedAt:  startedAt,
				FinishedAt: sql.NullTime{Valid: true, Time: now},
			}); err != nil {
				log.Error().
					Err(err).
					Stringer("run_id", runID).
					Msg("failed to mark unfinished test run as error")
			}

			var pgLogs pgtype.JSONB
			if err := pgLogs.Set(db.RunLog{
				Type: commonv1.Run_Log_TSTR.String(),
				Time: now.Format(time.RFC3339Nano),
				Data: []byte("failed to finish test run"),
			}); err != nil {
				log.Error().
					Err(err).
					Stringer("run_id", runID).
					Msg("failed to format finish failure log for run")
				return
			}
			if err := s.dbQuerier.AppendLogsToRun(context.Background(), s.pgxPool, db.AppendLogsToRunParams{
				Logs: pgLogs,
				ID:   runID,
			}); err != nil {
				log.Error().
					Err(err).
					Stringer("run_id", runID).
					Msg("failed to append finish failure log to run")
			}
		}
	}()

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&runnerv1.SubmitRunResponse{})
			}
			return err
		}

		if req.Id != "" {
			runnerID, err = uuid.Parse(req.Id)
			if err != nil {
				return status.Error(codes.InvalidArgument, "invalid id")
			}
		}

		if req.RunId != "" {
			runID, err = uuid.Parse(req.RunId)
			if err != nil {
				return status.Error(codes.InvalidArgument, "invalid run id")
			}
		}

		result = types.FromRunResult(req.Result)

		if req.StartedAt != nil {
			startedAt.Valid = true
			startedAt.Time = req.StartedAt.AsTime()
		}

		if req.FinishedAt != nil {
			finishedAt.Valid = true
			finishedAt.Time = req.FinishedAt.AsTime()
		}

		// We expect runner and run ids to be populated after the first received message
		if runnerID == uuid.Nil {
			return status.Error(codes.InvalidArgument, "missing id")
		}

		if runID == uuid.Nil {
			return status.Error(codes.InvalidArgument, "missing run id")
		}

		log.Debug().
			Stringer("run_id", runID).
			Stringer("runner_id", runnerID).
			Str("result", string(result)).
			Str("started_at", startedAt.Time.String()).
			Str("finished_at", finishedAt.Time.String()).
			Msg("received request")

		if err := s.dbQuerier.UpdateRun(stream.Context(), s.pgxPool, db.UpdateRunParams{
			ID:         runID,
			Result:     db.NullRunResult{Valid: true, RunResult: result},
			StartedAt:  startedAt,
			FinishedAt: finishedAt,
		}); err != nil {
			log.Error().Err(err).Msg("failed to save updated run")
			return status.Error(codes.Internal, "failed to update run")
		}

		if len(req.Logs) > 0 {
			var logs []db.RunLog
			for _, l := range req.Logs {
				logs = append(logs, db.RunLog{
					Type: l.OutputType.String(),
					Time: l.Time,
					Data: l.Data,
				})
			}
			fmt.Printf("%#v\n", logs)

			var pgLogs pgtype.JSONB
			if err := pgLogs.Set(logs); err != nil {
				log.Error().Err(err).Msg("failed to format run logs")
				return status.Error(codes.Internal, "failed to format run logs")
			}
			if err := s.dbQuerier.AppendLogsToRun(stream.Context(), s.pgxPool, db.AppendLogsToRunParams{
				ID:   runID,
				Logs: pgLogs,
			}); err != nil {
				log.Error().Err(err).Msg("failed to save run logs")
				return status.Error(codes.Internal, "failed to save run logs")
			}
		}
	}
}
