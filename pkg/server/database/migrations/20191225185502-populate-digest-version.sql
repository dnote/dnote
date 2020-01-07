-- populate-digest-version.sql populates the `version` column for the digests
-- by assigining an incremental integer scoped to a repetition rule that each
-- digest belongs, ordered by created_at timestamp of the digests.

-- +migrate Up
UPDATE digests
SET version=t1.version
FROM (
  SELECT
    digests.uuid,
    ROW_NUMBER() OVER (PARTITION BY digests.rule_id ORDER BY digests.created_at) AS version
  FROM digests
  WHERE digests.rule_id IS NOT NULL
) AS t1
WHERE digests.uuid = t1.uuid;

-- +migrate Down
UPDATE digests
SET version=0;
