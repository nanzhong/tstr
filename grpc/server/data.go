package server

import (
	"context"
	"database/sql"

	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	commonv1 "github.com/nanzhong/tstr/api/common/v1"
	datav1 "github.com/nanzhong/tstr/api/data/v1"
	"github.com/nanzhong/tstr/db"
	"github.com/nanzhong/tstr/grpc/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	var env map[string]string
	if err := test.Env.AssignTo(&env); err != nil {
		log.Error().
			Err(err).
			Stringer("test_id", test.ID).
			Msg("failed to parse env")
		return nil, status.Error(codes.Internal, "failed to format run config env")
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
		RunConfig: &commonv1.Test_RunConfig{
			Id:             test.TestRunConfigID.UUID.String(),
			ContainerImage: test.ContainerImage.String,
			Command:        test.Command.String,
			Args:           test.Args,
			Env:            env,
			CreatedAt:      types.ToProtoTimestamp(test.CreatedAt),
		},
		NextRunAt:    types.ToProtoTimestamp(test.NextRunAt),
		RegisteredAt: types.ToProtoTimestamp(test.RegisteredAt),
		UpdatedAt:    types.ToProtoTimestamp(test.UpdatedAt),
		ArchivedAt:   types.ToProtoTimestamp(test.ArchivedAt),
	}

	runSummaries, err := s.dbQuerier.RunSummaryForTest(ctx, s.pgxPool, db.RunSummaryForTestParams{
		TestID: test.ID,
		// TODO Configure default + query param handling
		Limit: 200,
	})
	if err != nil {
		log.Error().Err(err).Stringer("test_id", test.ID).Msg("failed to summarize runs for test")
		return nil, status.Error(codes.Internal, "failed to summarize runs for test")
	}
	var pbRunSummaries []*datav1.RunSummary
	for _, s := range runSummaries {
		pbRunSummaries = append(pbRunSummaries, &datav1.RunSummary{
			Id:              s.ID.String(),
			TestId:          s.TestID.String(),
			TestRunConfigId: s.TestRunConfigID.String(),
			RunnerId:        s.RunnerID.UUID.String(),
			Result:          types.ToRunResult(s.Result.RunResult),
			ScheduledAt:     types.ToProtoTimestamp(s.ScheduledAt),
			StartedAt:       types.ToProtoTimestamp(s.StartedAt),
			FinishedAt:      types.ToProtoTimestamp(s.FinishedAt),
			ComputedData:    types.ToProtoComputedData(s.ComputedData),
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
		var env map[string]string
		if err := test.Env.AssignTo(&env); err != nil {
			log.Error().
				Err(err).
				Stringer("test_id", test.ID).
				Msg("failed to parse env")
			return nil, status.Error(codes.Internal, "failed to format run config env")
		}
		pbTests = append(pbTests, &commonv1.Test{
			Id:           test.ID.String(),
			Name:         test.Name,
			Labels:       labels,
			CronSchedule: test.CronSchedule.String,
			RunConfig: &commonv1.Test_RunConfig{
				Id:             test.TestRunConfigID.UUID.String(),
				ContainerImage: test.ContainerImage.String,
				Command:        test.Command.String,
				Args:           test.Args,
				Env:            env,
				CreatedAt:      types.ToProtoTimestamp(test.CreatedAt),
			},
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

	var env map[string]string
	if err := run.Env.AssignTo(&env); err != nil {
		log.Error().
			Err(err).
			Stringer("run_id", run.ID).
			Msg("failed to parse env")
		return nil, status.Error(codes.Internal, "failed to format run config env")
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
		Result:       types.ToRunResult(run.Result.RunResult),
		Logs:         types.ToRunLogs(logs),
		ScheduledAt:  types.ToProtoTimestamp(run.ScheduledAt),
		StartedAt:    types.ToProtoTimestamp(run.StartedAt),
		FinishedAt:   types.ToProtoTimestamp(run.FinishedAt),
		ComputedData: types.ToProtoComputedData(run.ComputedData),
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
		var env map[string]string
		if err := run.Env.AssignTo(&env); err != nil {
			log.Error().
				Err(err).
				Stringer("run_id", run.ID).
				Msg("failed to parse env")
			return nil, status.Error(codes.Internal, "failed to format run config env")
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
			Result:      types.ToRunResult(run.Result.RunResult),
			Logs:        types.ToRunLogs(logs),
			ScheduledAt: types.ToProtoTimestamp(run.ScheduledAt),
			StartedAt:   types.ToProtoTimestamp(run.StartedAt),
			FinishedAt:  types.ToProtoTimestamp(run.FinishedAt),
		}
		if run.RunnerID.Valid {
			pbRun.RunnerId = run.RunnerID.UUID.String()
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

	runSummaries, err := s.dbQuerier.RunSummaryForRunner(ctx, s.pgxPool, db.RunSummaryForRunnerParams{
		RunnerID: runner.ID,
		// TODO Configure default + query param handling
		Limit: 200,
	})
	if err != nil {
		log.Error().Err(err).Stringer("runner_id", runner.ID).Msg("failed to summarize runs for runner")
		return nil, status.Error(codes.Internal, "failed to summarize runs for test")
	}
	var pbRunSummaries []*datav1.RunSummary
	for _, s := range runSummaries {
		pbRunSummaries = append(pbRunSummaries, &datav1.RunSummary{
			Id:              s.ID.String(),
			TestId:          s.TestID.String(),
			TestRunConfigId: s.TestRunConfigID.String(),
			RunnerId:        s.RunnerID.UUID.String(),
			Result:          types.ToRunResult(s.Result.RunResult),
			ScheduledAt:     types.ToProtoTimestamp(s.ScheduledAt),
			StartedAt:       types.ToProtoTimestamp(s.StartedAt),
			FinishedAt:      types.ToProtoTimestamp(s.FinishedAt),
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
