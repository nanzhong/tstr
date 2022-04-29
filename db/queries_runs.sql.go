// Code generated by pggen. DO NOT EDIT.

package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"time"
)

const getRunSQL = `SELECT runs.*, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at
FROM runs
JOIN test_run_configs
ON runs.test_run_config_id = test_run_configs.id
WHERE runs.id = $1;`

type GetRunRow struct {
	ID              string       `json:"id"`
	TestID          string       `json:"test_id"`
	TestRunConfigID int          `json:"test_run_config_id"`
	RunnerID        string       `json:"runner_id"`
	Result          RunResult    `json:"result"`
	Logs            pgtype.JSONB `json:"logs"`
	ScheduledAt     time.Time    `json:"scheduled_at"`
	StartedAt       time.Time    `json:"started_at"`
	FinishedAt      time.Time    `json:"finished_at"`
	ContainerImage  string       `json:"container_image"`
	Command         string       `json:"command"`
	Args            []string     `json:"args"`
	Env             pgtype.JSONB `json:"env"`
	CreatedAt       time.Time    `json:"created_at"`
}

// GetRun implements Querier.GetRun.
func (q *DBQuerier) GetRun(ctx context.Context, id string) (GetRunRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetRun")
	row := q.conn.QueryRow(ctx, getRunSQL, id)
	var item GetRunRow
	if err := row.Scan(&item.ID, &item.TestID, &item.TestRunConfigID, &item.RunnerID, &item.Result, &item.Logs, &item.ScheduledAt, &item.StartedAt, &item.FinishedAt, &item.ContainerImage, &item.Command, &item.Args, &item.Env, &item.CreatedAt); err != nil {
		return item, fmt.Errorf("query GetRun: %w", err)
	}
	return item, nil
}

// GetRunBatch implements Querier.GetRunBatch.
func (q *DBQuerier) GetRunBatch(batch genericBatch, id string) {
	batch.Queue(getRunSQL, id)
}

// GetRunScan implements Querier.GetRunScan.
func (q *DBQuerier) GetRunScan(results pgx.BatchResults) (GetRunRow, error) {
	row := results.QueryRow()
	var item GetRunRow
	if err := row.Scan(&item.ID, &item.TestID, &item.TestRunConfigID, &item.RunnerID, &item.Result, &item.Logs, &item.ScheduledAt, &item.StartedAt, &item.FinishedAt, &item.ContainerImage, &item.Command, &item.Args, &item.Env, &item.CreatedAt); err != nil {
		return item, fmt.Errorf("scan GetRunBatch row: %w", err)
	}
	return item, nil
}

const listRunsSQL = `SELECT runs.*, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at
FROM runs
JOIN test_run_configs
ON runs.test_run_config_id = test_run_configs.id
WHERE
  CASE WHEN $1
    THEN runs.test_id = ANY ($2::uuid[])
    ELSE TRUE
  END AND
  CASE WHEN $3
    THEN runs.test_id = ANY (
      SELECT tests.id
      FROM test_suites
      JOIN tests
      ON tests.labels @> test_suites.labels
      WHERE test_suites.id = ANY ($4::uuid[])
    )
    ELSE TRUE
  END AND
  CASE WHEN $5
    THEN runner_id = ANY ($6::uuid[])
    ELSE TRUE
  END AND
  CASE WHEN $7
    THEN result = ANY ($8::run_result[])
    ELSE TRUE
  END AND
  CASE WHEN $9
    THEN scheduled_at < $10::timestamptz
    ELSE TRUE
  END AND
  CASE WHEN $11
    THEN scheduled_at > $12::timestamptz
    ELSE TRUE
  END AND
  CASE WHEN $13
    THEN started_at < $14::timestamptz
    ELSE TRUE
  END AND
  CASE WHEN $15
    THEN started_at > $16::timestamptz
    ELSE TRUE
  END AND
  CASE WHEN $17
    THEN finished_at < $18::timestamptz
    ELSE TRUE
  END AND
  CASE WHEN $19
    THEN finished_at > $20::timestamptz
    ELSE TRUE
  END;`

type ListRunsParams struct {
	FilterTestIds         bool
	TestIds               pgtype.UUIDArray
	FilterTestSuiteIds    bool
	TestSuiteIds          pgtype.UUIDArray
	FilterRunnerIds       bool
	RunnerIds             pgtype.UUIDArray
	FilterResults         bool
	Results               []RunResult
	FilterScheduledBefore bool
	ScheduledBefore       time.Time
	FilterScheduledAfter  bool
	ScheduledAfter        time.Time
	FilterStartedBefore   bool
	StartedBefore         time.Time
	FilterStartedAfter    bool
	StartedAfter          time.Time
	FilterFinishedBefore  bool
	FinishedBefore        time.Time
	FilterFinishedAfter   bool
	FinishedAfter         time.Time
}

