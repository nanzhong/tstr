-- migrate:up
ALTER TABLE runs ADD result_data jsonb DEFAULT '{}'::jsonb;



-- migrate:down

ALTER TABLE runs DROP COLUMN result_data;