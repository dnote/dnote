-- this migration is noop because digests have been removed

-- delete-outdated-digests.sql deletes digests that do not belong to any repetition rules,
-- along with digest_notes associations.

-- +migrate Up

-- +migrate Down
