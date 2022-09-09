-- migrate:up

DROP TABLE test_suites;

-- migrate:down

CREATE TABLE test_suites (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  name varchar NOT NULL UNIQUE,
  labels jsonb,
  created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamptz DEFAULT CURRENT_TIMESTAMP
);
