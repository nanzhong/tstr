-- migrate:up
ALTER TABLE runs ADD computed_data jsonb DEFAULT '{}'::jsonb;



-- migrate:down

ALTER TABLE runs DROP COLUMN computed_data;