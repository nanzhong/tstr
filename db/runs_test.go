//go:build integration

package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunsQueries(t *testing.T) {
	ctx := context.Background()

	withTestDB(t, func(db DBTX) {
		querier := New()

		var runner Runner

		var (
			testA Test
			testB Test
			testC Test
		)

		var (
			runA Run
			runB Run
			runC Run
		)

		var (
			runConfig    pgtype.JSONB
			labels       pgtype.JSONB
			matrix       pgtype.JSONB
			testMatrixID uuid.UUID
		)

		err := runConfig.Set(&TestRunConfig{
			ContainerImage: "image",
			Command:        "cmd",
			Args:           []string{"a"},
			Env:            map[string]string{"env": "value"},
			TimeoutSeconds: 60,
		})
		require.NoError(t, err)

		err = labels.Set(map[string]string{"label": "value"})
		require.NoError(t, err)

		err = matrix.Set(&TestMatrix{
			Labels: map[string][]string{"label": {"value_1", "value_2"}},
		})
		require.NoError(t, err)

		testMatrixID = uuid.New()

		testA, err = querier.RegisterTest(ctx, db, RegisterTestParams{
			Namespace:    "ns1",
			Name:         "testA",
			RunConfig:    runConfig,
			Labels:       labels,
			Matrix:       matrix,
			CronSchedule: sql.NullString{Valid: true, String: "* * * * *"},
		})
		require.NoError(t, err)

		testB, err = querier.RegisterTest(ctx, db, RegisterTestParams{
			Namespace:    "ns1",
			Name:         "testB",
			RunConfig:    runConfig,
			Labels:       labels,
			Matrix:       matrix,
			CronSchedule: sql.NullString{Valid: true, String: "* * * * *"},
		})
		require.NoError(t, err)

		testC, err = querier.RegisterTest(ctx, db, RegisterTestParams{
			Namespace:    "ns2",
			Name:         "testC",
			RunConfig:    runConfig,
			Labels:       labels,
			Matrix:       matrix,
			CronSchedule: sql.NullString{Valid: true, String: "* * * * *"},
		})
		require.NoError(t, err)

		t.Run("ScheduleRun", func(t *testing.T) {
			var err error
			runA, err = querier.ScheduleRun(ctx, db, ScheduleRunParams{
				Labels:       labels,
				TestMatrixID: uuid.NullUUID{UUID: testMatrixID, Valid: true},
				TestID:       testA.ID,
				Namespace:    testA.Namespace,
			})
			require.NoError(t, err)
			assert.NotZero(t, runA.ID)
			assert.Equal(t, testA.ID, runA.TestID)
			assert.Equal(t, testA.RunConfig, runA.TestRunConfig)
			assert.Equal(t, testA.RunConfig, runA.TestRunConfig)
			assert.NotZero(t, runA.TestMatrixID)
			assert.Equal(t, testA.Labels, runA.Labels)
			assert.Zero(t, runA.RunnerID.UUID)
			assert.Equal(t, RunResultUnknown, runA.Result.RunResult)
			assert.Equal(t, pgtype.JSONB{
				Bytes:  nil,
				Status: pgtype.Null,
			}, runA.Logs)
			assert.Equal(t, pgtype.JSONB{
				Bytes:  []byte("{}"),
				Status: pgtype.Present,
			}, runA.ResultData)
			assert.NotZero(t, runA.ScheduledAt)

			runB, err = querier.ScheduleRun(ctx, db, ScheduleRunParams{
				Labels:       labels,
				TestMatrixID: uuid.NullUUID{UUID: testMatrixID, Valid: true},
				TestID:       testB.ID,
				Namespace:    testB.Namespace,
			})
			require.NoError(t, err)

			runC, err = querier.ScheduleRun(ctx, db, ScheduleRunParams{
				Labels:       labels,
				TestMatrixID: uuid.NullUUID{UUID: testMatrixID, Valid: true},
				TestID:       testC.ID,
				Namespace:    testC.Namespace,
			})
			require.NoError(t, err)
		})

		t.Run("GetRun", func(t *testing.T) {
			run, err := querier.GetRun(ctx, db, GetRunParams{
				Namespace: testA.Namespace,
				ID:        runA.ID,
			})
			require.NoError(t, err)
			assert.Equal(t, runA, run)
		})

		t.Run("ListRuns", func(t *testing.T) {
			runs, err := querier.ListRuns(ctx, db, testA.Namespace)
			require.NoError(t, err)
			assert.Equal(t, []Run{runA, runB}, runs)
		})

		t.Run("ListPendingRuns", func(t *testing.T) {
			runs, err := querier.ListPendingRuns(ctx, db)
			require.NoError(t, err)
			assert.Equal(t, []ListPendingRunsRow{
				toListPendingRunsRow(runA, testA.Namespace),
				toListPendingRunsRow(runB, testB.Namespace),
				toListPendingRunsRow(runC, testC.Namespace),
			}, runs)
		})

		t.Run("DeleteRunsForTest", func(t *testing.T) {
			err := querier.DeleteRunsForTest(ctx, db, testC.ID)
			require.NoError(t, err)

			_, err = querier.GetRun(ctx, db, GetRunParams{
				Namespace: testC.Namespace,
				ID:        runC.ID,
			})
			assert.ErrorIs(t, err, pgx.ErrNoRows)
		})

		t.Run("AssignRun", func(t *testing.T) {
			var acceptSelectors pgtype.JSONB
			err := acceptSelectors.Set(map[string]string{"label": ".*"})
			require.NoError(t, err)

			runner, err = querier.RegisterRunner(ctx, db, RegisterRunnerParams{
				Name:                     "a",
				NamespaceSelectors:       []string{"a"},
				AcceptTestLabelSelectors: acceptSelectors,
				RejectTestLabelSelectors: pgtype.JSONB{
					Bytes:  nil,
					Status: pgtype.Null,
				},
			})
			require.NoError(t, err)

			run, err := querier.AssignRun(ctx, db, AssignRunParams{
				RunnerID: runner.ID,
				RunIDs:   []uuid.UUID{runA.ID, runB.ID},
			})
			require.NoError(t, err)
			runA.RunnerID = uuid.NullUUID{UUID: runner.ID, Valid: true}
			assert.Equal(t, runA, run)

			run, err = querier.AssignRun(ctx, db, AssignRunParams{
				RunnerID: runner.ID,
				RunIDs:   []uuid.UUID{runA.ID, runB.ID},
			})
			require.NoError(t, err)
			runB.RunnerID = uuid.NullUUID{UUID: runner.ID, Valid: true}
			assert.Equal(t, runB, run)
		})

		t.Run("UpdateRun", func(t *testing.T) {
			start := time.Now().Truncate(time.Second)
			finish := start.Add(time.Minute)

			err := querier.UpdateRun(ctx, db, UpdateRunParams{
				ID: runA.ID,
				Result: NullRunResult{
					RunResult: RunResultPass,
					Valid:     true,
				},
				StartedAt:  sql.NullTime{Time: start, Valid: true},
				FinishedAt: sql.NullTime{Time: finish, Valid: true},
			})
			require.NoError(t, err)

			run, err := querier.GetRun(ctx, db, GetRunParams{
				Namespace: testA.Namespace,
				ID:        runA.ID,
			})
			require.NoError(t, err)
			runA = run
			assert.Equal(t, NullRunResult{RunResult: RunResultPass, Valid: true}, run.Result)
			assert.Equal(t, sql.NullTime{Time: start, Valid: true}, run.StartedAt)
			assert.Equal(t, sql.NullTime{Time: finish, Valid: true}, run.FinishedAt)
		})

		t.Run("AppendLogsToRun", func(t *testing.T) {
			logs := []RunLog{{
				Type: "tstr",
				Time: time.Now().Format(time.RFC3339),
				Data: []byte("log"),
			}}
			var runLogs pgtype.JSONB
			err = runLogs.Set(logs)
			require.NoError(t, err)

			err := querier.AppendLogsToRun(ctx, db, AppendLogsToRunParams{
				Logs: runLogs,
				ID:   runA.ID,
			})
			require.NoError(t, err)

			run, err := querier.GetRun(ctx, db, GetRunParams{
				Namespace: testA.Namespace,
				ID:        runA.ID,
			})
			require.NoError(t, err)
			runA = run
			var gotLogs []RunLog
			err = run.Logs.AssignTo(&gotLogs)
			require.NoError(t, err)
			assert.Equal(t, logs, gotLogs)
		})

		t.Run("UpdateResultData", func(t *testing.T) {
			data := map[string]string{"a": "123"}
			var resultData pgtype.JSONB
			err := resultData.Set(data)
			require.NoError(t, err)

			err = querier.UpdateResultData(ctx, db, UpdateResultDataParams{
				ResultData: resultData,
				ID:         runA.ID,
			})
			require.NoError(t, err)

			run, err := querier.GetRun(ctx, db, GetRunParams{
				Namespace: testA.Namespace,
				ID:        runA.ID,
			})
			require.NoError(t, err)
			runA = run
			var gotData map[string]string
			err = run.ResultData.AssignTo(&gotData)
			require.NoError(t, err)
			assert.Equal(t, data, gotData)
		})

		t.Run("ResetOrphanedRuns", func(t *testing.T) {
			err := querier.ResetOrphanedRuns(ctx, db, time.Now().Add(time.Hour))
			require.NoError(t, err)

			run, err := querier.GetRun(ctx, db, GetRunParams{
				Namespace: testB.Namespace,
				ID:        runB.ID,
			})
			require.NoError(t, err)
			assert.Zero(t, run.RunnerID.UUID)
		})

		t.Run("TimeoutRuns", func(t *testing.T) {
			_, err := querier.AssignRun(ctx, db, AssignRunParams{
				RunnerID: runner.ID,
				RunIDs:   []uuid.UUID{runB.ID},
			})
			require.NoError(t, err)

			err = querier.UpdateRun(ctx, db, UpdateRunParams{
				ID: runB.ID,
				Result: NullRunResult{
					RunResult: RunResultUnknown,
					Valid:     true,
				},
				StartedAt: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
			})
			require.NoError(t, err)

			logs := []RunLog{{
				Type: "tstr",
				Time: time.Now().Format(time.RFC3339),
				Data: []byte("timeout"),
			}}
			var timeoutLogs pgtype.JSONB
			err = timeoutLogs.Set(logs)
			require.NoError(t, err)
			err = querier.TimeoutRuns(ctx, db, TimeoutRunsParams{
				TimeoutLog:     timeoutLogs,
				DefaultTimeout: int32(time.Minute.Seconds()),
			})
			require.NoError(t, err)

			run, err := querier.GetRun(ctx, db, GetRunParams{
				Namespace: testB.Namespace,
				ID:        runB.ID,
			})
			require.NoError(t, err)
			runB = run
			assert.Equal(t, NullRunResult{RunResult: RunResultError, Valid: true}, run.Result)
			assert.NotZero(t, run.FinishedAt)
			var gotLogs []RunLog
			err = run.Logs.AssignTo(&gotLogs)
			require.NoError(t, err)
			assert.Equal(t, logs, gotLogs)
		})

		t.Run("RunSummariesForTest", func(t *testing.T) {
			runs, err := querier.RunSummariesForTest(ctx, db, RunSummariesForTestParams{
				Namespace:      testA.Namespace,
				TestID:         testA.ID,
				ScheduledAfter: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
			})
			require.NoError(t, err)
			assert.Equal(t, []RunSummariesForTestRow{toRunSummariesForTestRow(runA, testA)}, runs)
		})

		t.Run("RunSummariesForRunner", func(t *testing.T) {
			runs, err := querier.RunSummariesForRunner(ctx, db, RunSummariesForRunnerParams{
				Namespace:      testA.Namespace,
				RunnerID:       uuid.NullUUID{UUID: runner.ID, Valid: true},
				ScheduledAfter: sql.NullTime{Time: time.Now().Add(-time.Hour), Valid: true},
			})
			require.NoError(t, err)
			assert.Equal(t, []RunSummariesForRunnerRow{
				toRunSummariesForRunnerRow(runA, testA),
				toRunSummariesForRunnerRow(runB, testB),
			}, runs)
		})

		t.Run("QueryRuns", func(t *testing.T) {
			t.Skipf("TODO: need to implement")
		})

		t.Run("SummarizeRunsBreakdownResult", func(t *testing.T) {
			t.Skipf("TODO: need to implement")
		})

		t.Run("SummarizeRunsBreakdownTest", func(t *testing.T) {
			t.Skipf("TODO: need to implement")
		})
	})
}

