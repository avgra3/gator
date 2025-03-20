-- name: Unfollow :exec
DELETE FROM feed_follow
WHERE feed_follow.user_id = $1
AND feed_id = (
	SELECT id
	FROM feeds
	WHERE url = $2
);

