-- +goose Up
	ALTER TABLE files ADD COLUMN url text DEFAULT identifier(7);
	ALTER TABLE files ALTER COLUMN url SET STORAGE plain;

-- +goose Down
	ALTER TABLE files DROP COLUMN url;
