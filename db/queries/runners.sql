-- name: RegisterRunner :one
INSERT INTO runners (name, accept_test_label_selectors, reject_test_label_selectors, last_heartbeat_at)
VALUES (
  sqlc.arg('name'),
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
WHERE last_heartbeat_at > sqlc.arg('heartbeat_since');