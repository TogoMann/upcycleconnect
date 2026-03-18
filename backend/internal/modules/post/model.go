package post

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Post struct {
	Id        int64            `db:"id" json:"id"`
	ThreadId  int64            `db:"thread_id" json:"thread_id"`
	CreatedBy int64            `db:"created_by" json:"created_by"`
	Content   string           `db:"content" json:"content"`
	Upvotes   int32            `db:"upvotes" json:"upvotes"`
	Downvotes int32            `db:"downvotes" json:"downvotes"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
	EditedAt  pgtype.Timestamp `db:"edited_at" json:"edited_at"`
}
