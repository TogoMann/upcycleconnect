package contract

import "github.com/jackc/pgx/v5/pgtype"

type Contract struct {
	Id        pgtype.Int8            `db:"id" json:"id"`
	Name      string           `db:"name" json:"name"`
	CreatedBy pgtype.Int8            `db:"created_by" json:"created_by"`
	Content   string           `db:"content" json:"content"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
	Until     pgtype.Timestamp `db:"until" json:"until"`
}
