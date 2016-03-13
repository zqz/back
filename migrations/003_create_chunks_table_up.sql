-- Initial 'chunks' table.
CREATE TABLE chunks (
  id UUID  DEFAULT uuid_generate_v4() PRIMARY KEY,
  file_id  UUID REFERENCES files (id),
  size     INTEGER,
  hash     TEXT,
  position INTEGER,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX index_chunks_on_file_id ON chunks (file_id);

-- Auto update created_at and updated_at
CREATE TRIGGER chunks_trigger_set_created_at
  BEFORE INSERT ON chunks
  FOR EACH ROW EXECUTE PROCEDURE set_created_at();

CREATE TRIGGER chunks_trigger_set_updated_at
  BEFORE UPDATE ON chunks
  FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

