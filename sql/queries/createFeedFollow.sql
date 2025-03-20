-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (INSERT INTO feed_follow (id, created_at, updated_at, user_id, feed_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
)
RETURNING *
)
SELECT inserted_feed_follow.id, inserted_feed_follow.created_at, inserted_feed_follow.updated_at, inserted_feed_follow.user_id, users.name AS user_name, inserted_feed_follow.feed_id, feeds.name AS feed_name
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;
