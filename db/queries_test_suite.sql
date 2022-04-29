-- name: DefineTestSuite :one
INSERT INTO test_suites (name, labels)
VALUES (pggen.arg('name')::varchar, pggen.arg('labels')::jsonb)
RETURNING *;

-- name: UpdateTestSuite :exec
UPDATE test_suites
SET
  name = pggen.arg('name')::varchar,
  labels = pggen.arg('labels')::jsonb,
  updated_at = CURRENT_TIMESTAMP
WHERE id = pggen.arg('id');

-- name: GetTestSuite :one
SELECT *
FROM test_suites
WHERE id = pggen.arg('id')::uuid;

-- name: listTestSuites :many
SELECT *
FROM test_suites
WHERE labels @> pggen.arg('labels')::jsonb
ORDER BY name ASC;

-- name: ArchiveTestSuite :exec
UPDATE test_suites
SET archived_at = CURRENT_TIMESTAMP
WHERE id = pggen.arg('id');
