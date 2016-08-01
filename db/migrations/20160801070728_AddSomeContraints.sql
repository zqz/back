-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE files ALTER COLUMN state SET NOT NULL;
ALTER TABLE files ALTER COLUMN hash SET NOT NULL;
ALTER TABLE files ALTER COLUMN "type" SET NOT NULL;
ALTER TABLE files ALTER COLUMN name SET NOT NULL;
ALTER TABLE files ALTER COLUMN num_chunks SET NOT NULL;
ALTER TABLE files ALTER COLUMN slug SET NOT NULL;
ALTER TABLE files ALTER COLUMN "size" SET NOT NULL;

ALTER TABLE users ALTER COLUMN first_name SET NOT NULL;
ALTER TABLE users ALTER COLUMN last_name SET NOT NULL;
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ALTER COLUMN phone SET NOT NULL;
ALTER TABLE users ALTER COLUMN banned SET NOT NULL;
ALTER TABLE users ALTER COLUMN hash SET NOT NULL;

ALTER TABLE chunks ALTER COLUMN file_id SET NOT NULL;
ALTER TABLE chunks ALTER COLUMN size SET NOT NULL;
ALTER TABLE chunks ALTER COLUMN hash SET NOT NULL;
ALTER TABLE chunks ALTER COLUMN position SET NOT NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
-- Nah bro.

