-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP DEFAULT NULL;

-- +goose Down
ALTER TABLE IF EXISTS feeds
DROP COLUMN IF EXISTS last_fetched_at;
