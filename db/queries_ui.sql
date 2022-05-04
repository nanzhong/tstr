-- name: UIListRecentRuns :many
SELECT runs.*, tests.name as test_name, tests.labels, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at, runners.name AS runner_name
FROM runs
JOIN test_run_configs ON runs.test_run_config_id = test_run_configs.id
JOIN runners ON runs.runner_id = runners.id
JOIN tests on runs.test_id = tests.id
ORDER BY runs.started_at DESC
LIMIT pggen.arg('limit');

-- name: UITestsByLabels :many
SELECT labels, array_agg(tests.*) AS tests FROM tests WHERE labels IN (SELECT DISTINCT(labels) FROM tests) GROUP BY labels;

-- name: UITestResults :many
SELECT test_id,array_agg(result) AS results FROM runs where test_id is not null GROUP BY test_id;
