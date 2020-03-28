-- this migration is noop because digests have been removed

-- -use-id-in-digest-notes-joining-table.sql replaces uuids with ids
-- as foreign keys in the digest_notes joining table.

-- +migrate Up

-- +migrate Down
