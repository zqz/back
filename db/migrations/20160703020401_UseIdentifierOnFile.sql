-- +goose Up
	ALTER TABLE files ADD COLUMN slug text DEFAULT identifier(7);
	ALTER TABLE files ALTER COLUMN slug SET STORAGE plain;

-- +goose Down
	ALTER TABLE files DROP COLUMN slug;
