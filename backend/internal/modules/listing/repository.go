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
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, item_id, created_by, created_at, approved, approved_by, approved_at, status, price FROM listing")
	if err != nil {
		return nil, fmt.Errorf("package listing/repo GetAll query: %w", err)
	}

	listings, err := pgx.CollectRows(rows, pgx.RowToStructByName[Listing])

	if err != nil {
		return nil, fmt.Errorf("package listing/repo GetAll: %v", err.Error())
	}

	return listings, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Listing, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, item_id, created_by, created_at, approved, approved_by, approved_at, status, price FROM listing WHERE id = $1", id)
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
		"INSERT INTO listing (name, description, item_id, created_by, price) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		listingDto.Name, listingDto.Description, listingDto.ItemId, listingDto.CreatedBy, listingDto.Price).Scan(&id)

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
		return fmt.Errorf("package listing/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) Update(id pgtype.Int8, l Listing) error {
	tag, err := r.db.Exec(db.Ctx,
		"UPDATE listing SET name=$1, description=$2, item_id=$3, created_by=$4, approved=$5, approved_by=$6, approved_at=$7, status=$8, price=$9 WHERE id=$10",
		l.Name, l.Description, l.ItemId, l.CreatedBy, l.Approved, l.ApprovedBy, l.ApprovedAt, l.Status, l.Price, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package listing/repo Update: Id invalide: %d", id)
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
