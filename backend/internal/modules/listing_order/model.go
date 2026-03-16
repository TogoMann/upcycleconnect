package listingorder

import (
	"time"

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
	Id        int64              `db:"id" json:"id"`
	ListingId int64              `db:"listing_id" json:"listing_id"`
	UserId    int64              `db:"user_id" json:"user_id"`
	Price     pgtype.Numeric     `db:"price" json:"price"`
	CreatedAt time.Time          `db:"created_at" json:"created_at"`
	status    ListingOrderStatus `db:"status" json:"status"`
}
