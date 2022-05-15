-- name: RegisterTest :one
WITH data (name, labels, cron_schedule, next_run_at, container_image, command, args, env) AS (
  VALUES (
    sqlc.arg('name')::varchar,
    sqlc.arg('labels')::jsonb,
    sqlc.arg('cron_schedule')::varchar,
    sqlc.narg('next_run_at')::timestamptz,
    sqlc.arg('container_image')::varchar,
    sqlc.arg('command')::varchar,
    sqlc.arg('args')::varchar[],
    sqlc.arg('env')::jsonb
  )
), test AS (
  INSERT INTO tests (name, labels, cron_schedule, next_run_at)
  SELECT name, labels, cron_schedule, next_run_at
  FROM data
  RETURNING id, name, labels, cron_schedule, next_run_at, registered_at, updated_at
), test_run_config AS (
  INSERT INTO test_run_configs (test_id, container_image, command, args, env)
  SELECT test.id, container_image, command, args, env
  FROM data, test
  RETURNING id AS test_run_config_id, container_image, command, args, env, created_at AS test_run_config_created_at
)
SELECT * FROM test, test_run_config;

-- name: GetTest :one
SELECT tests.*, latest_configs.id AS test_run_config_id, latest_configs.container_image, latest_configs.command, latest_configs.args, latest_configs.env, latest_configs.created_at
FROM tests
JOIN test_run_configs AS latest_configs
ON tests.id = latest_configs.test_id
LEFT JOIN test_run_configs
ON test_run_configs.test_id = latest_configs.test_id AND latest_configs.created_at > test_run_configs.created_at
WHERE tests.id = sqlc.arg('id')::uuid
ORDER BY test_run_configs.created_at DESC
LIMIT 1;

-- name: ListTests :many
SELECT tests.*, latest_configs.id AS test_run_config_id, latest_configs.container_image, latest_configs.command, latest_configs.args, latest_configs.env, latest_configs.created_at
FROM tests
JOIN test_run_configs AS latest_configs
ON tests.id = latest_configs.test_id
LEFT JOIN test_run_configs
ON test_run_configs.test_id = latest_configs.test_id AND latest_configs.created_at > test_run_configs.created_at
WHERE test_run_configs IS NULL AND tests.labels @> sqlc.arg('labels')::jsonb
ORDER BY tests.name ASC;

-- name: ListTestsIDsMatchingLabelKeys :many
SELECT tests.id, tests.labels
FROM tests
WHERE
  tests.labels ?& sqlc.arg('include_label_keys')::varchar[] AND
  NOT tests.labels ?| COALESCE(sqlc.arg('filter_label_keys')::varchar[], '{}')::varchar[];

-- name: UpdateTest :exec
UPDATE tests
SET
  name = sqlc.arg('name')::varchar,
  labels = sqlc.arg('labels')::jsonb,
  cron_schedule = sqlc.arg('cron_schedule')::varchar,
  next_run_at = sqlc.narg('next_run_at')::timestamptz,
  updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')::uuid;

-- name: CreateTestRunConfig :one
INSERT INTO test_run_configs (container_image, command, args, env)
VALUES (
  sqlc.arg('container_image')::varchar,
  sqlc.arg('command')::varchar,
  sqlc.arg('args')::varchar[],
  sqlc.arg('env')::jsonb
)
RETURNING *;

-- name: ArchiveTest :exec
UPDATE tests
SET archived_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')::uuid;

-- name: ListTestsToSchedule :many
SELECT tests.*, latest_configs.id AS test_run_config_id, latest_configs.container_image, latest_configs.command, latest_configs.args, latest_configs.env, latest_configs.created_at
FROM tests
JOIN test_run_configs AS latest_configs
ON tests.id = latest_configs.test_id
LEFT JOIN test_run_configs
ON test_run_configs.test_id = latest_configs.test_id AND latest_configs.created_at > test_run_configs.created_at
LEFT JOIN runs
ON runs.test_id = tests.id AND runs.result = 'unknown' AND runs.started_at IS NULL
WHERE next_run_at < CURRENT_TIMESTAMP AND runs.id IS NULL
FOR UPDATE OF tests SKIP LOCKED;