type ListRunsRow struct {
	ID              string       `json:"id"`
	TestID          string       `json:"test_id"`
	TestRunConfigID int          `json:"test_run_config_id"`
	RunnerID        string       `json:"runner_id"`
	Result          RunResult    `json:"result"`
	Logs            pgtype.JSONB `json:"logs"`
	ScheduledAt     time.Time    `json:"scheduled_at"`
	StartedAt       time.Time    `json:"started_at"`
	FinishedAt      time.Time    `json:"finished_at"`
	ContainerImage  string       `json:"container_image"`
	Command         string       `json:"command"`
	Args            []string     `json:"args"`
	Env             pgtype.JSONB `json:"env"`
	CreatedAt       time.Time    `json:"created_at"`
}

// ListRuns implements Querier.ListRuns.
func (q *DBQuerier) ListRuns(ctx context.Context, params ListRunsParams) ([]ListRunsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "ListRuns")
	rows, err := q.conn.Query(ctx, listRunsSQL, params.FilterTestIds, params.TestIds, params.FilterTestSuiteIds, params.TestSuiteIds, params.FilterRunnerIds, params.RunnerIds, params.FilterResults, q.types.newRunResultArrayInit(params.Results), params.FilterScheduledBefore, params.ScheduledBefore, params.FilterScheduledAfter, params.ScheduledAfter, params.FilterStartedBefore, params.StartedBefore, params.FilterStartedAfter, params.StartedAfter, params.FilterFinishedBefore, params.FinishedBefore, params.FilterFinishedAfter, params.FinishedAfter)
	if err != nil {
		return nil, fmt.Errorf("query ListRuns: %w", err)
	}
	defer rows.Close()
	items := []ListRunsRow{}
	for rows.Next() {
		var item ListRunsRow
		if err := rows.Scan(&item.ID, &item.TestID, &item.TestRunConfigID, &item.RunnerID, &item.Result, &item.Logs, &item.ScheduledAt, &item.StartedAt, &item.FinishedAt, &item.ContainerImage, &item.Command, &item.Args, &item.Env, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan ListRuns row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close ListRuns rows: %w", err)
	}
	return items, err
}

// ListRunsBatch implements Querier.ListRunsBatch.
func (q *DBQuerier) ListRunsBatch(batch genericBatch, params ListRunsParams) {
	batch.Queue(listRunsSQL, params.FilterTestIds, params.TestIds, params.FilterTestSuiteIds, params.TestSuiteIds, params.FilterRunnerIds, params.RunnerIds, params.FilterResults, q.types.newRunResultArrayInit(params.Results), params.FilterScheduledBefore, params.ScheduledBefore, params.FilterScheduledAfter, params.ScheduledAfter, params.FilterStartedBefore, params.StartedBefore, params.FilterStartedAfter, params.StartedAfter, params.FilterFinishedBefore, params.FinishedBefore, params.FilterFinishedAfter, params.FinishedAfter)
}

// ListRunsScan implements Querier.ListRunsScan.
func (q *DBQuerier) ListRunsScan(results pgx.BatchResults) ([]ListRunsRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query ListRunsBatch: %w", err)
	}
	defer rows.Close()
	items := []ListRunsRow{}
	for rows.Next() {
		var item ListRunsRow
		if err := rows.Scan(&item.ID, &item.TestID, &item.TestRunConfigID, &item.RunnerID, &item.Result, &item.Logs, &item.ScheduledAt, &item.StartedAt, &item.FinishedAt, &item.ContainerImage, &item.Command, &item.Args, &item.Env, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan ListRunsBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close ListRunsBatch rows: %w", err)
	}
	return items, err
}

const scheduleRunSQL = `INSERT INTO runs (test_id, test_run_config_id)
VALUES ($1::uuid, $2::integer)
RETURNING *;`

type ScheduleRunRow struct {
	ID              string       `json:"id"`
	TestID          string       `json:"test_id"`
	TestRunConfigID int          `json:"test_run_config_id"`
	RunnerID        string       `json:"runner_id"`
	Result          RunResult    `json:"result"`
	Logs            pgtype.JSONB `json:"logs"`
	ScheduledAt     time.Time    `json:"scheduled_at"`
	StartedAt       time.Time    `json:"started_at"`
	FinishedAt      time.Time    `json:"finished_at"`
}

// ScheduleRun implements Querier.ScheduleRun.
func (q *DBQuerier) ScheduleRun(ctx context.Context, testID string, testRunConfigID int) (ScheduleRunRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "ScheduleRun")
	row := q.conn.QueryRow(ctx, scheduleRunSQL, testID, testRunConfigID)
	var item ScheduleRunRow
	if err := row.Scan(&item.ID, &item.TestID, &item.TestRunConfigID, &item.RunnerID, &item.Result, &item.Logs, &item.ScheduledAt, &item.StartedAt, &item.FinishedAt); err != nil {
		return item, fmt.Errorf("query ScheduleRun: %w", err)
	}
	return item, nil
}

