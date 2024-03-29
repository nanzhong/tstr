-- name: IssueAccessToken :one
-- TODO re: ::text[] https://github.com/kyleconroy/sqlc/issues/1256
INSERT INTO access_tokens (name, token_hash, namespace_selectors, scopes, expires_at)
VALUES (
  sqlc.arg('name'), 
  sqlc.arg('token_hash'), 
  sqlc.arg('namespace_selectors'),
  sqlc.arg('scopes')::text[]::access_token_scope[], 
  sqlc.arg('expires_at')
)
RETURNING id, name, namespace_selectors, scopes::text[], issued_at, expires_at;

-- name: AuthAccessToken :one
-- TODO re: ::text[] https://github.com/kyleconroy/sqlc/issues/1256
SELECT id, name, namespace_selectors, scopes::text[], issued_at, expires_at, revoked_at
FROM access_tokens
WHERE
  token_hash = sqlc.arg('token_hash') AND
  revoked_at IS NULL AND
  (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP);

-- name: GetAccessToken :one
-- TODO re: ::text[] https://github.com/kyleconroy/sqlc/issues/1256
SELECT id, name, namespace_selectors, scopes::text[], issued_at, expires_at, revoked_at
FROM access_tokens
WHERE id = sqlc.arg('id');

-- name: ListAccessTokens :many
-- TODO re: ::text[] https://github.com/kyleconroy/sqlc/issues/1256
SELECT id, name, namespace_selectors, scopes::text[], issued_at, expires_at, revoked_at
FROM access_tokens
WHERE
  CASE WHEN sqlc.arg('include_expired')::bool
   THEN TRUE
   ELSE expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP
  END AND
  CASE WHEN sqlc.arg('include_revoked')::bool
   THEN TRUE
   ELSE revoked_at IS NULL OR revoked_at > CURRENT_TIMESTAMP
  END
ORDER BY issued_at ASC;

-- name: RevokeAccessToken :exec
UPDATE access_tokens
SET revoked_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id');
