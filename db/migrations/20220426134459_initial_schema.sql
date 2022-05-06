-- migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tests (
       id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
       name varchar NOT NULL UNIQUE,
       labels jsonb,
       cron_schedule varchar,
       registered_at timestamptz DEFAULT CURRENT_TIMESTAMP,
       updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
       archived_at timestamptz
);
CREATE INDEX ON tests USING GIN (labels);
CREATE UNIQUE INDEX tests_unique_name ON tests(name) WHERE archived_at IS NULL;

CREATE TABLE test_run_configs (
       id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
       test_id uuid REFERENCES tests(id),
       container_image varchar NOT NULL,
       command varchar,
       args varchar[],
       env jsonb,
       created_at timestamptz DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON test_run_configs(created_at);

CREATE TABLE test_suites (
       id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
       name varchar NOT NULL UNIQUE,
       labels jsonb,
       created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
       updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
       archived_at timestamptz
);
CREATE INDEX ON test_suites(archived_at);
CREATE UNIQUE INDEX test_suites_unique_name ON test_suites(name) WHERE archived_at IS NULL;

CREATE TABLE runners (
       id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
       name varchar NOT NULL,
       accept_test_label_selectors jsonb,
       reject_test_label_selectors jsonb,
       registered_at timestamptz DEFAULT CURRENT_TIMESTAMP,
       last_heartbeat_at timestamptz DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON runners(last_heartbeat_at);

CREATE TYPE run_result AS ENUM ('unknown', 'pass', 'fail', 'error');
CREATE TABLE runs (
       id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
       test_id uuid REFERENCES tests(id),
       test_run_config_id uuid REFERENCES test_run_configs(id),
       runner_id uuid REFERENCES runners(id),
       result run_result DEFAULT 'unknown'::run_result,
       logs jsonb,
       scheduled_at timestamptz DEFAULT CURRENT_TIMESTAMP,
       started_at timestamptz,
       finished_at timestamptz
);
CREATE INDEX ON runs(test_id, test_run_config_id);
CREATE INDEX ON runs(runner_id);
CREATE INDEX ON runs(result);
CREATE INDEX ON runs(scheduled_at, started_at, finished_at);

CREATE TYPE access_token_scope AS ENUM ('admin', 'control_r', 'control_rw', 'runner');
CREATE TABLE access_tokens (
       id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
       name varchar NOT NULL,
       token_hash varchar NOT NULL,
       scopes access_token_scope[],
       issued_at timestamptz DEFAULT CURRENT_TIMESTAMP,
       expires_at timestamptz,
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
