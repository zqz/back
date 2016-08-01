-- +goose Up
ALTER TABLE files RENAME COLUMN chunks TO num_chunks;

-- +goose Down
ALTER TABLE files RENAME COLUMN num_chunks TO chunks;
