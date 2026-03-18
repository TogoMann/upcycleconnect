package entry

import "github.com/jackc/pgx/v5/pgtype"

type Entry struct {
	Id        int64            `db:"id" json:"id"`
	CreatedBy int64            `db:"created_by" json:"created_by"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
	Schedule  pgtype.Timestamp `db:"schedule" json:"schedule"`
	Start     pgtype.Timestamp `db:"start" json:"start"`
	Ending    pgtype.Timestamp `db:"ending" json:"ending"`
}