// ScheduleRunBatch implements Querier.ScheduleRunBatch.
func (q *DBQuerier) ScheduleRunBatch(batch genericBatch, testID string, testRunConfigID int) {
	batch.Queue(scheduleRunSQL, testID, testRunConfigID)
}

// ScheduleRunScan implements Querier.ScheduleRunScan.
func (q *DBQuerier) ScheduleRunScan(results pgx.BatchResults) (ScheduleRunRow, error) {
	row := results.QueryRow()
	var item ScheduleRunRow
	if err := row.Scan(&item.ID, &item.TestID, &item.TestRunConfigID, &item.RunnerID, &item.Result, &item.Logs, &item.ScheduledAt, &item.StartedAt, &item.FinishedAt); err != nil {
		return item, fmt.Errorf("scan ScheduleRunBatch row: %w", err)
	}
	return item, nil
}

const nextRunSQL = `UPDATE runs
SET runner_id = $1, scheduled_at = CURRENT_TIMESTAMP
WHERE id = (
  SELECT runs.id
  FROM runs
  JOIN tests
  ON runs.test_id = tests.id
  WHERE scheduled_at IS NULL AND tests.labels @> $2
  ORDER BY scheduled_at ASC
  LIMIT 1
  FOR UPDATE
)
RETURNING *;`

type NextRunRow struct {
	ID              string       `json:"id"`
	TestID          string       `json:"test_id"`
	TestRunConfigID int          `json:"test_run_config_id"`
	RunnerID        string       `json:"runner_id"`
	Result          RunResult    `json:"result"`
	Logs            pgtype.JSONB `json:"logs"`
	ScheduledAt     time.Time    `json:"scheduled_at"`
	StartedAt       time.Time    `json:"started_at"`
	FinishedAt      time.Time    `json:"finished_at"`
}

// NextRun implements Querier.NextRun.
func (q *DBQuerier) NextRun(ctx context.Context, runnerID string, labels pgtype.JSONB) (NextRunRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "NextRun")
	row := q.conn.QueryRow(ctx, nextRunSQL, runnerID, labels)
	var item NextRunRow
	if err := row.Scan(&item.ID, &item.TestID, &item.TestRunConfigID, &item.RunnerID, &item.Result, &item.Logs, &item.ScheduledAt, &item.StartedAt, &item.FinishedAt); err != nil {
		return item, fmt.Errorf("query NextRun: %w", err)
	}
	return item, nil
}

// NextRunBatch implements Querier.NextRunBatch.
func (q *DBQuerier) NextRunBatch(batch genericBatch, runnerID string, labels pgtype.JSONB) {
	batch.Queue(nextRunSQL, runnerID, labels)
}

// NextRunScan implements Querier.NextRunScan.
func (q *DBQuerier) NextRunScan(results pgx.BatchResults) (NextRunRow, error) {
	row := results.QueryRow()
	var item NextRunRow
	if err := row.Scan(&item.ID, &item.TestID, &item.TestRunConfigID, &item.RunnerID, &item.Result, &item.Logs, &item.ScheduledAt, &item.StartedAt, &item.FinishedAt); err != nil {
		return item, fmt.Errorf("scan NextRunBatch row: %w", err)
	}
	return item, nil
}

const updateRunSQL = `UPDATE runs
SET
  result = $1,
  logs = $2,
  started_at = $3::timestamptz,
  finished_at = $4::timestamptz
WHERE id = $5::uuid;`

type UpdateRunParams struct {
	Result     RunResult
	Logs       pgtype.JSONB
	StartedAt  time.Time
	FinishedAt time.Time
	ID         string
}

// UpdateRun implements Querier.UpdateRun.
func (q *DBQuerier) UpdateRun(ctx context.Context, params UpdateRunParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateRun")
	cmdTag, err := q.conn.Exec(ctx, updateRunSQL, params.Result, params.Logs, params.StartedAt, params.FinishedAt, params.ID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query UpdateRun: %w", err)
	}
	return cmdTag, err
}

// UpdateRunBatch implements Querier.UpdateRunBatch.
func (q *DBQuerier) UpdateRunBatch(batch genericBatch, params UpdateRunParams) {
	batch.Queue(updateRunSQL, params.Result, params.Logs, params.StartedAt, params.FinishedAt, params.ID)
}

// UpdateRunScan implements Querier.UpdateRunScan.
func (q *DBQuerier) UpdateRunScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec UpdateRunBatch: %w", err)
	}
	return cmdTag, err
}
