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

-- name: UIRunsSummary :many
SELECT id, test_run_config_id, runner_id, result, scheduled_at, started_at,finished_at
FROM runs
WHERE runs.test_id = sqlc.narg('test_id')::uuid 
ORDER by runs.started_at desc
LIMIT sqlc.arg('limit');


-- name: UIRunnerSummary :many
SELECT T.name AS test_name, T.id as test_id, R.id as run_id, R.result, R.finished_at , R.started_at from runs R 
INNER JOIN tests T on T.id = R.test_id 
WHERE runner_id = sqlc.narg('runner_id')::uuid
ORDER BY R.started_at DESC
LIMIT sqlc.arg('limit');

