package event

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Event struct {
	Id         pgtype.Int8      `db:"id" json:"id"`
	Approved   bool             `db:"approved" json:"approved"`
	ApprovedBy pgtype.Int8      `db:"approved_by" json:"approved_by"`
	ApprovedAt pgtype.Timestamp `db:"approved_at" json:"approved_at"`
	Price      pgtype.Numeric   `db:"price" json:"price"`
	Date       pgtype.Date      `db:"date" json:"date"`
	StartTime  pgtype.Time      `db:"start_time" json:"start_time"`
	EndTime    pgtype.Time      `db:"end_time" json:"end_time"`
	Location   string           `db:"location" json:"location"`
	CreatedBy  pgtype.Int8      `db:"created_by" json:"created_by"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
}
