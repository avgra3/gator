// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: unfollow.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const unfollow = `-- name: Unfollow :exec
DELETE FROM feed_follow
WHERE feed_follow.user_id = $1
AND feed_id = (
	SELECT id
	FROM feeds
	WHERE url = $2
)
`

type UnfollowParams struct {
	UserID uuid.NullUUID
	Url    string
}

func (q *Queries) Unfollow(ctx context.Context, arg UnfollowParams) error {
	_, err := q.db.ExecContext(ctx, unfollow, arg.UserID, arg.Url)
	return err
}
