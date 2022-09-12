-- name: RegisterTest :one
INSERT INTO tests (namespace, name, run_config, labels, matrix, cron_schedule, next_run_at)
VALUES (
  sqlc.arg('namespace'),
  sqlc.arg('name'),
  sqlc.arg('run_config'),
  sqlc.arg('labels'),
  sqlc.arg('matrix'),
  sqlc.narg('cron_schedule'),
  sqlc.arg('next_run_at')
)
RETURNING *;

-- name: GetTest :one
SELECT *
FROM tests
WHERE tests.id = sqlc.arg('id') AND tests.namespace = sqlc.arg('namespace');

-- name: ListTests :many
SELECT *
FROM tests
WHERE tests.namespace = sqlc.arg('namespace')
ORDER BY tests.name ASC;

-- name: UpdateTest :exec
UPDATE tests
SET
  name = sqlc.arg('name')::varchar,
  run_config = sqlc.arg('run_config')::jsonb,
  labels = sqlc.arg('labels')::jsonb,
  matrix = sqlc.arg('matrix')::jsonb,
  cron_schedule = sqlc.narg('cron_schedule')::varchar,
  next_run_at = sqlc.narg('next_run_at')::timestamptz,
  updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')::uuid AND tests.namespace = sqlc.arg('namespace');

-- name: DeleteTest :exec
DELETE FROM tests
WHERE id = sqlc.arg('id')::uuid AND tests.namespace = sqlc.arg('namespace');

-- name: ListTestsToSchedule :many
SELECT tests.*
FROM tests
LEFT JOIN runs
ON runs.test_id = tests.id AND runs.result = 'unknown' AND runs.started_at IS NULL
WHERE tests.next_run_at < CURRENT_TIMESTAMP AND runs.id IS NULL
FOR UPDATE OF tests SKIP LOCKED;

-- name: QueryTests :many
SELECT *
FROM tests
WHERE
  tests.namespace = sqlc.arg('namespace') AND
  (sqlc.narg('ids')::uuid[] IS NULL OR tests.id = ANY (sqlc.narg('ids')::uuid[])) AND
  (sqlc.narg('labels')::jsonb IS NULL OR tests.labels @> sqlc.narg('labels')::jsonb)
ORDER BY tests.name ASC;

-- name: ListAllNamespaces :many
SELECT DISTINCT(namespace)
FROM tests;
