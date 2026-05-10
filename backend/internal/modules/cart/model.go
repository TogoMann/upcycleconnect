package cart

import (
	"github.com/jackc/pgx/v5/pgtype"
	"backend/internal/modules/listing"
)

type CartItem struct {
	Id        pgtype.Int8      `db:"id" json:"id"`
	UserId    pgtype.Int8      `db:"user_id" json:"user_id"`
	ListingId pgtype.Int8      `db:"listing_id" json:"listing_id"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
}

type CartItemWithListing struct {
	CartItem
	Listing listing.Listing `json:"listing"`
}
