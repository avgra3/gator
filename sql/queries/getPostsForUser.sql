-- name: GetPostsForUser :many
SELECT *
FROM posts
ORDER BY published_at DESC NULLS LAST, created_at DESC, updated_at DESC
LIMIT $1;
