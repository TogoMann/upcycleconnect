package cart

import (
	db "backend/internal/database"
	"backend/internal/modules/listing"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByUserId(userId pgtype.Int8) ([]CartItemWithListing, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			ci.id, ci.user_id, ci.listing_id, ci.created_at,
			l.id, l.name, l.description, l.category, l.item_id, l.city_id, l.created_by, l.created_at, l.approved, l.approved_by, l.approved_at, l.status, l.price, l.image_url
		FROM cart_item ci
		JOIN listing l ON ci.listing_id = l.id
		WHERE ci.user_id = $1
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []CartItemWithListing{}
	for rows.Next() {
		var item CartItemWithListing
		var l listing.Listing
		err := rows.Scan(
			&item.Id, &item.UserId, &item.ListingId, &item.CreatedAt,
			&l.Id, &l.Name, &l.Description, &l.Category, &l.ItemId, &l.CityId, &l.CreatedBy, &l.CreatedAt, &l.Approved, &l.ApprovedBy, &l.ApprovedAt, &l.Status, &l.Price, &l.ImageUrl,
		)
		if err != nil {
			return nil, err
		}
		item.Listing = l
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) Add(userId, listingId pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "INSERT INTO cart_item (user_id, listing_id) VALUES ($1, $2) ON CONFLICT (user_id, listing_id) DO NOTHING", userId, listingId)
	return err
}

func (r *Repository) Remove(userId, listingId pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM cart_item WHERE user_id = $1 AND listing_id = $2", userId, listingId)
	return err
}

func (r *Repository) Clear(userId pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM cart_item WHERE user_id = $1", userId)
	return err
}
