-- name: GetRunner :one
SELECT *
FROM runners
WHERE id = pggen.arg('id');

-- name: ListRunners :many
SELECT *
FROM runners
WHERE pggen.arg('filter_revoked') OR revoked_at IS NOT NULL;

-- name: ApproveRunner :exec
UPDATE runners
SET approved_at = CURRENT_TIMESTAMP
WHERE id = pggen.arg('id')::uuid;

-- name: RevokeRunner :exec
UPDATE runners
SET revoked_at = CURRENT_TIMESTAMP
WHERE id = pggen.arg('id')::uuid;
