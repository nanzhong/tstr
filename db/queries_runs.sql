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
INSERT INTO runs (test_id, test_run_config_id)
VALUES (pggen.arg('test_id')::uuid, pggen.arg('test_run_config_id')::integer)
RETURNING *;
