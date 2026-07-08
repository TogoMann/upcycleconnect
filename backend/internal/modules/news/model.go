package news

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type NewsType string

const (
	Actualite NewsType = "actualite"
	Conseil   NewsType = "conseil"
)

type NewsFrontend struct {
	Id        int64   `db:"id" json:"id"`
	CreatedBy *int64  `db:"created_by" json:"created_by"`
	Title     string  `db:"title" json:"title"`
	Content   string  `db:"content" json:"content"`
	Type      string  `db:"type" json:"type"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	Upvotes   int32   `db:"upvotes" json:"upvotes"`
	Downvotes int32   `db:"downvotes" json:"downvotes"`
	MyVote    *string `db:"my_vote" json:"my_vote"`
}

type NewsVote struct {
	NewsId   pgtype.Int8 `db:"news_id" json:"news_id"`
	UserId   pgtype.Int8 `db:"user_id" json:"user_id"`
	VoteType string      `db:"vote_type" json:"vote_type"`
}

type News struct {
	Id        pgtype.Int8      `db:"id" json:"id"`
	CreatedBy pgtype.Int8      `db:"created_by" json:"created_by"`
	Title     string           `db:"title" json:"title"`
	Content   string           `db:"content" json:"content"`
	Type      NewsType         `db:"type" json:"type"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
	Upvotes   int32            `db:"upvotes" json:"upvotes"`
	Downvotes int32            `db:"downvotes" json:"downvotes"`
}
