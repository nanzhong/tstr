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

	runConfig := pgtype.JSONB{}
	if err := runConfig.Set(types.FromProtoTestRunConfig(r.RunConfig)); err != nil {
		log.Error().Err(err).Msg("failed to parse run config")
		return nil, status.Error(codes.InvalidArgument, "failed to parse run config")
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
		Name:         r.Name,
		Labels:       labels,
		RunConfig:    runConfig,
		CronSchedule: r.CronSchedule,
		NextRunAt:    nextRunAt,
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
	var resultRunConfig db.TestRunConfig
	if err := result.RunConfig.AssignTo(&resultRunConfig); err != nil {
		log.Error().Err(err).Msg("failed to format stored run config")
		return nil, status.Error(codes.Internal, "failed to format run config")
	}

	return &controlv1.RegisterTestResponse{
		Test: &commonv1.Test{
			Id:           result.ID.String(),
			Name:         result.Name,
			Labels:       resultLabels,
			CronSchedule: result.CronSchedule.String,
			NextRunAt:    types.ToProtoTimestamp(result.NextRunAt),
			RunConfig:    types.ToProtoTestRunConfig(resultRunConfig),
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

	run, err := s.dbQuerier.ScheduleRun(ctx, s.pgxPool, test.ID)
	if err != nil {
		log.Error().
			Err(err).
			Stringer("test_id", test.ID).
			Msg("failed to schedule run")
		return nil, status.Error(codes.Internal, "failed to schedule run")
	}

	var runConfig db.TestRunConfig
	if err := run.TestRunConfig.AssignTo(&runConfig); err != nil {
		log.Error().
			Err(err).
			Stringer("run_id", run.ID).
			Msg("failed to parse run config")
		return nil, status.Error(codes.Internal, "failed to format run config")
	}

	return &controlv1.ScheduleRunResponse{
		Run: &commonv1.Run{
			Id:            run.ID.String(),
			TestId:        run.TestID.String(),
			TestRunConfig: types.ToProtoTestRunConfig(runConfig),
			ScheduledAt:   types.ToProtoTimestamp(run.ScheduledAt),
		},
	}, nil
}
