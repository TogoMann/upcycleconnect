package course

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Course struct {
	Id         pgtype.Int8            `db:"id" json:"id"`
	CreatedBy  pgtype.Int8            `db:"created_by" json:"created_by"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
	Approved   bool             `db:"approved" json:"approved"`
	ApprovedBy pgtype.Int8            `db:"approved_by" json:"approved_by"`
	ApprovedAt pgtype.Timestamp `db:"approved_at" json:"approved_at"`
	Price      pgtype.Numeric   `db:"price" json:"price"`
}
