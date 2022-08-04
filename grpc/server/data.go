package server

import (
	"context"
	"database/sql"
	"sort"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	datav1 "github.com/nanzhong/tstr/api/data/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type DataServer struct {
	datav1.UnimplementedDataServiceServer

	pgxPool   *pgxpool.Pool
	dbQuerier db.Querier
	clock     clock.Clock
}

func NewDataServer(pgxPool *pgxpool.Pool) datav1.DataServiceServer {
	return &DataServer{
		pgxPool:   pgxPool,
		dbQuerier: db.New(),
		clock:     clock.New(),
	}
}

func (s *DataServer) GetTest(ctx context.Context, r *datav1.GetTestRequest) (*datav1.GetTestResponse, error) {
	id, err := uuid.Parse(r.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid test id")
	}

	test, err := s.dbQuerier.GetTest(ctx, s.pgxPool, id)
	if err != nil {
		log.Error().
			Err(err).
			Stringer("test_id", id).
			Msg("failed to get test")
		return nil, status.Error(codes.Internal, "failed to get test")
	}

	var runConfig db.TestRunConfig
	if err := test.RunConfig.AssignTo(&runConfig); err != nil {
		log.Error().
			Err(err).
			Stringer("test_id", test.ID).
			Msg("failed to parse run config")
		return nil, status.Error(codes.Internal, "failed to format run config")
	}
	var labels map[string]string
	if err := test.Labels.AssignTo(&labels); err != nil {
		log.Error().
			Err(err).
			Stringer("test_id", test.ID).
			Msg("failed to parse labels")
		return nil, status.Error(codes.Internal, "failed to format labels")
	}
	pbTest := &commonv1.Test{
		Id:           test.ID.String(),
		Name:         test.Name,
		Labels:       labels,
		CronSchedule: test.CronSchedule.String,
		RunConfig:    types.ToProtoTestRunConfig(runConfig),
		NextRunAt:    types.ToProtoTimestamp(test.NextRunAt),
		RegisteredAt: types.ToProtoTimestamp(test.RegisteredAt),
		UpdatedAt:    types.ToProtoTimestamp(test.UpdatedAt),
	}

	runSummaries, err := s.dbQuerier.RunSummariesForTest(ctx, s.pgxPool, db.RunSummariesForTestParams{
		TestID: test.ID,
		// TODO Configure default + query param handling
		ScheduledAfter: sql.NullTime{Valid: true, Time: s.clock.Now().Add(-24 * time.Hour)},
	})
	if err != nil {
		log.Error().Err(err).Stringer("test_id", test.ID).Msg("failed to summarize runs for test")
		return nil, status.Error(codes.Internal, "failed to summarize runs for test")
	}
	var pbRunSummaries []*datav1.RunSummary
	for _, s := range runSummaries {
		var runConfig db.TestRunConfig
		if err := test.RunConfig.AssignTo(&runConfig); err != nil {
			log.Error().
				Err(err).
				Stringer("run_id", s.ID).
				Msg("failed to parse run config")
			return nil, status.Error(codes.Internal, "failed to format run config")
		}

		pbRunSummaries = append(pbRunSummaries, &datav1.RunSummary{
			Id:            s.ID.String(),
			TestId:        s.TestID.String(),
			TestName:      s.TestName,
			TestRunConfig: types.ToProtoTestRunConfig(runConfig),
			RunnerId:      s.RunnerID.UUID.String(),
			Result:        types.ToRunResult(s.Result.RunResult),
			ScheduledAt:   types.ToProtoTimestamp(s.ScheduledAt),
			StartedAt:     types.ToProtoTimestamp(s.StartedAt),
			FinishedAt:    types.ToProtoTimestamp(s.FinishedAt),
			ResultData:    types.ToProtoResultData(s.ResultData),
		})
	}

	return &datav1.GetTestResponse{
		Test:         pbTest,
		RunSummaries: pbRunSummaries,
	}, nil
}

