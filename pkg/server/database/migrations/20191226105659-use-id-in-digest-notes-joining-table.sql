-- -use-id-in-digest-notes-joining-table.sql replaces uuids with ids
-- as foreign keys in the digest_notes joining table.

-- +migrate Up

DO $$

-- +migrate StatementBegin
BEGIN
  PERFORM column_name FROM information_schema.columns WHERE table_name= 'digest_notes' and column_name = 'digest_id';
  IF NOT found THEN
    ALTER TABLE digest_notes ADD COLUMN digest_id int;
  END IF;
  PERFORM column_name FROM information_schema.columns WHERE table_name= 'digest_notes' and column_name = 'note_id';
  IF NOT found THEN
    ALTER TABLE digest_notes ADD COLUMN note_id int;
  END IF;

  -- migrate if digest_notes.digest_uuid exists
  PERFORM column_name FROM information_schema.columns WHERE table_name= 'digest_notes' and column_name = 'digest_uuid';
  IF found THEN
    -- update note_id
    UPDATE digest_notes
    SET note_id=t1.note_id
    FROM (
      SELECT notes.id AS note_id, notes.uuid AS note_uuid
      FROM digest_notes
      INNER JOIN notes ON notes.uuid = digest_notes.note_uuid
    ) AS t1
    WHERE digest_notes.note_uuid = t1.note_uuid;

    -- update digest_id
    UPDATE digest_notes
    SET digest_id=t1.digest_id
    FROM (
      SELECT digests.id AS digest_id, digests.uuid AS digest_uuid
      FROM digest_notes
      INNER JOIN digests ON digests.uuid = digest_notes.digest_uuid
    ) AS t1
    WHERE digest_notes.digest_uuid = t1.digest_uuid;

    ALTER TABLE digest_notes DROP COLUMN digest_uuid;
    ALTER TABLE digest_notes DROP COLUMN note_uuid;
  END IF;
END; $$
-- +migrate StatementEnd


-- +migrate Down
