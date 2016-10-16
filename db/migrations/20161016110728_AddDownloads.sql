
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE downloads (
  id SERIAL PRIMARY KEY,
  ip INET,
  cache_hit BOOLEAN NOT NULL,
  file_id UUID references files(id),
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX index_downloads_on_file_id ON downloads (file_id);

-- Auto update created_at
CREATE TRIGGER downloads_trigger_set_created_at
  BEFORE INSERT ON downloads
  FOR EACH ROW EXECUTE PROCEDURE set_created_at();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TRIGGER downloads_trigger_set_created_at ON downloads;
DROP TABLE downloads;

