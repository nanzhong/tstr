-- migrate:up

CREATE TABLE tests (
       id uuid PRIMARY KEY,
       name varchar NOT NULL,
       labels jsonb,
       cron_schedule varchar,
       registered_at timestamptz NOT NULL,
       updated_at timestamptz NOT NULL,
       archived_at timestamptz
);
CREATE INDEX ON tests USING GIN (labels);

CREATE TABLE test_run_configs (
       id serial PRIMARY KEY,
       test_id uuid references tests(id),
       container_image varchar NOT NULL,
       command varchar,
       env jsonb,
       created_at timestamptz NOT NULL
);

CREATE TABLE test_suites (
       id uuid PRIMARY KEY,
       name varchar NOT NULL,
       labels jsonb,
       created_at timestamptz NOT NULL,
       updated_at timestamptz NOT NULL,
       archived_at timestamptz
);
CREATE INDEX ON test_suites(archived_at);

CREATE TABLE runners (
       id uuid PRIMARY KEY,
       name varchar NOT NULL,
       accept_test_labels jsonb,
       reject_test_labels jsonb,
       registered_at timestamptz NOT NULL,
       approved_at timestamptz NOT NULL,
       last_heartbeat_at timestamptz NOT NULL
);
CREATE INDEX ON runners(last_heartbeat_at);

CREATE TYPE run_result AS ENUM ('pass', 'fail', 'error');
CREATE TABLE runs (
       id uuid PRIMARY KEY,
       test_id uuid references tests(id),
       test_run_config_id integer references test_run_configs(id),
       runner_id uuid references runners(id),
       result run_result,
       logs jsonb,
       scheduled_at timestamptz NOT NULL,
       started_at timestamptz,
       finished_at timestamptz
);
CREATE INDEX ON runs(test_id, test_run_config_id);
CREATE INDEX ON runs(runner_id);
CREATE INDEX ON runs(result);
CREATE INDEX ON runs(scheduled_at, started_at, finished_at);

CREATE TYPE access_token_scope AS ENUM ('admin', 'control_r', 'control_rw');
CREATE TABLE access_tokens (
       id uuid PRIMARY KEY,
       name varchar NOT NULL,
       token_hash varchar NOT NULL,
       scopes access_token_scope[],
       issued_at timestamptz NOT NULL,
       expires_at timestamptz NOT NULL,
       revoked_at timestamptz
);
CREATE INDEX ON access_tokens (revoked_at, expires_at, token_hash);

-- migrate:down

DROP TABLE access_tokens;
DROP TYPE access_token_scope;

DROP TABLE runs;
DROP TYPE run_result;

DROP TABLE runners;

DROP TABLE test_suites;

DROP TABLE test_run_configs;

DROP TABLE tests;
