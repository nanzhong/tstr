-- name: GetRun :one
SELECT runs.*, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at
FROM runs
JOIN test_run_configs
ON runs.test_run_config_id = test_run_configs.id
WHERE runs.id = sqlc.arg('id');

-- name: ListRuns :many
SELECT runs.*, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at
FROM runs
JOIN test_run_configs
ON runs.test_run_config_id = test_run_configs.id
WHERE
  (sqlc.narg('test_ids')::uuid[] IS NULL OR runs.test_id = ANY (sqlc.narg('test_ids')::uuid[])) AND
  (sqlc.narg('test_suite_ids')::uuid[] IS NULL OR runs.test_id = ANY (
      SELECT tests.id
      FROM test_suites
      JOIN tests
      ON tests.labels @> test_suites.labels
      WHERE test_suites.id = ANY (sqlc.narg('test_suite_ids')::uuid[])
    )) AND
  (sqlc.narg('runner_ids')::uuid[] IS NULL OR runner_id = ANY (sqlc.narg('runner_ids')::uuid[])) AND
  (sqlc.narg('results')::run_result[] IS NULL OR result = ANY (sqlc.narg('results')::run_result[])) AND
  (sqlc.arg('scheduled_before')::timestamptz IS NULL OR scheduled_at < sqlc.narg('scheduled_before')::timestamptz) AND
  (sqlc.narg('scheduled_after')::timestamptz IS NULL OR scheduled_at > sqlc.narg('scheduled_after')::timestamptz) AND
  (sqlc.narg('started_before')::timestamptz IS NULL OR started_at < sqlc.narg('started_before')::timestamptz) AND
  (sqlc.narg('started_after')::timestamptz IS NULL OR started_at > sqlc.narg('started_after')::timestamptz) AND
  (sqlc.narg('finished_before')::timestamptz IS NULL OR finished_at < sqlc.narg('finished_before')::timestamptz) AND
  (sqlc.narg('finished_after')::timestamptz IS NULL OR finished_at > sqlc.narg('finished_after')::timestamptz);

-- name: ScheduleRun :one
WITH scheduled_run AS (
  INSERT INTO runs (test_id, test_run_config_id)
  VALUES (sqlc.arg('test_id')::uuid, sqlc.arg('test_run_config_id')::uuid)
  RETURNING *
)
SELECT scheduled_run.*, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at AS test_run_config_created_at
FROM scheduled_run
JOIN test_run_configs
ON scheduled_run.test_run_config_id = test_run_configs.id;

-- name: AssignRun :one
UPDATE runs
SET runner_id = sqlc.arg('runner_id')::uuid
FROM test_run_configs
WHERE runs.id = (
  SELECT id
  FROM runs AS selected_runs
  WHERE selected_runs.test_id = ANY(sqlc.arg('test_ids')::uuid[]) AND selected_runs.runner_id IS NULL
  ORDER BY selected_runs.scheduled_at ASC
  LIMIT 1
) AND runs.test_run_config_id = test_run_configs.id
RETURNING runs.*, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at AS test_run_config_created_at;

-- name: UpdateRun :exec
UPDATE runs
SET
  result = sqlc.arg('result'),
  logs = sqlc.arg('logs'),
  started_at = sqlc.arg('started_at')::timestamptz,
  finished_at = sqlc.arg('finished_at')::timestamptz
WHERE id = sqlc.arg('id')::uuid;
