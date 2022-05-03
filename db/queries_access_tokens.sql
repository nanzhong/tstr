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
  CASE WHEN pggen.arg('include_expired')
   THEN TRUE
   ELSE expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP
  END AND
  CASE WHEN pggen.arg('include_revoked')
   THEN TRUE
   ELSE revoked_at IS NULL OR revoked_at > CURRENT_TIMESTAMP
  END;
