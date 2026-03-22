package listing

import (
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *Repository) GetAll() ([]Listing, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, created_at, approved, approved_by, approved_at, status, price FROM listing")
	if err != nil {
		return nil, fmt.Errorf("package listing/repo GetAll query: %w", err)
	}

	listings, err := pgx.CollectRows(rows, pgx.RowToStructByName[Listing])

	if err != nil {
		return nil, fmt.Errorf("package news/repo GetAll: %v", err.Error())
	}

	return listings, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Listing, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, created_at, approved, approved_by, approved_at, status, price  FROM listing WHERE id = $1", id)
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
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO listing (created_by, price) VALUES ($1, $2)",
		listingDto.CreatedBy, listingDto.Price)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: tag.RowsAffected(), Valid: true}, err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE listing WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package listing/repo: Id invalide: %d", id)
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
