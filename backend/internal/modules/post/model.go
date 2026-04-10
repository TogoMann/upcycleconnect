package post

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Post struct {
	Id        pgtype.Int8      `db:"id" json:"id"`
	ThreadId  pgtype.Int8      `db:"thread_id" json:"thread_id"`
	CreatedBy pgtype.Int8      `db:"created_by" json:"created_by"`
	ParentId  pgtype.Int8      `db:"parent_id" json:"parent_id"`
	Content   string           `db:"content" json:"content"`
	Upvotes   int32            `db:"upvotes" json:"upvotes"`
	Downvotes int32            `db:"downvotes" json:"downvotes"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
	EditedAt  pgtype.Timestamp `db:"edited_at" json:"edited_at"`
}

type ThreadPosts struct {
	Post
	Title   string `db:"title" json:"title"`
	Content string `db:"thread_content" json:"thread_content"`
}
