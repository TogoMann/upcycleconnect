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
	Date       pgtype.Timestamp `db:"date" json:"date"`
	StartDate  pgtype.Timestamp `db:"start_date" json:"start_date"`
	EndDate    pgtype.Timestamp `db:"end_date" json:"end_date"`
	Location   string           `db:"location" json:"location"`
	CreatedBy  pgtype.Int8      `db:"created_by" json:"created_by"`
	CreatedAt  pgtype.Timestamp `db:"created_at" json:"created_at"`
}
