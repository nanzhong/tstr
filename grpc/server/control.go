package server

import (
	"context"
	"database/sql"

	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	controlv1 "github.com/nanzhong/tstr/api/control/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ControlServer struct {
	controlv1.UnimplementedControlServiceServer
	pgxPool    *pgxpool.Pool
	dbQuerier  db.Querier
	cronParser cron.Parser
	clock      clock.Clock
}

func NewControlServer(pgxPool *pgxpool.Pool) controlv1.ControlServiceServer {
	return &ControlServer{
		pgxPool:    pgxPool,
		dbQuerier:  db.New(),
		cronParser: cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
		clock:      clock.New(),
	}
}

func (s *ControlServer) RegisterTest(ctx context.Context, r *controlv1.RegisterTestRequest) (*controlv1.RegisterTestResponse, error) {
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

	var nextRunAt sql.NullTime
	if r.CronSchedule != "" {
		schedule, err := s.cronParser.Parse(r.CronSchedule)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid cron schedule")
		}
		nextRunAt.Valid = true
		nextRunAt.Time = schedule.Next(s.clock.Now())
	}

	result, err := s.dbQuerier.RegisterTest(ctx, s.pgxPool, db.RegisterTestParams{
		Name:           r.Name,
		Labels:         labels,
		CronSchedule:   r.CronSchedule,
		NextRunAt:      nextRunAt,
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

	return &controlv1.RegisterTestResponse{
		Test: &commonv1.Test{
			Id:           result.ID.String(),
			Name:         result.Name,
			Labels:       resultLabels,
			CronSchedule: result.CronSchedule.String,
			NextRunAt:    types.ToProtoTimestamp(result.NextRunAt),
			RunConfig: &commonv1.Test_RunConfig{
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

func (s *ControlServer) ScheduleRun(ctx context.Context, req *controlv1.ScheduleRunRequest) (*controlv1.ScheduleRunResponse, error) {
	testID, err := uuid.Parse(req.TestId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid test id")
	}

	test, err := s.dbQuerier.GetTest(ctx, s.pgxPool, testID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to lookup test to schedule run for")
	}

	run, err := s.dbQuerier.ScheduleRun(ctx, s.pgxPool, db.ScheduleRunParams{
		TestID:          test.ID,
		TestRunConfigID: test.TestRunConfigID.UUID,
	})
	if err != nil {
		log.Error().
			Err(err).
			Stringer("test_id", test.ID).
			Stringer("test_urn_config_id", test.TestRunConfigID.UUID).
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

	return &controlv1.ScheduleRunResponse{
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
			ScheduledAt: types.ToProtoTimestamp(run.ScheduledAt),
		},
	}, nil
}
