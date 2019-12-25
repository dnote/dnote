-- delete-outdated-digests.sql deletes digests that do not belong to any repetition rules,
-- along with digest_notes associations.

-- +migrate Up
DELETE
FROM digest_notes
USING digests
WHERE
  digests.rule_id IS NULL AND
  digests.id = digest_notes.digest_id;

DELETE FROM digests
WHERE digests.rule_id IS NULL;

-- +migrate Down
