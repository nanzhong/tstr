package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/api/control/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ControlServer struct {
	control.UnimplementedControlServiceServer

	dbQuerier db.Querier
}

func NewControlServer(dbQuerier db.Querier) control.ControlServiceServer {
	return &ControlServer{
		dbQuerier: dbQuerier,
	}
}

func (s *ControlServer) RegisterTest(ctx context.Context, r *control.RegisterTestRequest) (*control.RegisterTestResponse, error) {
	labels := pgtype.JSONB{}
	if err := labels.Set(r.Labels); err != nil {
		log.Error().Err(err).Msg("failed to parse test labels")
		return nil, status.Error(codes.InvalidArgument, "failed to parse labels")
	}

	env := pgtype.JSONB{}
	if err := env.Set(r.RunConfig.Env); err != nil {
		log.Error().Err(err).Msg("failed to parse env")
		return nil, status.Error(codes.InvalidArgument, "failed to parse env")
	}

	result, err := s.dbQuerier.RegisterTest(ctx, db.RegisterTestParams{
		Name:           r.Name,
		Labels:         labels,
		CronSchedule:   r.CronSchedule,
		ContainerImage: r.RunConfig.ContainerImage,
		Command:        r.RunConfig.Command,
		Args:           r.RunConfig.Args,
		Env:            env,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to store test")
		return nil, status.Error(codes.Internal, "failed to register test")
	}

	var resultLabels map[string]string
	if err := result.Labels.AssignTo(&resultLabels); err != nil {
		log.Error().Err(err).Msg("failed to format stored labels")
		return nil, status.Error(codes.Internal, "failed to format labels")
	}
	var resultEnv map[string]string
	if err := result.Env.AssignTo(&resultEnv); err != nil {
		log.Error().Err(err).Msg("failed to format stored env")
		return nil, status.Error(codes.Internal, "failed to format env")
	}

	return &control.RegisterTestResponse{
		Test: &common.Test{
			Id:           result.ID.String(),
			Name:         result.Name,
			Labels:       resultLabels,
			CronSchedule: result.CronSchedule.String,
			RunConfig: &common.Test_RunConfig{
				Id:             result.TestRunConfigID.String(),
				ContainerImage: result.ContainerImage,
				Command:        result.Command.String,
				Args:           result.Args,
				Env:            resultEnv,
				CreatedAt:      types.ToProtoTimestamp(result.TestRunConfigCreatedAt),
			},
			RegisteredAt: types.ToProtoTimestamp(result.RegisteredAt),
			UpdatedAt:    types.ToProtoTimestamp(result.UpdatedAt),
		},
	}, nil
}

func (s *ControlServer) ScheduleRun(ctx context.Context, req *control.ScheduleRunRequest) (*control.ScheduleRunResponse, error) {
	testID, err := uuid.Parse(req.TestId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid test id")
	}

	test, err := s.dbQuerier.GetTest(ctx, testID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to lookup test to schedule run for")
	}

	run, err := s.dbQuerier.ScheduleRun(ctx, db.ScheduleRunParams{
		TestID:          test.ID,
		TestRunConfigID: test.TestRunConfigID,
	})
	if err != nil {
		log.Error().
			Err(err).
			Stringer("test_id", test.ID).
			Stringer("test_urn_config_id", test.TestRunConfigID).
			Msg("failed to schedule run")
		return nil, status.Error(codes.Internal, "failed to schedule run")
	}

	var env map[string]string
	if err := run.Env.AssignTo(&env); err != nil {
		log.Error().
			Err(err).
			Stringer("run_id", run.ID).
			Msg("failed to parse env")
		return nil, status.Error(codes.Internal, "failed to format run config env")
	}

	return &control.ScheduleRunResponse{
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
			ScheduledAt: types.ToProtoTimestamp(run.ScheduledAt),
		},
	}, nil
}
