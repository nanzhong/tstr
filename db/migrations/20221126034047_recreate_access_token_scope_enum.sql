-- migrate:up
ALTER TYPE access_token_scope RENAME VALUE 'control_rw' TO 'control';
UPDATE access_tokens SET scopes = array_remove(scopes, 'control_r');
CREATE TYPE access_token_scope_new AS ENUM ('admin', 'control', 'data', 'runner');
ALTER TABLE access_tokens ALTER COLUMN scopes TYPE access_token_scope_new[] USING (scopes::text[]::access_token_scope_new[]);
DROP TYPE access_token_scope;
ALTER TYPE access_token_scope_new RENAME TO access_token_scope;

-- migrate:down
ALTER TYPE access_token_scope RENAME VALUE 'control' TO 'control_rw';
ALTER TYPE access_token_scope ADD VALUE 'control_r' AFTER 'control_rw';
