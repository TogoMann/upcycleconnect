package courseorder

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CourseOrder struct {
	Id       int64            `db:"id" json:"id"`
	CourseId int64            `db:"course_id" json:"course_id"`
	BuyerId  int64            `db:"buyer_id" json:"buyer_id"`
	Price    pgtype.Numeric   `db:"price" json:"price"`
	BookedAt pgtype.Timestamp `db:"booked_at" json:"booked_at"`
}
