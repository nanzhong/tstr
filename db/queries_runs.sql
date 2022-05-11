-- name: GetRun :one
SELECT runs.*, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at
FROM runs
JOIN test_run_configs
ON runs.test_run_config_id = test_run_configs.id
WHERE runs.id = pggen.arg('id');

-- name: ListRuns :many
SELECT runs.*, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at
FROM runs
JOIN test_run_configs
ON runs.test_run_config_id = test_run_configs.id
WHERE
  CASE WHEN pggen.arg('filter_test_ids')
    THEN runs.test_id = ANY (pggen.arg('test_ids')::uuid[])
    ELSE TRUE
  END AND
  CASE WHEN pggen.arg('filter_test_suite_ids')
    THEN runs.test_id = ANY (
      SELECT tests.id
      FROM test_suites
      JOIN tests
      ON tests.labels @> test_suites.labels
      WHERE test_suites.id = ANY (pggen.arg('test_suite_ids')::uuid[])
    )
    ELSE TRUE
  END AND
  CASE WHEN pggen.arg('filter_runner_ids')
    THEN runner_id = ANY (pggen.arg('runner_ids')::uuid[])
    ELSE TRUE
  END AND
  CASE WHEN pggen.arg('filter_results')
    THEN result = ANY (pggen.arg('results')::run_result[])
    ELSE TRUE
  END AND
  CASE WHEN pggen.arg('filter_scheduled_before')
    THEN scheduled_at < pggen.arg('scheduled_before')::timestamptz
    ELSE TRUE
  END AND
  CASE WHEN pggen.arg('filter_scheduled_after')
    THEN scheduled_at > pggen.arg('scheduled_after')::timestamptz
    ELSE TRUE
  END AND
  CASE WHEN pggen.arg('filter_started_before')
    THEN started_at < pggen.arg('started_before')::timestamptz
    ELSE TRUE
  END AND
  CASE WHEN pggen.arg('filter_started_after')
    THEN started_at > pggen.arg('started_after')::timestamptz
    ELSE TRUE
  END AND
  CASE WHEN pggen.arg('filter_finished_before')
    THEN finished_at < pggen.arg('finished_before')::timestamptz
    ELSE TRUE
  END AND
  CASE WHEN pggen.arg('filter_finished_after')
    THEN finished_at > pggen.arg('finished_after')::timestamptz
    ELSE TRUE
  END;

-- name: ScheduleRun :one
WITH scheduled_run AS (
  INSERT INTO runs (test_id, test_run_config_id)
  VALUES (pggen.arg('test_id')::uuid, pggen.arg('test_run_config_id')::uuid)
  RETURNING *
)
SELECT scheduled_run.*, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at AS test_run_config_created_at
FROM scheduled_run
JOIN test_run_configs
ON scheduled_run.test_run_config_id = test_run_configs.id;

-- name: AssignRun :one
UPDATE runs
SET runner_id = pggen.arg('runner_id')
FROM test_run_configs
WHERE runs.id = (
  SELECT id
  FROM runs
  WHERE test_id = ANY(pggen.arg('test_ids')) AND runner_id IS NULL
  ORDER BY scheduled_at ASC
  LIMIT 1
) AND runs.test_run_config_id = test_run_configs.id
RETURNING runs.*, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at AS test_run_config_created_at;

-- name: UpdateRun :exec
UPDATE runs
SET
  result = pggen.arg('result'),
  logs = pggen.arg('logs'),
  started_at = pggen.arg('started_at')::timestamptz,
  finished_at = pggen.arg('finished_at')::timestamptz
WHERE id = pggen.arg('id')::uuid;
