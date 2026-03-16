package contract

import "time"

type Contract struct {
	Id        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedBy int64     `db:"created_by" json:"created_by"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Until     time.Time `db:"until" json:"until"`
}
