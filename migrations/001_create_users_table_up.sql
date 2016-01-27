-- Required for UUID's
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Initial 'user' table.
CREATE TABLE users (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  first_name varchar(64),
  last_name varchar(64),
  username varchar(32) UNIQUE NOT NULL,
  address varchar(256),
  phone varchar(256),
  email varchar(256),
  apikey varchar(256),
  password varchar(128),
  hash varchar(256),
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  banned boolean DEFAULT false
);

CREATE INDEX index_users_on_username ON users (username);

CREATE FUNCTION set_created_at()
  RETURNS TRIGGER
  LANGUAGE plpgsql
  AS $$
    BEGIN
      NEW.created_at = now() at time zone 'utc';
      NEW.updated_at = NEW.created_at;
      RETURN NEW;
    END;
  $$;

CREATE FUNCTION set_updated_at()
  RETURNS TRIGGER
  LANGUAGE plpgsql
  AS $$
    BEGIN
      NEW.updated_at = now() at time zone 'utc';
      RETURN NEW;
    END;
  $$;

CREATE TRIGGER users_trigger_set_created_at
  BEFORE INSERT ON users
  FOR EACH ROW EXECUTE PROCEDURE set_created_at();

CREATE TRIGGER users_trigger_set_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW EXECUTE PROCEDURE set_updated_at();

