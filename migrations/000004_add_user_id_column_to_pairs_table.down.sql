DROP INDEX IF EXISTS idx_url_pairs_user_id;

ALTER TABLE url_pairs
DROP COLUMN IF EXISTS user_id;