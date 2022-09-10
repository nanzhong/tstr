-- name: GetRun :one
SELECT *
FROM runs
WHERE runs.id = sqlc.arg('id');

-- name: ListRuns :many
SELECT *
FROM runs;

-- name: ScheduleRun :one
INSERT INTO runs (test_id, test_run_config, labels, test_matrix_id)
SELECT tests.id, tests.run_config, sqlc.arg('labels'), sqlc.narg('test_matrix_id')
FROM tests
WHERE tests.id = sqlc.arg('test_id')
RETURNING *;

-- name: ListPendingRuns :many
SELECT runs.*, tests.namespace
FROM runs
JOIN tests
ON runs.test_id = tests.id
WHERE runner_id IS NULL;

-- name: AssignRun :one
UPDATE runs
SET runner_id = sqlc.arg('runner_id')::uuid
WHERE runs.id = (
  SELECT matching_runs.id
  FROM runs AS matching_runs
  WHERE
    matching_runs.id = ANY (sqlc.arg('run_IDs')::uuid[]) AND
    matching_runs.runner_id IS NULL
  LIMIT 1
  FOR UPDATE SKIP LOCKED
)
RETURNING *;

-- name: UpdateRun :exec
UPDATE runs
SET
  result = sqlc.arg('result'),
  started_at = sqlc.narg('started_at')::timestamptz,
  finished_at = sqlc.narg('finished_at')::timestamptz
WHERE id = sqlc.arg('id');

-- name: DeleteRunsForTest :exec
DELETE FROM runs
WHERE test_id = sqlc.arg('test_id');

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

-- name: TimeoutRuns :exec
UPDATE runs
SET
  result = 'error',
  finished_at = CURRENT_TIMESTAMP,
  logs = COALESCE(logs, '[]'::jsonb) || sqlc.arg('timeout_log')
WHERE
  result = 'unknown' AND
  runner_id IS NOT NULL AND
  CURRENT_TIMESTAMP > started_at + make_interval(secs => COALESCE(test_run_config['timeout_seconds']::int, sqlc.arg('default_timeout')::int));

-- name: RunSummariesForTest :many
SELECT runs.id, tests.id AS test_id, tests.name AS test_name, runs.test_run_config, runs.labels, runs.runner_id, runs.result, runs.scheduled_at, runs.started_at, runs.finished_at, runs.result_data
FROM runs
JOIN tests
ON runs.test_id = tests.id
WHERE runs.test_id = sqlc.arg('test_id') AND runs.scheduled_at > sqlc.arg('scheduled_after')
ORDER by runs.scheduled_at desc;

-- name: RunSummariesForRunner :many
SELECT runs.id, tests.id AS test_id, tests.name AS test_name, runs.test_run_config, runs.labels, runs.runner_id, runs.result, runs.scheduled_at, runs.started_at, runs.finished_at, runs.result_data
FROM runs
JOIN tests
ON runs.test_id = tests.id
WHERE runs.runner_id = sqlc.arg('runner_id') AND runs.scheduled_at > sqlc.arg('scheduled_after')
ORDER by runs.scheduled_at desc;

-- name: QueryRuns :many
SELECT *
FROM runs
WHERE
  (sqlc.narg('ids')::uuid[] IS NULL OR runs.id = ANY (sqlc.narg('ids')::uuid[])) AND
  (sqlc.narg('test_ids')::uuid[] IS NULL OR runs.test_id = ANY (sqlc.narg('test_ids')::uuid[])) AND
  (sqlc.narg('runner_ids')::uuid[] IS NULL OR runner_id = ANY (sqlc.narg('runner_ids')::uuid[])) AND
  (sqlc.narg('results')::text[] IS NULL OR result::text = ANY (sqlc.narg('results')::text[])) AND
  (sqlc.arg('scheduled_before')::timestamptz IS NULL OR scheduled_at < sqlc.narg('scheduled_before')::timestamptz) AND
  (sqlc.narg('scheduled_after')::timestamptz IS NULL OR scheduled_at > sqlc.narg('scheduled_after')::timestamptz) AND
  (sqlc.narg('started_before')::timestamptz IS NULL OR started_at < sqlc.narg('started_before')::timestamptz) AND
  (sqlc.narg('started_after')::timestamptz IS NULL OR started_at > sqlc.narg('started_after')::timestamptz) AND
  (sqlc.narg('finished_before')::timestamptz IS NULL OR finished_at < sqlc.narg('finished_before')::timestamptz) AND
  (sqlc.narg('finished_after')::timestamptz IS NULL OR finished_at > sqlc.narg('finished_after')::timestamptz)
ORDER BY scheduled_at DESC;

-- name: SummarizeRunsBreakdownResult :many
WITH intervals AS (
  SELECT generate_series(
    date_trunc(sqlc.arg('precision'), sqlc.arg('start_time')::timestamptz) + make_interval(secs => sqlc.arg('interval')),
    date_trunc(sqlc.arg('precision'), sqlc.arg('end_time')::timestamptz),
    make_interval(secs => sqlc.arg('interval'))
  ) as start
)
SELECT 
  intervals.start::timestamptz,
  COUNT(id) FILTER (WHERE result = 'pass') as pass,
  COUNT(id) FILTER (WHERE result = 'fail') as fail,
  COUNT(id) FILTER (WHERE result = 'error') as error,
  COUNT(id) FILTER (WHERE result = 'unknown') as unknown
FROM intervals
LEFT JOIN runs
ON
  intervals.start = date_trunc(sqlc.arg('precision'), runs.scheduled_at) AND
  runs.scheduled_at > sqlc.arg('start_time') AND
  runs.scheduled_at < sqlc.arg('end_time')
GROUP BY intervals.start
ORDER BY intervals.start ASC;

-- name: SummarizeRunsBreakdownTest :many
WITH intervals AS (
  SELECT generate_series(
    date_trunc(sqlc.arg('precision'), sqlc.arg('start_time')::timestamptz) + make_interval(secs => sqlc.arg('interval')),
    date_trunc(sqlc.arg('precision'), sqlc.arg('end_time')::timestamptz),
    make_interval(secs => sqlc.arg('interval'))
  ) as start
)
SELECT 
  intervals.start::timestamptz,
  tests.id,
  tests.name,
  COUNT(tests.id) FILTER (WHERE runs.result = 'pass') as pass,
  COUNT(tests.id) FILTER (WHERE runs.result = 'fail') as fail,
  COUNT(tests.id) FILTER (WHERE runs.result = 'error') as error,
  COUNT(tests.id) FILTER (WHERE runs.result = 'unknown') as unknown
FROM intervals
LEFT JOIN runs
ON
  intervals.start = date_trunc(sqlc.arg('precision'), runs.scheduled_at) AND
  runs.scheduled_at > sqlc.arg('start_time') AND
  runs.scheduled_at < sqlc.arg('end_time')
JOIN tests
ON runs.test_id = tests.id
GROUP BY intervals.start, tests.id
ORDER BY intervals.start ASC;
