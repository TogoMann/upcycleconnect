package comments

import "github.com/jackc/pgx/v5/pgtype"

type Comments struct {
	Id        pgtype.Int8      `db:"id" json:"id"`
	NewsId    pgtype.Int8      `db:"news_id" json:"news_id"`
	CreatedBy pgtype.Int8      `db:"created_by" json:"created_by"`
	ParentId  pgtype.Int8      `db:"parent_id" json:"parent_id"`
	Content   string           `db:"content" json:"content"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
	Upvotes   int32            `db:"upvotes" json:"upvotes"`
	Downvotes int32            `db:"downvotes" json:"downvotes"`
}
