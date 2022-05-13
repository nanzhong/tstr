-- name: DefineTestSuite :one
INSERT INTO test_suites (name, labels)
VALUES (sqlc.arg('name'), sqlc.arg('labels'))
RETURNING *;

-- name: UpdateTestSuite :exec
UPDATE test_suites
SET
  name = sqlc.arg('name')::varchar,
  labels = sqlc.arg('labels'),
  updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id');

-- name: GetTestSuite :one
SELECT *
FROM test_suites
WHERE id = sqlc.arg('id')::uuid;

-- name: ListTestSuites :many
SELECT *
FROM test_suites
WHERE labels @> sqlc.arg('labels')
ORDER BY name ASC;

-- name: ArchiveTestSuite :exec
UPDATE test_suites
SET archived_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id');
