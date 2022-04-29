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
  ((pggen.arg('filter_test_ids') AND runs.test_id = ANY (pggen.arg('test_ids')::uuid[])) OR TRUE) AND
  ((pggen.arg('filter_test_suite_ids') AND runs.test_id = ANY (
    SELECT tests.id
    FROM test_suites
    JOIN tests
    ON tests.labels @> test_suites.labels
    WHERE test_suites.id = ANY (pggen.arg('test_suite_ids')::uuid[])
  )) OR TRUE) AND
  ((pggen.arg('filter_runner_ids') AND runner_id = ANY (pggen.arg('runner_ids')::uuid[])) OR TRUE) AND
  ((pggen.arg('filter_results') AND result = ANY (pggen.arg('results')::run_result[])) OR TRUE) AND
  ((pggen.arg('filter_scheduled_before') AND scheduled_at < (pggen.arg('scheduled_before')::timestamptz)) OR TRUE) AND
  ((pggen.arg('filter_scheduled_after') AND scheduled_at > (pggen.arg('scheduled_after')::timestamptz)) OR TRUE) AND
  ((pggen.arg('filter_started_before') AND started_at < (pggen.arg('started_before')::timestamptz)) OR TRUE) AND
  ((pggen.arg('filter_started_after') AND started_at > (pggen.arg('started_after')::timestamptz)) OR TRUE) AND
  ((pggen.arg('filter_finished_before') AND finished_at < (pggen.arg('finished_before')::timestamptz)) OR TRUE) AND
  ((pggen.arg('filter_finished_after') AND finished_at > (pggen.arg('finished_after')::timestamptz)) OR TRUE);

-- name: ScheduleRun :one
INSERT INTO runs (test_id, test_run_config_id)
VALUES (pggen.arg('test_id')::uuid, pggen.arg('test_run_config_id')::integer)
RETURNING *;
