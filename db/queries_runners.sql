-- name: RegisterRunner :one
INSERT INTO runners (name, accept_test_labels, reject_test_labels)
VALUES (pggen.arg('name'), pggen.arg('accept_test_labels'), pggen.arg('reject_test_labels'))
RETURNING *;

-- name: GetRunner :one
SELECT *
FROM runners
WHERE id = pggen.arg('id');

-- name: ListRunners :many
SELECT *
FROM runners
WHERE
  CASE WHEN pggen.arg('filter_revoked')
    THEN revoked_at IS NOT NULL
    ELSE TRUE
  END;

-- name: ApproveRunner :exec
UPDATE runners
SET approved_at = CURRENT_TIMESTAMP
WHERE id = pggen.arg('id')::uuid;

-- name: RevokeRunner :exec
UPDATE runners
SET revoked_at = CURRENT_TIMESTAMP
WHERE id = pggen.arg('id')::uuid;
