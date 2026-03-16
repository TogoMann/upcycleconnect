package news

import "time"

type News struct {
	Id        int64     `db:"id" json:"id"`
	CreatedBy int64     `db:"created_by" json:"created_by"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Upvotes   int32     `db:"upvotes" json:"upvotes"`
	Downvotes int32     `db:"downvotes" json:"downvotes"`
}