func (s *DataServer) QueryTests(ctx context.Context, r *datav1.QueryTestsRequest) (*datav1.QueryTestsResponse, error) {
	var testIDs []uuid.UUID
	for _, rid := range r.Ids {
		id, err := uuid.Parse(rid)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to parse test id")
		}
		testIDs = append(testIDs, id)
	}

	var testSuiteIDs []uuid.UUID
	for _, rid := range r.TestSuiteIds {
		id, err := uuid.Parse(rid)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to parse test suite id")
		}
		testSuiteIDs = append(testSuiteIDs, id)
	}

	labels := pgtype.JSONB{Status: pgtype.Null}
	if r.Labels != nil {
		if err := labels.Set(&r.Labels); err != nil {
			log.Error().Err(err).Msg("failed to parse labels")
			return nil, status.Error(codes.InvalidArgument, "failed to parse labels")
		}
	}

	tests, err := s.dbQuerier.QueryTests(ctx, s.pgxPool, db.QueryTestsParams{
		Ids:          testIDs,
		TestSuiteIds: testSuiteIDs,
		Labels:       labels,
	})
	if err != nil {
		// NOTE []uuid.UUID can't be directly used as []fmt.Stringer
		testIDStrings := make([]string, len(testIDs))
		for i, s := range testIDs {
			testIDStrings[i] = s.String()
		}
		testSuiteIDStrings := make([]string, len(testSuiteIDs))
		for i, s := range testSuiteIDs {
			testSuiteIDStrings[i] = s.String()
		}
		log.Error().
			Err(err).
			Strs("test_ids", testIDStrings).
			Strs("test_suite_ids", testSuiteIDStrings).
			Dict("labels", zerolog.Dict().Fields(labels)).
			Msg("failed to query tests")
	}

	var pbTests []*commonv1.Test
	for _, test := range tests {
		var labels map[string]string
		if err := test.Labels.AssignTo(&labels); err != nil {
			log.Error().
				Err(err).
				Stringer("test_id", test.ID).
				Msg("failed to parse labels")
			return nil, status.Error(codes.Internal, "failed to format labels")
		}
		var runConfig db.TestRunConfig
		if err := test.RunConfig.AssignTo(&runConfig); err != nil {
			log.Error().
				Err(err).
				Stringer("test_id", test.ID).
				Msg("failed to parse run config")
			return nil, status.Error(codes.Internal, "failed to format run config")
		}
		pbTests = append(pbTests, &commonv1.Test{
			Id:           test.ID.String(),
			Name:         test.Name,
			Labels:       labels,
			CronSchedule: test.CronSchedule.String,
			RunConfig:    types.ToProtoTestRunConfig(runConfig),
			NextRunAt:    types.ToProtoTimestamp(test.NextRunAt),
			RegisteredAt: types.ToProtoTimestamp(test.RegisteredAt),
			UpdatedAt:    types.ToProtoTimestamp(test.UpdatedAt),
		})
	}

	return &datav1.QueryTestsResponse{
		Tests: pbTests,
	}, nil
}

func (s *DataServer) GetTestSuite(ctx context.Context, r *datav1.GetTestSuiteRequest) (*datav1.GetTestSuiteResponse, error) {
	id, err := uuid.Parse(r.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid test suite id")
	}

	testSuite, err := s.dbQuerier.GetTestSuite(ctx, s.pgxPool, id)
	if err != nil {
		log.Error().
			Err(err).
			Stringer("test_suite_id", id).
			Msg("failed to get test suite")
		return nil, status.Error(codes.Internal, "failed to get test suite")
	}

	var labels map[string]string
	if err := testSuite.Labels.AssignTo(&labels); err != nil {
		log.Error().
			Err(err).
			Stringer("test_suite_id", testSuite.ID).
			Msg("failed to parse labels")
		return nil, status.Error(codes.Internal, "failed to format labels")
	}
	pbTestSuite := &commonv1.TestSuite{
		Id:        testSuite.ID.String(),
		Name:      testSuite.Name,
		Labels:    labels,
		CreatedAt: types.ToProtoTimestamp(testSuite.CreatedAt),
		UpdatedAt: types.ToProtoTimestamp(testSuite.UpdatedAt),
	}

	return &datav1.GetTestSuiteResponse{
		TestSuite: pbTestSuite,
	}, nil
}

