
-- +migrate Up

-- Configure full text search
CREATE TEXT SEARCH DICTIONARY english_nostop (
  Template = snowball,
  Language = english
);

CREATE TEXT SEARCH CONFIGURATION public.english_nostop ( COPY = pg_catalog.english );

ALTER TEXT SEARCH CONFIGURATION public.english_nostop
ALTER MAPPING FOR asciiword, asciihword, hword_asciipart, hword, hword_part, word WITH english_nostop;


-- Create a trigger
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION note_tsv_trigger() RETURNS trigger AS $$
begin
  new.tsv := setweight(to_tsvector('english_nostop', new.body), 'A');
  return new;
end
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tsvectorupdate ON notes;
CREATE TRIGGER tsvectorupdate
BEFORE INSERT OR UPDATE ON notes
FOR EACH ROW EXECUTE PROCEDURE note_tsv_trigger();
-- +migrate StatementEnd

-- index tsv
CREATE INDEX IF NOT EXISTS idx_notes_tsv
ON notes
USING gin(tsv);

-- initialize tsv
UPDATE notes
SET tsv = setweight(to_tsvector('english_nostop', notes.body), 'A')
WHERE notes.encrypted = false;

-- +migrate Down
