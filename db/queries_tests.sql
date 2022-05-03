-- name: RegisterTest :one
WITH data (name, labels, cron_schedule, container_image, command, args, env) AS (
  VALUES (
    pggen.arg('name')::varchar,
    pggen.arg('labels')::jsonb,
    pggen.arg('cron_schedule')::varchar,
    pggen.arg('container_image')::varchar,
    pggen.arg('command')::varchar,
    pggen.arg('args')::varchar[],
    pggen.arg('env')::jsonb
  )
), test AS (
  INSERT INTO tests (name, labels, cron_schedule)
  SELECT name, labels, cron_schedule
  FROM data
  RETURNING id, name, labels, cron_schedule, registered_at, updated_at
), test_run_config AS (
  INSERT INTO test_run_configs (test_id, container_image, command, args, env)
  SELECT test.id, container_image, command, args, env
  FROM data, test
  RETURNING id AS test_run_config_version, container_image, command, args, env, created_at AS test_run_config_created_at
)
SELECT * FROM test, test_run_config;

-- name: GetTest :one
SELECT tests.*, test_run_configs.id AS test_run_config_version, container_image, command, args, env, created_at
FROM tests
JOIN test_run_configs
ON tests.id = test_run_configs.test_id
WHERE tests.id = pggen.arg('id')::uuid
ORDER BY test_run_configs.created_at DESC
LIMIT 1;

-- name: ListTests :many
SELECT tests.*, latest_configs.id AS test_run_config_version, latest_configs.container_image, latest_configs.command, latest_configs.args, latest_configs.env, latest_configs.created_at
FROM tests
JOIN test_run_configs AS latest_configs
ON tests.id = latest_configs.test_id
LEFT JOIN test_run_configs
ON test_run_configs.test_id = latest_configs.test_id AND latest_configs.created_at > test_run_configs.created_at
WHERE test_run_configs IS NULL AND tests.labels @> pggen.arg('labels')::jsonb
ORDER BY tests.name ASC;

-- name: UpdateTest :exec
UPDATE tests
SET
  name = pggen.arg('name')::varchar,
  labels = pggen.arg('labels')::jsonb,
  cron_schedule = pggen.arg('cron_schedule')::varchar,
  updated_at = CURRENT_TIMESTAMP
WHERE id = pggen.arg('id')::uuid;

-- name: CreateTestRunConfig :one
INSERT INTO test_run_configs (container_image, command, args, env)
VALUES (
  pggen.arg('container_image')::varchar,
  pggen.arg('command')::varchar,
  pggen.arg('args')::varchar[],
  pggen.arg('env')::jsonb
)
RETURNING *;

-- name: ArchiveTest :exec
UPDATE tests
SET archived_at = CURRENT_TIMESTAMP
WHERE id = pggen.arg('id')::uuid;
