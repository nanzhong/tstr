-- name: RegisterRunner :one
INSERT INTO runners (name, accept_test_label_selectors, reject_test_label_selectors, last_heartbeat_at)
VALUES (
  pggen.arg('name'),
  pggen.arg('accept_test_label_selectors'),
  pggen.arg('reject_test_label_selectors'),
  CURRENT_TIMESTAMP
)
RETURNING *;

-- name: GetRunner :one
SELECT *
FROM runners
WHERE id = pggen.arg('id');

-- name: ListRunners :many
SELECT *
FROM runners
WHERE last_heartbeat_at > pggen.arg('heartbeat_since');
