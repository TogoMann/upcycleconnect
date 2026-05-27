package cart

import (
	"backend/internal/modules/course"
	"backend/internal/modules/event"
	"backend/internal/modules/listing"
	"github.com/jackc/pgx/v5/pgtype"
)

type CartItem struct {
	Id        pgtype.Int8      `db:"id" json:"id"`
	UserId    pgtype.Int8      `db:"user_id" json:"user_id"`
	ListingId pgtype.Int8      `db:"listing_id" json:"listing_id"`
	EventId   pgtype.Int8      `db:"event_id" json:"event_id"`
	CourseId  pgtype.Int8      `db:"course_id" json:"course_id"`
	CreatedAt pgtype.Timestamp `db:"created_at" json:"created_at"`
}

type CartItemDetailed struct {
	CartItem
	Listing *listing.Listing `json:"listing,omitempty"`
	Event   *event.Event     `json:"event,omitempty"`
	Course  *course.Course   `json:"course,omitempty"`
}
