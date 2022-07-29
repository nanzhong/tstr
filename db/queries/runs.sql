-- name: GetRun :one
SELECT *
FROM runs
WHERE runs.id = sqlc.arg('id');

-- name: ListRuns :many
SELECT *
FROM runs;

-- name: ScheduleRun :one
WITH latest_run_config AS (
  SELECT run_config
  FROM tests
  WHERE tests.id = sqlc.arg('test_id')
)
INSERT INTO runs (test_id, test_run_config)
VALUES (sqlc.arg('test_id'), latest_run_config)
RETURNING *;

-- name: AssignRun :one
UPDATE runs
SET runner_id = sqlc.arg('runner_id')::uuid
WHERE runs.id = (
  SELECT id
  FROM runs AS selected_runs
  WHERE selected_runs.test_id = ANY(sqlc.arg('test_ids')::uuid[]) AND selected_runs.runner_id IS NULL
  ORDER BY selected_runs.scheduled_at ASC
  LIMIT 1
)
RETURNING runs.*;

-- name: UpdateRun :exec
UPDATE runs
SET
  result = sqlc.arg('result'),
  started_at = sqlc.narg('started_at')::timestamptz,
  finished_at = sqlc.narg('finished_at')::timestamptz
WHERE id = sqlc.arg('id');

-- name: AppendLogsToRun :exec
UPDATE runs
SET logs = COALESCE(logs, '[]'::jsonb) || sqlc.arg('logs')
WHERE id = sqlc.arg('id');

-- name: UpdateResultData :exec
UPDATE runs
SET result_data = result_data || sqlc.arg('result_data')::jsonb 
WHERE id = sqlc.arg('id');

-- name: ResetOrphanedRuns :exec
UPDATE runs
SET runner_id = NULL
WHERE
  result = 'unknown' AND
  started_at IS NULL AND
  scheduled_at < sqlc.arg('before')::timestamptz;

-- name: RunSummaryForTest :many
SELECT id, test_id, test_run_config, runner_id, result, scheduled_at, started_at, finished_at, result_data
FROM runs
WHERE runs.test_id = sqlc.arg('test_id')
ORDER by runs.scheduled_at desc
LIMIT sqlc.arg('limit');

-- name: RunSummaryForRunner :many
SELECT id, test_id, test_run_config, runner_id, result, scheduled_at, started_at, finished_at, result_data
FROM runs
WHERE runs.runner_id = sqlc.arg('runner_id')::uuid
ORDER by runs.scheduled_at desc
LIMIT sqlc.arg('limit');

-- name: QueryRuns :many
SELECT *
FROM runs
WHERE
  (sqlc.narg('ids')::uuid[] IS NULL OR runs.id = ANY (sqlc.narg('ids')::uuid[])) AND
  (sqlc.narg('test_ids')::uuid[] IS NULL OR runs.test_id = ANY (sqlc.narg('test_ids')::uuid[])) AND
  (sqlc.narg('test_suite_ids')::uuid[] IS NULL OR runs.test_id = ANY (
      SELECT tests.id
      FROM test_suites
      JOIN tests
      ON tests.labels @> test_suites.labels
      WHERE test_suites.id = ANY (sqlc.narg('test_suite_ids')::uuid[])
    )) AND
  (sqlc.narg('runner_ids')::uuid[] IS NULL OR runner_id = ANY (sqlc.narg('runner_ids')::uuid[])) AND
  (sqlc.narg('results')::text[] IS NULL OR result::text = ANY (sqlc.narg('results')::text[])) AND
  (sqlc.arg('scheduled_before')::timestamptz IS NULL OR scheduled_at < sqlc.narg('scheduled_before')::timestamptz) AND
  (sqlc.narg('scheduled_after')::timestamptz IS NULL OR scheduled_at > sqlc.narg('scheduled_after')::timestamptz) AND
  (sqlc.narg('started_before')::timestamptz IS NULL OR started_at < sqlc.narg('started_before')::timestamptz) AND
  (sqlc.narg('started_after')::timestamptz IS NULL OR started_at > sqlc.narg('started_after')::timestamptz) AND
  (sqlc.narg('finished_before')::timestamptz IS NULL OR finished_at < sqlc.narg('finished_before')::timestamptz) AND
  (sqlc.narg('finished_after')::timestamptz IS NULL OR finished_at > sqlc.narg('finished_after')::timestamptz);


