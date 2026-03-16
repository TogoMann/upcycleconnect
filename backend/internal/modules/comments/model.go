package comments

import "time"

type Comments struct {
	Id        int64     `db:"id" json:"id"`
	NewsId    int64     `db:"news_id" json:"news_id"`
	CreatedBy int64     `db:"created_by" json:"created_by"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Upvotes   int32     `db:"upvotes" json:"upvotes"`
	Downvotes int32     `db:"downvotes" json:"downvotes"`
}
