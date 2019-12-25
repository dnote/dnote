
-- +migrate Up

ALTER TABLE digests DROP CONSTRAINT digests_pkey;
ALTER TABLE digests ADD PRIMARY KEY (id);

-- +migrate Down

ALTER TABLE digests DROP CONSTRAINT digests_pkey;
ALTER TABLE digests ADD PRIMARY KEY (uuid);
