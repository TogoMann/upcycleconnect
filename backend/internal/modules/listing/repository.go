package listing

import (
	db "backend/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]Listing, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			l.id, l.name, l.description, l.category, l.item_id, l.city_id, l.created_by, l.created_at, 
			l.approved, l.approved_by, l.approved_at, l.status, l.price, COALESCE(l.image_url, '') as image_url,
			COALESCE(c.name, '') as city_name,
			COALESCE(u.username, '') as created_by_name
		FROM listing l
		LEFT JOIN city c ON l.city_id = c.id
		LEFT JOIN users u ON l.created_by = u.id`)
	if err != nil {
		return nil, fmt.Errorf("package listing/repo GetAll query: %w", err)
	}

	listings, err := pgx.CollectRows(rows, pgx.RowToStructByName[Listing])

	if err != nil {
		return nil, fmt.Errorf("package listing/repo GetAll: %v", err.Error())
	}

	return listings, nil
}

func (r *Repository) GetAllApproved() ([]Listing, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			l.id, l.name, l.description, l.category, l.item_id, l.city_id, l.created_by, l.created_at, 
			l.approved, l.approved_by, l.approved_at, l.status, l.price, COALESCE(l.image_url, '') as image_url,
			COALESCE(c.name, '') as city_name,
			COALESCE(u.username, '') as created_by_name
		FROM listing l
		LEFT JOIN city c ON l.city_id = c.id
		LEFT JOIN users u ON l.created_by = u.id
		WHERE l.approved = true`)
	if err != nil {
		return nil, fmt.Errorf("package listing/repo GetAllApproved query: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Listing])
}

func (r *Repository) GetByUserId(userId pgtype.Int8) ([]Listing, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			l.id, l.name, l.description, l.category, l.item_id, l.city_id, l.created_by, l.created_at, 
			l.approved, l.approved_by, l.approved_at, l.status, l.price, COALESCE(l.image_url, '') as image_url,
			COALESCE(c.name, '') as city_name,
			COALESCE(u.username, '') as created_by_name
		FROM listing l
		LEFT JOIN city c ON l.city_id = c.id
		LEFT JOIN users u ON l.created_by = u.id
		WHERE l.created_by = $1`, userId)
	if err != nil {
		return nil, fmt.Errorf("package listing/repo GetByUserId query: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Listing])
}

func (r *Repository) GetById(id pgtype.Int8) (*Listing, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			l.id, l.name, l.description, l.category, l.item_id, l.city_id, l.created_by, l.created_at, 
			l.approved, l.approved_by, l.approved_at, l.status, l.price, COALESCE(l.image_url, '') as image_url,
			COALESCE(c.name, '') as city_name,
			COALESCE(u.username, '') as created_by_name
		FROM listing l
		LEFT JOIN city c ON l.city_id = c.id
		LEFT JOIN users u ON l.created_by = u.id
		WHERE l.id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("package listing/repo GetById query: %w", err)
	}

	listing, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Listing])

	if err != nil {
		return nil, fmt.Errorf("package listing/repo GetById: %v", err.Error())
	}
	return &listing, nil
}

func (r *Repository) Create(listingDto Listing) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO listing (name, description, category, item_id, city_id, created_by, price, image_url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		listingDto.Name, listingDto.Description, listingDto.Category, listingDto.ItemId, listingDto.CityId, listingDto.CreatedBy, listingDto.Price, listingDto.ImageUrl).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM listing WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package listing/repo: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) Update(id pgtype.Int8, l Listing) error {
	tag, err := r.db.Exec(db.Ctx,
		"UPDATE listing SET name=$1, description=$2, category=$3, item_id=$4, city_id=$5, created_by=$6, approved=$7, approved_by=$8, approved_at=$9, status=$10, price=$11, image_url=$12 WHERE id=$13",
		l.Name, l.Description, l.Category, l.ItemId, l.CityId, l.CreatedBy, l.Approved, l.ApprovedBy, l.ApprovedAt, l.Status, l.Price, l.ImageUrl, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package listing/repo Update: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) Approve(id pgtype.Int8, adminId pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "UPDATE listing SET approved = true, approved_by = $1, approved_at = NOW() WHERE id = $2", adminId, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package listing/repo Approve: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) Disapprove(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "UPDATE listing SET approved = false, approved_by = NULL, approved_at = NULL WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package listing/repo Disapprove: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64

	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM listing WHERE id = $1", id).Scan(&idFound)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package listing/repo ExistsById query: %w", err)
	}

	return true, nil
}
