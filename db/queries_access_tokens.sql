-- name: IssueAccessToken :one
INSERT INTO access_tokens (name, token_hash, scopes, expires_at)
VALUES (pggen.arg('name'), pggen.arg('token_hash'), pggen.arg('scopes'), pggen.arg('expires_at'))
RETURNING id, name, scopes, issued_at, expires_at;

-- name: ValidateAccessToken :one
SELECT TRUE
FROM access_tokens
WHERE token_hash = pggen.arg('token_hash') AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP);

-- name: GetAccessToken :one
SELECT id, name, scopes, issued_at, expires_at, revoked_at
FROM access_tokens
WHERE id = pggen.arg('id')::uuid;

-- name: ListAccessTokens :many
SELECT id, name, scopes, issued_at, expires_at, revoked_at
FROM access_tokens
WHERE
  CASE WHEN pggen.arg('filter_expired')
   THEN expires_at IS NOT NULL
   ELSE TRUE
  END AND
  CASE WHEN pggen.arg('filter_revoked')
   THEN revoked_at IS NOT NULL
   ELSE TRUE
  END;
