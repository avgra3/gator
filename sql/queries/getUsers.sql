-- name: GetUsers :many
SELECT name
FROM users
ORDER BY updated_at DESC;
