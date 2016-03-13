-- Initial 'files' table.
CREATE TABLE files (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  size   INTEGER,
  chunks INTEGER,
  state  INTEGER,
  name   TEXT,
  hash   TEXT,
  type   TEXT,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX index_files_on_hash ON files (hash);

-- Auto update created_at and updated_at
CREATE TRIGGER files_trigger_set_created_at
  BEFORE INSERT ON files
  FOR EACH ROW EXECUTE PROCEDURE set_created_at();

CREATE TRIGGER files_trigger_set_updated_at
  BEFORE UPDATE ON files
  FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

