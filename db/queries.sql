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
