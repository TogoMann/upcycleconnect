package city

import (
	db "backend/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]City, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, zip_code FROM city ORDER BY name ASC")
	if err != nil {
		return nil, fmt.Errorf("package city/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[City])
	if err != nil {
		return nil, fmt.Errorf("package city/repo GetAll: %v", err.Error())
	}

	return items, nil
}
