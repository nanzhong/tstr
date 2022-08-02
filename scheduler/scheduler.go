package scheduler

import (
	"context"
	"database/sql"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

const defaultRunTimeout = 5 * time.Minute

type Scheduler struct {
	pgxPool      *pgxpool.Pool
	dbQuerier    db.Querier
	cronParser   cron.Parser
	clock        clock.Clock
	nextSchedule time.Time
	retryDelay   time.Duration

	doneCh       chan struct{}
	stopCh       chan struct{}
	stopCancelFn context.CancelFunc
}

func New(pgxPool *pgxpool.Pool) *Scheduler {
	return &Scheduler{
		pgxPool:    pgxPool,
		dbQuerier:  db.New(),
		cronParser: cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor),
		clock:      clock.New(),
		doneCh:     make(chan struct{}),
		stopCh:     make(chan struct{}),
	}
}

func (s *Scheduler) Start() error {
	defer close(s.doneCh)

	var ctx context.Context
	ctx, s.stopCancelFn = context.WithCancel(context.Background())

	s.nextSchedule = s.clock.Now()
	for {
		select {
		case <-s.stopCh:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-s.clock.After(s.clock.Until(s.nextSchedule)):
			log.Debug().Msg("starting schedule pass")

			// First, reset orphaned runs.
			orphanThreshold := s.clock.Now().Add(-5 * time.Second)
			if err := s.dbQuerier.ResetOrphanedRuns(ctx, s.pgxPool, orphanThreshold); err != nil {
				log.Error().Err(err).Msg("failed to reset orphaned runs")
			}

			// Next, mark timeout tests as having errored.
			var timeoutLog pgtype.JSONB
			if err := timeoutLog.Set(db.RunLog{
				Type: commonv1.Run_Log_TSTR.String(),
				Time: s.clock.Now().Format(time.RFC3339Nano),
				Data: []byte("run timed out"),
			}); err != nil {
				log.Error().
					Err(err).
					Msg("failed to format timeout error log for run")
			} else {
				if err := s.dbQuerier.TimeoutRuns(ctx, s.pgxPool, db.TimeoutRunsParams{
					TimeoutLog:     timeoutLog,
					DefaultTimeout: int32(defaultRunTimeout.Seconds()),
				}); err != nil {
					log.Error().Err(err).Msg("failed to time out runs")
				}
			}

			// Next, schedule new runs.
			minNextRunAt := s.clock.Now().Add(time.Minute)
			err := s.pgxPool.BeginFunc(ctx, func(tx pgx.Tx) error {
				tests, err := s.dbQuerier.ListTestsToSchedule(ctx, tx)
				if err != nil {
					log.Error().Err(err).Msg("failed to list tests to schedule")
					return err
				}

				for _, test := range tests {
					runParams, err := RunsForTest(test)
					if err != nil {
						log.Error().Err(err).Stringer("test_id", test.ID).Msg("generating runs to schedule")
						return err
					}
					for _, runParam := range runParams {
						run, err := s.dbQuerier.ScheduleRun(ctx, tx, runParam)
						if err != nil {
							log.Error().
								Err(err).
								Stringer("test_id", test.ID).
								Stringer("test_matrix_id", runParam.TestMatrixID.UUID).
								Msg("failed to schedule run for test")
							return err
						}
						log.Info().
							Stringer("test_id", runParam.TestID).
							Stringer("run_id", run.ID).
							Stringer("test_matrix_id", runParam.TestMatrixID.UUID).
							Msg("scheduled run for test")
					}

					schedule, err := s.cronParser.Parse(test.CronSchedule.String)
					if err != nil {
						log.Error().
							Err(err).
							Stringer("test_id", test.ID).
							Str("test_cron_schedule", test.CronSchedule.String).
							Msg("failed to parse test cron schedule to calculate next run at")
						return err
					}

					nextRunAt := schedule.Next(s.clock.Now())
					if nextRunAt.Before(minNextRunAt) {
						minNextRunAt = nextRunAt
					}

					err = s.dbQuerier.UpdateTest(ctx, tx, db.UpdateTestParams{
						ID:           test.ID,
						Name:         test.Name,
						Labels:       test.Labels,
						Matrix:       test.Matrix,
						RunConfig:    test.RunConfig,
						CronSchedule: test.CronSchedule,
						NextRunAt:    sql.NullTime{Valid: true, Time: nextRunAt},
					})
					if err != nil {
						log.Error().
							Err(err).
							Stringer("test_id", test.ID).
							Msg("failed to update test with next run at after scheduling")
						return err
					}
				}
				return nil
			})
			if err != nil {
				log.Error().
					Err(err).
					Msg("failed to schedule runs, retrying next pass")
				s.nextSchedule = s.clock.Now().Add(s.retryDelay)
			} else {
				log.Debug().Msg("finished schedule pass")
				s.nextSchedule = minNextRunAt
			}
		}
	}
}

func (s *Scheduler) Stop(ctx context.Context) {
	close(s.stopCh)

	select {
	case <-s.doneCh:
	case <-ctx.Done():
		if s.stopCancelFn != nil {
			s.stopCancelFn()
		}
	}
}
