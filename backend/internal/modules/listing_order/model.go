package listingorder

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ListingOrderStatus string

const (
	Pending   = "pending"
	Paid      = "paid"
	Shipped   = "shipped"
	Completed = "completed"
)

type ListingOrder struct {
	Id                    pgtype.Int8        `db:"id" json:"id"`
	ListingId             pgtype.Int8        `db:"listing_id" json:"listing_id"`
	UserId                pgtype.Int8        `db:"user_id" json:"user_id"`
	StripePaymentIntentId string             `db:"stripe_payment_intent_id" json:"stripe_payment_intent_id"`
	Price                 pgtype.Numeric     `db:"price" json:"price"`
	CreatedAt             pgtype.Timestamp   `db:"created_at" json:"created_at"`
	Status                ListingOrderStatus `db:"status" json:"status"`
}
