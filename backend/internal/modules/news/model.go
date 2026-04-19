package news

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type NewsFrontend struct {
	Id        int64  `db:"id" json:"id"`
	CreatedBy *int64 `db:"created_by" json:"created_by"`
	Title     string `db:"title" json:"title"`
	Content   string `db:"content" json:"content"`
	CreatedAt string `db:"created_at" json:"created_at"`
	Upvotes   int32  `db:"upvotes" json:"upvotes"`
	Downvotes int32  `db:"downvotes" json:"downvotes"`
}

type News struct {
	Id        pgtype.Int8      `db:"id" json:"id"`
	CreatedBy pgtype.Int8      `db:"created_by" json:"created_by"`
	Title     string           `db:"title" json:"title"`
	Content   string           `db:"content" json:"content"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
	Upvotes   int32            `db:"upvotes" json:"upvotes"`
	Downvotes int32            `db:"downvotes" json:"downvotes"`
}