func (s *DataServer) QueryTestSuites(ctx context.Context, r *datav1.QueryTestSuitesRequest) (*datav1.QueryTestSuitesResponse, error) {
	var testSuiteIDs []uuid.UUID
	for _, rid := range r.Ids {
		id, err := uuid.Parse(rid)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to parse test suite id")
		}
		testSuiteIDs = append(testSuiteIDs, id)
	}

	labels := pgtype.JSONB{Status: pgtype.Null}
	if r.Labels != nil {
		if err := labels.Set(&r.Labels); err != nil {
			log.Error().Err(err).Msg("failed to parse labels")
			return nil, status.Error(codes.InvalidArgument, "failed to parse labels")
		}
	}

	testSuites, err := s.dbQuerier.QueryTestSuites(ctx, s.pgxPool, db.QueryTestSuitesParams{
		Ids:    testSuiteIDs,
		Labels: labels,
	})
	if err != nil {
		log.Error().
			Err(err).
			Dict("labels", zerolog.Dict().Fields(labels)).
			Msg("failed to query test suites")
		return nil, status.Error(codes.Internal, "failed to query test suites")
	}
	var pbTestSuites []*commonv1.TestSuite
	for _, ts := range testSuites {
		var pbLabels map[string]string
		if err := ts.Labels.AssignTo(&pbLabels); err != nil {
			log.Error().Err(err).Stringer("test_suite_id", ts.ID).Msg("failed to format labels")
			return nil, status.Error(codes.Internal, "failed to format labels for test suite")
		}
		pbTestSuites = append(pbTestSuites, &commonv1.TestSuite{
			Id:        ts.ID.String(),
			Name:      ts.Name,
			Labels:    pbLabels,
			CreatedAt: types.ToProtoTimestamp(ts.CreatedAt),
			UpdatedAt: types.ToProtoTimestamp(ts.UpdatedAt),
		})
	}

	return &datav1.QueryTestSuitesResponse{
		TestSuites: pbTestSuites,
	}, nil
}

func (s *DataServer) GetRun(ctx context.Context, r *datav1.GetRunRequest) (*datav1.GetRunResponse, error) {
	id, err := uuid.Parse(r.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid run id")
	}

	run, err := s.dbQuerier.GetRun(ctx, s.pgxPool, id)
	if err != nil {
		log.Error().
			Err(err).
			Stringer("run_id", id).
			Msg("failed to get run")
		return nil, status.Error(codes.Internal, "failed to get run")
	}

	var runConfig db.TestRunConfig
	if err := run.TestRunConfig.AssignTo(&runConfig); err != nil {
		log.Error().
			Err(err).
			Stringer("run_id", run.ID).
			Msg("failed to parse run config")
		return nil, status.Error(codes.Internal, "failed to format run config")
	}

	var logs []db.RunLog
	if err := run.Logs.AssignTo(&logs); err != nil {
		log.Error().
			Err(err).
			Stringer("run_id", run.ID).
			Msg("failed to parse logs")
		return nil, status.Error(codes.Internal, "failed to format run logs")
	}
	pbRun := &commonv1.Run{
		Id:            run.ID.String(),
		TestId:        run.TestID.String(),
		TestRunConfig: types.ToProtoTestRunConfig(runConfig),
		Result:        types.ToRunResult(run.Result.RunResult),
		Logs:          types.ToRunLogs(logs),
		ScheduledAt:   types.ToProtoTimestamp(run.ScheduledAt),
		StartedAt:     types.ToProtoTimestamp(run.StartedAt),
		FinishedAt:    types.ToProtoTimestamp(run.FinishedAt),
		ResultData:    types.ToProtoResultData(run.ResultData),
	}
	if run.RunnerID.Valid {
		pbRun.RunnerId = run.RunnerID.UUID.String()
	}

	return &datav1.GetRunResponse{
		Run: pbRun,
	}, nil
}

