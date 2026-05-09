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

type EventFull struct {
	Id               pgtype.Int8      `db:"id" json:"id"`
	Approved         bool             `db:"approved" json:"approved"`
	ApprovedBy       pgtype.Int8      `db:"approved_by" json:"approved_by"`
	ApprovedAt       pgtype.Timestamp `db:"approved_at" json:"approved_at"`
	Price            pgtype.Numeric   `db:"price" json:"price"`
	Date             pgtype.Date      `db:"date" json:"date"`
	StartTime        pgtype.Time      `db:"start_time" json:"start_time"`
	EndTime          pgtype.Time      `db:"end_time" json:"end_time"`
	Location         string           `db:"location" json:"location"`
	CreatedBy        pgtype.Int8      `db:"created_by" json:"created_by"`
	CreatedAt        pgtype.Timestamp `db:"created_at" json:"created_at"`
	CreatorUsername  pgtype.Text      `db:"creator_username" json:"creator_username"`
	CreatorEmail     pgtype.Text      `db:"creator_email" json:"creator_email"`
	ApproverUsername pgtype.Text      `db:"approver_username" json:"approver_username"`
	ApproverEmail    pgtype.Text      `db:"approver_email" json:"approver_email"`
}
