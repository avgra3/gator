// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: deleteAllFeeds.sql

package database

import (
	"context"
)

const deleteAllFeeds = `-- name: DeleteAllFeeds :exec
DELETE FROM feeds
`

func (q *Queries) DeleteAllFeeds(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllFeeds)
	return err
}