func (s *DataServer) QueryRuns(ctx context.Context, r *datav1.QueryRunsRequest) (*datav1.QueryRunsResponse, error) {
	var runIDs []uuid.UUID
	for _, rid := range r.Ids {
		id, err := uuid.Parse(rid)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to parse run id")
		}
		runIDs = append(runIDs, id)
	}

	var testIDs []uuid.UUID
	for _, rid := range r.TestIds {
		id, err := uuid.Parse(rid)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to parse test id")
		}
		testIDs = append(testIDs, id)
	}

	var testSuiteIDs []uuid.UUID
	for _, rid := range r.TestSuiteIds {
		id, err := uuid.Parse(rid)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to parse test suite id")
		}
		testSuiteIDs = append(testSuiteIDs, id)
	}

	var runnerIDs []uuid.UUID
	for _, rid := range r.RunnerIds {
		id, err := uuid.Parse(rid)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to parse runner id")
		}
		runnerIDs = append(runnerIDs, id)
	}

	results := types.FromRunResults(r.Results)
	scheduledBefore := types.FromProtoTimestampAsNullTime(r.ScheduledBefore)
	scheduledAfter := types.FromProtoTimestampAsNullTime(r.ScheduledAfter)
	startedBefore := types.FromProtoTimestampAsNullTime(r.StartedBefore)
	startedAfter := types.FromProtoTimestampAsNullTime(r.StartedAfter)
	finishedBefore := types.FromProtoTimestampAsNullTime(r.FinishedBefore)
	finishedAfter := types.FromProtoTimestampAsNullTime(r.FinishedAfter)

	runs, err := s.dbQuerier.QueryRuns(ctx, s.pgxPool, db.QueryRunsParams{
		Ids:             runIDs,
		TestIds:         testIDs,
		TestSuiteIds:    testSuiteIDs,
		RunnerIds:       runnerIDs,
		Results:         results,
		ScheduledBefore: scheduledBefore,
		ScheduledAfter:  scheduledAfter,
		StartedBefore:   startedBefore,
		StartedAfter:    startedAfter,
		FinishedBefore:  finishedBefore,
		FinishedAfter:   finishedAfter,
	})
	if err != nil {
		resultStrings := make([]string, len(results))
		for i, r := range results {
			resultStrings[i] = string(r)
		}
		log.Error().
			Err(err).
			Strs("ids", r.Ids).
			Strs("test_ids", r.TestIds).
			Strs("test_suite_ids", r.TestSuiteIds).
			Strs("results", resultStrings).
			Time("scheduled_before", scheduledBefore.Time).
			Time("scheduled_after", scheduledAfter.Time).
			Time("started_before", startedBefore.Time).
			Time("started_after", startedAfter.Time).
			Time("finished_before", finishedBefore.Time).
			Time("finished_after", finishedAfter.Time).
			Msg("failed to query runs")
		return nil, status.Error(codes.Internal, "failed to query runs")
	}

	var pbRuns []*commonv1.Run
	for _, run := range runs {
		pbRun, err := types.ToProtoRun(&run)
		if err != nil {
			log.Error().
				Err(err).
				Stringer("run_id", run.ID).
				Msg("failed to format run")
			return nil, status.Error(codes.Internal, "failed to format run")
		}
		pbRuns = append(pbRuns, pbRun)
	}

	return &datav1.QueryRunsResponse{
		Runs: pbRuns,
	}, nil
}

