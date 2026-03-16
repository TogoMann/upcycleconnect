package thread

import (
	"time"
)

type Thread struct {
	Id         int64     `db:"id" json:"id"`
	CreatedBy  int64     `db:"created_by" json:"created_by"`
	Title      string    `db:"title" json:"title"`
	Content    string    `db:"content" json:"content"`
	Upvotes    int32     `db:"upvotes" json:"upvotes"`
	Downvotes  int32     `db:"downvotes" json:"downvotes"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	LastPostAt time.Time `db:"last_post_at" json:"last_post_at"`
}
