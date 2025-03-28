// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: getPostsForUser.sql

package database

import (
	"context"
)

const getPostsForUser = `-- name: GetPostsForUser :many
SELECT id, created_at, updated_at, title, url, description, published_at, feed_id
FROM posts
ORDER BY published_at DESC NULLS LAST, created_at DESC, updated_at DESC
LIMIT $1
`

func (q *Queries) GetPostsForUser(ctx context.Context, limit int32) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
