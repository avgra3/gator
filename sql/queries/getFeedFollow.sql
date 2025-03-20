-- name: GetFeedFollow :many
SELECT feed_follow.id, feed_follow.created_at, feed_follow.updated_at, users.name AS user_name, feeds.name AS feed_name
FROM feed_follow
INNER JOIN users ON feed_follow.user_id = users.id
INNER JOIN feeds ON feed_follow.feed_id = feeds.id
WHERE feed_follow.user_id = $1;
