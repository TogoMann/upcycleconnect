package cart

import (
	db "backend/internal/database"
	"backend/internal/modules/course"
	"backend/internal/modules/event"
	"backend/internal/modules/listing"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByUserId(userId pgtype.Int8) ([]CartItemDetailed, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			ci.id, ci.user_id, ci.listing_id, ci.event_id, ci.course_id, ci.created_at,
			l.id, COALESCE(l.name, ''), COALESCE(l.description, ''), COALESCE(l.category::text, ''), l.item_id, l.city_id, l.created_by, l.created_at, COALESCE(l.approved, false), l.approved_by, l.approved_at, COALESCE(l.status::text, ''), l.price, l.image_url,
			e.id, COALESCE(e.approved, false), e.approved_by, e.approved_at, e.price, e.date, e.start_time, e.end_time, COALESCE(e.location, ''), e.created_by, e.created_at,
			c.id, COALESCE(c.name, ''), COALESCE(c.description, ''), COALESCE(c.max_capacity, 0), c.created_by, c.created_at, COALESCE(c.approved, false), c.approved_by, c.approved_at, c.price, c.date, c.start_time, c.end_time
		FROM cart_item ci
		LEFT JOIN listing l ON ci.listing_id = l.id
		LEFT JOIN event e ON ci.event_id = e.id
		LEFT JOIN course c ON ci.course_id = c.id
		WHERE ci.user_id = $1
	`, userId)
	if err != nil {
		fmt.Printf("Cart Query Error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	items := []CartItemDetailed{}
	for rows.Next() {
		var item CartItemDetailed
		var l listing.Listing
		var e event.Event
		var c course.Course

		err := rows.Scan(
			&item.Id, &item.UserId, &item.ListingId, &item.EventId, &item.CourseId, &item.CreatedAt,
			&l.Id, &l.Name, &l.Description, &l.Category, &l.ItemId, &l.CityId, &l.CreatedBy, &l.CreatedAt, &l.Approved, &l.ApprovedBy, &l.ApprovedAt, &l.Status, &l.Price, &l.ImageUrl,
			&e.Id, &e.Approved, &e.ApprovedBy, &e.ApprovedAt, &e.Price, &e.Date, &e.StartTime, &e.EndTime, &e.Location, &e.CreatedBy, &e.CreatedAt,
			&c.Id, &c.Name, &c.Description, &c.MaxCapacity, &c.CreatedBy, &c.CreatedAt, &c.Approved, &c.ApprovedBy, &c.ApprovedAt, &c.Price, &c.Date, &c.StartTime, &c.EndTime,
		)
		if err != nil {
			fmt.Printf("Cart Scan Error: %v\n", err)
			return nil, err
		}

		if l.Id.Valid {
			item.Listing = &l
		}
		if e.Id.Valid {
			item.Event = &e
		}
		if c.Id.Valid {
			item.Course = &c
		}

		items = append(items, item)
	}
	fmt.Printf("Cart fetched for user %d: %d items\n", userId.Int64, len(items))
	return items, nil
}

func (r *Repository) Add(userId, listingId, eventId, courseId pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, `
		INSERT INTO cart_item (user_id, listing_id, event_id, course_id) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT DO NOTHING
	`, userId, listingId, eventId, courseId)
	return err
}

func (r *Repository) Remove(userId, listingId, eventId, courseId pgtype.Int8) error {
	query := "DELETE FROM cart_item WHERE user_id = $1"
	args := []interface{}{userId}
	if listingId.Valid {
		query += " AND listing_id = $2"
		args = append(args, listingId)
	} else if eventId.Valid {
		query += " AND event_id = $2"
		args = append(args, eventId)
	} else if courseId.Valid {
		query += " AND course_id = $2"
		args = append(args, courseId)
	}

	_, err := r.db.Exec(db.Ctx, query, args...)
	return err
}

func (r *Repository) Clear(userId pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM cart_item WHERE user_id = $1", userId)
	return err
}

func (r *Repository) ClearDirectPay(userId pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM cart_item WHERE user_id = $1 AND (event_id IS NOT NULL OR course_id IS NOT NULL)", userId)
	return err
}
