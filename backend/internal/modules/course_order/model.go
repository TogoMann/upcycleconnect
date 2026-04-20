package courseorder

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CourseOrder struct {
	Id                    pgtype.Int8      `db:"id" json:"id"`
	CourseId              pgtype.Int8      `db:"course_id" json:"course_id"`
	BuyerId               pgtype.Int8      `db:"buyer_id" json:"buyer_id"`
	StripePaymentIntentId string           `db:"stripe_payment_intent_id" json:"stripe_payment_intent_id"`
	Price                 pgtype.Numeric   `db:"price" json:"price"`
	BookedAt              pgtype.Timestamp `db:"booked_at" json:"booked_at"`
}

type CourseOrderWithCourse struct {
	CourseOrder
	CourseName string `db:"course_name" json:"course_name"`
}
