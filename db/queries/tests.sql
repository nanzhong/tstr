-- name: RegisterTest :one
INSERT INTO tests (name, labels, run_config, cron_schedule, next_run_at)
VALUES (
  sqlc.arg('name')::varchar,
  sqlc.arg('labels')::jsonb,
  sqlc.arg('run_config')::jsonb,
  sqlc.arg('cron_schedule')::varchar,
  sqlc.narg('next_run_at')::timestamptz
)
RETURNING *;

-- name: GetTest :one
SELECT *
FROM tests
WHERE tests.id = sqlc.arg('id');

-- name: ListTests :many
SELECT *
FROM tests
WHERE archived_at IS NULL
ORDER BY tests.name ASC;

-- name: ListTestsIDsMatchingLabelKeys :many
SELECT tests.id, tests.labels
FROM tests
WHERE
  tests.labels ?& COALESCE(sqlc.arg('include_label_keys')::varchar[], '{}') AND
  NOT tests.labels ?| COALESCE(sqlc.arg('filter_label_keys')::varchar[], '{}')::varchar[];

-- name: UpdateTest :exec
UPDATE tests
SET
  name = sqlc.arg('name')::varchar,
  labels = sqlc.arg('labels')::jsonb,
  run_config = sqlc.arg('run_config')::jsonb,
  cron_schedule = sqlc.arg('cron_schedule')::varchar,
  next_run_at = sqlc.narg('next_run_at')::timestamptz,
  updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')::uuid;

-- name: ArchiveTest :exec
UPDATE tests
SET archived_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')::uuid;

-- name: ListTestsToSchedule :many
SELECT *
FROM tests
LEFT JOIN runs
ON runs.test_id = tests.id AND runs.result = 'unknown' AND runs.started_at IS NULL
WHERE tests.next_run_at < CURRENT_TIMESTAMP AND runs.id IS NULL
FOR UPDATE OF tests SKIP LOCKED;

-- name: QueryTests :many
SELECT *
FROM tests
WHERE
  (sqlc.narg('ids')::uuid[] IS NULL OR tests.id = ANY (sqlc.narg('ids')::uuid[])) AND
  (sqlc.narg('test_suite_ids')::uuid[] IS NULL OR tests.id = ANY (
    SELECT tests.id
    FROM test_suites
    JOIN tests
    ON tests.labels @> test_suites.labels
    WHERE test_suites.id = ANY (sqlc.narg('test_suite_ids')::uuid[])
    )) AND
  (sqlc.narg('labels')::jsonb IS NULL OR tests.labels @> sqlc.narg('labels')::jsonb)
ORDER BY tests.name ASC;
