package listingorder

import (
	db "backend/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]ListingOrder, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, listing_id, user_id, price, created_at, status FROM listing_order")
	if err != nil {
		return nil, fmt.Errorf("package listing_order/repo GetAll query: %w", err)
	}

	listingOrders, err := pgx.CollectRows(rows, pgx.RowToStructByName[ListingOrder])

	if err != nil {
		return nil, fmt.Errorf("package news/repo GetAll: %v", err.Error())
	}

	return listingOrders, nil
}

func (r *Repository) GetById(id int64) (*ListingOrder, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, listing_id, user_id, price, created_at, status FROM listing_order WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package listing_order/repo GetById query: %w", err)
	}

	listingOrder, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[ListingOrder])

	if err != nil {
		return nil, fmt.Errorf("package listing_order/repo GetById: %v", err.Error())
	}
	return &listingOrder, nil
}

func (r *Repository) Create(listingOrderDto ListingOrder) (int64, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO listing_order (listing_id, user_id, price, status) VALUES ($1, $2, $3, $4)",
		listingOrderDto.ListingId, listingOrderDto.UserId, listingOrderDto.Price, listingOrderDto.Status)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), err
}

func (r *Repository) Delete(id int64) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE listing_order WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package listing_order/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id int64) (bool, error) {
	var idFound int64

	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM listing_order WHERE id = $1", id).Scan(&idFound)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package listing_order/repo ExistsById query: %w", err)
	}

	return true, nil
}
