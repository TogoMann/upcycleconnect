package thread

import "time"

type ThreadDto struct {
	CreatedBy int64     `db:"created_by" json:"created_by"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