func (s *DataServer) GetRunner(ctx context.Context, r *datav1.GetRunnerRequest) (*datav1.GetRunnerResponse, error) {
	id, err := uuid.Parse(r.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid runner id")
	}

	runner, err := s.dbQuerier.GetRunner(ctx, s.pgxPool, id)
	if err != nil {
		log.Error().Err(err).Stringer("runner_id", id).Msg("failed to get runner")
		return nil, status.Error(codes.Internal, "failed to get runner")
	}

	var (
		acceptSelectors map[string]string
		rejectSelectors map[string]string
	)
	if err := runner.AcceptTestLabelSelectors.AssignTo(&acceptSelectors); err != nil {
		log.Error().Err(err).Stringer("runner_id", runner.ID).Msg("failed to format accept label selectors")
		return nil, status.Error(codes.Internal, "failed to format accept label selectors")
	}
	if err := runner.RejectTestLabelSelectors.AssignTo(&rejectSelectors); err != nil {
		log.Error().Err(err).Stringer("runner_id", runner.ID).Msg("failed to format reject label selectors")
		return nil, status.Error(codes.Internal, "failed to format reject label selectors")
	}

	pbRunner := &commonv1.Runner{
		Id:                       runner.ID.String(),
		Name:                     runner.Name,
		AcceptTestLabelSelectors: acceptSelectors,
		RejectTestLabelSelectors: rejectSelectors,
		RegisteredAt:             types.ToProtoTimestamp(runner.RegisteredAt),
		LastHeartbeatAt:          types.ToProtoTimestamp(runner.LastHeartbeatAt),
	}

	runSummaries, err := s.dbQuerier.RunSummariesForRunner(ctx, s.pgxPool, db.RunSummariesForRunnerParams{
		RunnerID: uuid.NullUUID{Valid: true, UUID: runner.ID},
		// TODO Configure default + query param handling
		ScheduledAfter: sql.NullTime{Valid: true, Time: s.clock.Now().Add(-24 * time.Hour)},
	})
	if err != nil {
		log.Error().Err(err).Stringer("runner_id", runner.ID).Msg("failed to summarize runs for runner")
		return nil, status.Error(codes.Internal, "failed to summarize runs for test")
	}
	var pbRunSummaries []*datav1.RunSummary
	for _, s := range runSummaries {
		var runConfig db.TestRunConfig
		if err := s.TestRunConfig.AssignTo(&runConfig); err != nil {
			log.Error().
				Err(err).
				Stringer("run_id", s.ID).
				Msg("failed to parse run config")
			return nil, status.Error(codes.Internal, "failed to format run config")
		}

		pbRunSummaries = append(pbRunSummaries, &datav1.RunSummary{
			Id:            s.ID.String(),
			TestId:        s.TestID.String(),
			TestName:      s.TestName,
			TestRunConfig: types.ToProtoTestRunConfig(runConfig),
			RunnerId:      s.RunnerID.UUID.String(),
			Result:        types.ToRunResult(s.Result.RunResult),
			ScheduledAt:   types.ToProtoTimestamp(s.ScheduledAt),
			StartedAt:     types.ToProtoTimestamp(s.StartedAt),
			FinishedAt:    types.ToProtoTimestamp(s.FinishedAt),
		})
	}

	return &datav1.GetRunnerResponse{
		Runner:       pbRunner,
		RunSummaries: pbRunSummaries,
	}, nil
}

func (s *DataServer) QueryRunners(ctx context.Context, r *datav1.QueryRunnersRequest) (*datav1.QueryRunnersResponse, error) {
	var runnerIDs []uuid.UUID
	for _, rid := range r.Ids {
		id, err := uuid.Parse(rid)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "failed to parse runner id")
		}
		runnerIDs = append(runnerIDs, id)
	}

	var lastHeartbeatSince sql.NullTime
	if r.LastHeartbeatWithin != nil {
		lastHeartbeatSince.Valid = true
		lastHeartbeatSince.Time = s.clock.Now().Add(-r.LastHeartbeatWithin.AsDuration())
	}

	runners, err := s.dbQuerier.QueryRunners(ctx, s.pgxPool, db.QueryRunnersParams{
		Ids:                runnerIDs,
		LastHeartbeatSince: lastHeartbeatSince,
	})
	if err != nil {
		log.Error().Err(err).Time("last_heartbeat_since", lastHeartbeatSince.Time).Msg("failed to query runners")
		return nil, status.Error(codes.Internal, "failed to query runners")
	}
	var pbRunners []*commonv1.Runner
	for _, r := range runners {
		var (
			acceptSelectors map[string]string
			rejectSelectors map[string]string
		)
		if err := r.AcceptTestLabelSelectors.AssignTo(&acceptSelectors); err != nil {
			log.Error().
				Err(err).
				Stringer("runner_id", r.ID).
				Msg("failed to format accept selectors")
			return nil, status.Error(codes.Internal, "failed to format accept selectorr")
		}
		if err := r.RejectTestLabelSelectors.AssignTo(&rejectSelectors); err != nil {
			log.Error().
				Err(err).
				Stringer("runner_id", r.ID).
				Msg("failed to format reject selectors")
			return nil, status.Error(codes.Internal, "failed to format reject selectorr")
		}

		pbRunners = append(pbRunners, &commonv1.Runner{
			Id:                       r.ID.String(),
			Name:                     r.Name,
			AcceptTestLabelSelectors: acceptSelectors,
			RejectTestLabelSelectors: rejectSelectors,
			RegisteredAt:             types.ToProtoTimestamp(r.RegisteredAt),
			LastHeartbeatAt:          types.ToProtoTimestamp(r.LastHeartbeatAt),
		})
	}

	return &datav1.QueryRunnersResponse{
		Runners: pbRunners,
	}, nil
}

