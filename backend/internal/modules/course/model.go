package course

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Course struct {
	Id         int64          `db:"id" json:"id"`
	CreatedBy  int64          `db:"created_by" json:"created_by"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
	Approved   bool           `db:"approved" json:"approved"`
	ApprovedBy int64          `db:"approved_by" json:"approved_by"`
	ApprovedAt time.Time      `db:"approved_at" json:"approved_at"`
	Price      pgtype.Numeric `db:"price" json:"price"`
}
