ALTER TABLE url_pairs
    ADD COLUMN user_id UUID;

CREATE INDEX idx_url_pairs_user_id ON url_pairs (user_id);