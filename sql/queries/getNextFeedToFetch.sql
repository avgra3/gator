-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at DESC NULLS FIRST, updated_at DESC NULLS FIRST
LIMIT 1;

