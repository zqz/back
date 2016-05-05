-- +goose Up
-- Initial 'thumbnails' table.
CREATE TABLE thumbnails (
  id UUID  DEFAULT uuid_generate_v4() PRIMARY KEY,
  file_id  UUID REFERENCES files (id),
  size     INTEGER,
  hash     TEXT,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX index_thumbnails_on_file_id ON thumbnails (file_id);

-- Auto update created_at and updated_at
CREATE TRIGGER thumbnails_trigger_set_created_at
  BEFORE INSERT ON thumbnails
  FOR EACH ROW EXECUTE PROCEDURE set_created_at();

CREATE TRIGGER thumbnails_trigger_set_updated_at
  BEFORE UPDATE ON thumbnails
  FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

-- +goose Down
DROP TRIGGER thumbnails_trigger_set_created_at ON thumbnails;
DROP TRIGGER thumbnails_trigger_set_updated_at ON thumbnails;
DROP TABLE thumbnails;