func toListPendingRunsRow(r Run, ns string) ListPendingRunsRow {
	return ListPendingRunsRow{
		ID:            r.ID,
		TestID:        r.TestID,
		TestRunConfig: r.TestRunConfig,
		TestMatrixID:  r.TestMatrixID,
		Labels:        r.Labels,
		RunnerID:      r.RunnerID,
		Result:        r.Result,
		Logs:          r.Logs,
		ResultData:    r.ResultData,
		ScheduledAt:   r.ScheduledAt,
		StartedAt:     r.StartedAt,
		FinishedAt:    r.FinishedAt,
		Namespace:     ns,
	}
}

func toRunSummariesForTestRow(r Run, t Test) RunSummariesForTestRow {
	return RunSummariesForTestRow{
		ID:            r.ID,
		TestNamespace: t.Namespace,
		TestID:        t.ID,
		TestName:      t.Name,
		TestRunConfig: t.RunConfig,
		Labels:        t.Labels,
		RunnerID:      r.RunnerID,
		Result:        r.Result,
		ScheduledAt:   r.ScheduledAt,
		StartedAt:     r.StartedAt,
		FinishedAt:    r.FinishedAt,
		ResultData:    r.ResultData,
	}
}

func toRunSummariesForRunnerRow(r Run, t Test) RunSummariesForRunnerRow {
	return RunSummariesForRunnerRow{
		ID:            r.ID,
		TestNamespace: t.Namespace,
		TestID:        t.ID,
		TestName:      t.Name,
		TestRunConfig: t.RunConfig,
		Labels:        t.Labels,
		RunnerID:      r.RunnerID,
		Result:        r.Result,
		ScheduledAt:   r.ScheduledAt,
		StartedAt:     r.StartedAt,
		FinishedAt:    r.FinishedAt,
		ResultData:    r.ResultData,
	}
}