func (s *DataServer) SummarizeRuns(ctx context.Context, r *datav1.SummarizeRunsRequest) (*datav1.SummarizeRunsResponse, error) {
	var precision string
	var interval time.Duration
	switch r.Interval {
	case datav1.SummarizeRunsRequest_HOUR:
		precision = "hour"
		interval = time.Hour
	case datav1.SummarizeRunsRequest_DAY:
		precision = "hour"
		interval = 24 * time.Hour
	case datav1.SummarizeRunsRequest_WEEK:
		precision = "week"
		interval = 7 * 24 * time.Hour
	case datav1.SummarizeRunsRequest_UNKNOWN:
		precision = "hour"
		interval = time.Hour
	}
	startTime := r.ScheduledAfter.AsTime()
	endTime := r.ScheduledAfter.AsTime().Add(r.Window.AsDuration())

	var stats []*datav1.SummarizeRunsResponse_IntervalStats
	err := s.pgxPool.BeginFunc(ctx, func(tx pgx.Tx) error {
		log := log.With().
			Str("precision", precision).
			Time("start_time", startTime).
			Time("end_time", endTime).
			Dur("duration", interval).
			Logger()

		resultBreakdwon, err := s.dbQuerier.SummarizeRunsBreakdownResult(ctx, tx, db.SummarizeRunsBreakdownResultParams{
			Precision: precision,
			StartTime: sql.NullTime{Valid: true, Time: startTime},
			EndTime:   sql.NullTime{Valid: true, Time: endTime},
			Interval:  interval.Seconds(),
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to sumarize runs - result breakdown")
			return err
		}

		testBreakdwon, err := s.dbQuerier.SummarizeRunsBreakdownTest(ctx, tx, db.SummarizeRunsBreakdownTestParams{
			Precision: precision,
			StartTime: sql.NullTime{Valid: true, Time: startTime},
			EndTime:   sql.NullTime{Valid: true, Time: endTime},
			Interval:  interval.Seconds(),
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to sumarize runs - test breakdown")
			return err
		}

		statsMap := map[time.Time]*datav1.SummarizeRunsResponse_IntervalStats{}
		for _, result := range resultBreakdwon {
			statsMap[result.IntervalsStart] = &datav1.SummarizeRunsResponse_IntervalStats{
				StartTime: types.ToProtoTimestamp(result.IntervalsStart),
				Duration:  durationpb.New(interval),
				ResultCount: []*datav1.SummarizeRunsResponse_IntervalStats_ResultBreakdown{
					{Result: commonv1.Run_PASS, Count: uint32(result.Pass)},
					{Result: commonv1.Run_FAIL, Count: uint32(result.Fail)},
					{Result: commonv1.Run_ERROR, Count: uint32(result.Error)},
					{Result: commonv1.Run_UNKNOWN, Count: uint32(result.Unknown)},
				},
			}
		}
		for _, test := range testBreakdwon {
			statsMap[test.IntervalsStart].TestCount = append(statsMap[test.IntervalsStart].TestCount, &datav1.SummarizeRunsResponse_IntervalStats_TestBreakdown{
				TestId:   test.ID.String(),
				TestName: test.Name,
				ResultCount: []*datav1.SummarizeRunsResponse_IntervalStats_ResultBreakdown{
					{Result: commonv1.Run_PASS, Count: uint32(test.Pass)},
					{Result: commonv1.Run_FAIL, Count: uint32(test.Fail)},
					{Result: commonv1.Run_ERROR, Count: uint32(test.Error)},
					{Result: commonv1.Run_UNKNOWN, Count: uint32(test.Unknown)},
				},
			})
		}

		for _, v := range statsMap {
			stats = append(stats, v)
		}
		sort.Slice(stats, func(i, j int) bool {
			return stats[i].StartTime.AsTime().Before(stats[j].StartTime.AsTime())
		})
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to summarize runs")
	}
	return &datav1.SummarizeRunsResponse{
		IntervalStats: stats,
	}, nil
}
