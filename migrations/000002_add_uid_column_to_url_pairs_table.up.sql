ALTER TABLE url_pairs
    ADD COLUMN uid UUID;

ALTER TABLE url_pairs
    ADD CONSTRAINT url_pairs_uid_key UNIQUE (uid);