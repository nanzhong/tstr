-- name: RegisterRunner :one
INSERT INTO runners (name, namespace_selectors, accept_test_label_selectors, reject_test_label_selectors, last_heartbeat_at)
VALUES (
  sqlc.arg('name'),
  sqlc.arg('namespace_selectors'),
  sqlc.arg('accept_test_label_selectors'),
  sqlc.arg('reject_test_label_selectors'),
  CURRENT_TIMESTAMP
)
RETURNING *;

-- name: GetRunner :one
SELECT *
FROM runners
WHERE id = sqlc.arg('id');

-- name: ListRunners :many
SELECT *
FROM runners
ORDER by last_heartbeat_at DESC, registered_at;

-- name: UpdateRunnerHeartbeat :exec
UPDATE runners
SET last_heartbeat_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id');

-- name: QueryRunners :many
SELECT *
FROM runners
WHERE
  (sqlc.narg('ids')::uuid[] IS NULL OR runners.id = ANY (sqlc.narg('ids')::uuid[])) AND
  (sqlc.narg('last_heartbeat_since')::timestamptz IS NULL) OR last_heartbeat_at > sqlc.narg('last_heartbeat_since')
ORDER by name ASC;
