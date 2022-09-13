-- migrate:up

ALTER TABLE tests ADD COLUMN namespace varchar NOT NULL DEFAULT 'default';
ALTER TABLE tests ALTER COLUMN namespace DROP DEFAULT;
CREATE INDEX tests_namespace_idx ON tests(namespace);

ALTER TABLE runners ADD COLUMN namespace_selectors varchar[] NOT NULL DEFAULT '{".*"}';
ALTER TABLE runners ALTER COLUMN namespace_selectors DROP DEFAULT;

ALTER TABLE access_tokens ADD COLUMN namespace_selectors varchar[] NOT NULL DEFAULT '{".*"}';
ALTER TABLE access_tokens ALTER COLUMN namespace_selectors DROP DEFAULT;
-- migrate:down

ALTER TABLE tests DROP COLUMN namespace;
ALTER TABLE runners DROP COLUMN namespace_selectors;
ALTER TABLE access_tokens DROP COLUMN namespace_selectors;
