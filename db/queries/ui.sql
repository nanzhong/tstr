-- name: UIListRecentRuns :many
SELECT runs.*, tests.name as test_name, tests.labels, test_run_configs.container_image, test_run_configs.command, test_run_configs.args, test_run_configs.env, test_run_configs.created_at, runners.name AS runner_name, (finished_at is NULL)::bool AS is_pending
FROM runs
JOIN test_run_configs ON runs.test_run_config_id = test_run_configs.id
JOIN runners ON runs.runner_id = runners.id
JOIN tests on runs.test_id = tests.id
ORDER BY is_pending, runs.started_at DESC
LIMIT sqlc.arg('limit');

-- name: UITestsByLabels :many
SELECT labels, array_agg(tests.id) AS tests FROM tests WHERE labels IN (SELECT DISTINCT(labels) FROM tests) GROUP BY labels;

-- name: UITestResults :many
SELECT test_id,array_agg(result) AS results FROM runs where test_id is not null GROUP BY test_id;

-- name: UIListTests :many
select * from tests;
