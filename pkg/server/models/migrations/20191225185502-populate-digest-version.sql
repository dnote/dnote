-- this migration is noop because digests have been removed

-- populate-digest-version.sql populates the `version` column for the digests
-- by assigining an incremental integer scoped to a repetition rule that each
-- digest belongs, ordered by created_at timestamp of the digests.

-- +migrate Up

-- +migrate Down
