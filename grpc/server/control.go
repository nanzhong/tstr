package server

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
	"github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/api/control/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	var (
		dbRunConfigCreatedAt time.Time
		dbRegisteredAt       time.Time
		dbUpdatedAt          time.Time
	)
	if err := result.TestRunConfigCreatedAt.AssignTo(&dbRunConfigCreatedAt); err != nil {
		log.Error().Err(err).Msg("failed to convert run config created at")
		return nil, status.Error(codes.Internal, "failed to format run config created at time")
	}
	if err := result.RegisteredAt.AssignTo(&dbRegisteredAt); err != nil {
		log.Error().Err(err).Msg("failed to convert registered at")
		return nil, status.Error(codes.Internal, "failed to format registered at time")
	}
	if err := result.UpdatedAt.AssignTo(&dbUpdatedAt); err != nil {
		log.Error().Err(err).Msg("failed to convert updated at")
		return nil, status.Error(codes.Internal, "failed to format updated at time")
	}

	return &control.RegisterTestResponse{
		Test: &common.Test{
			Id:           result.ID,
			Name:         result.Name,
			Labels:       resultLabels,
			CronSchedule: result.CronSchedule,
			RunConfig: &common.Test_RunConfig{
				Id:             result.TestRunConfigID,
				ContainerImage: result.ContainerImage,
				Command:        result.Command,
				Args:           result.Args,
				Env:            resultEnv,
				CreatedAt:      timestamppb.New(dbRunConfigCreatedAt),
			},
			RegisteredAt: timestamppb.New(dbRegisteredAt),
			UpdatedAt:    timestamppb.New(dbUpdatedAt),
		},
	}, nil
}

func (s *ControlServer) ScheduleRun(ctx context.Context, req *control.ScheduleRunRequest) (*control.ScheduleRunResponse, error) {
	test, err := s.dbQuerier.GetTest(ctx, req.TestId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to lookup test to schedule run for")
	}
	fmt.Printf("%#v\n", test)
	run, err := s.dbQuerier.ScheduleRun(ctx, test.ID, test.TestRunConfigID)
	if err != nil {
		log.Error().
			Err(err).
			Str("test_id", test.ID).
			Str("test_urn_config_id", test.TestRunConfigID).
			Msg("failed to schedule run")
		return nil, status.Error(codes.Internal, "failed to schedule run")
	}

	var env map[string]string
	if err := run.Env.AssignTo(&env); err != nil {
		log.Error().
			Err(err).
			Str("run_id", run.ID).
			Msg("failed to parse env")
		return nil, status.Error(codes.Internal, "failed to format run config env")
	}

	return &control.ScheduleRunResponse{
		Run: &common.Run{
			Id:     run.ID,
			TestId: run.TestID,
			TestRunConfig: &common.Test_RunConfig{
				Id:             run.TestRunConfigID,
				ContainerImage: run.ContainerImage,
				Command:        run.Command,
				Args:           run.Args,
				Env:            env,
				CreatedAt:      toProtoTimestamp(run.TestRunConfigCreatedAt),
			},
			ScheduledAt: toProtoTimestamp(run.ScheduledAt),
		},
	}, nil
}
