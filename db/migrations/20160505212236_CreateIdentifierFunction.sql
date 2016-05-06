-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION identifier(len int DEFAULT 10) RETURNS text
  LANGUAGE plpgsql
  AS $$
    DECLARE
      alphabet text;
      alphabet_size int;
    BEGIN
      alphabet := 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
      alphabet_size := character_length(alphabet);

      RETURN array_to_string(
        array(
          SELECT substr(alphabet, trunc(random() * alphabet_size)::integer + 1, 1)
          FROM generate_series(1, len)
        ),
        ''
      );
    END;
  $$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION IF EXISTS identifier(int);
