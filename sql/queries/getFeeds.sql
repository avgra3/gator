-- name: GetFeeds :many
SELECT feeds.name as feedName, feeds.url as feedUrl, users.name as userName
FROM feeds INNER JOIN users
ON feeds.user_id = users.id
ORDER BY feeds.updated_at DESC;
