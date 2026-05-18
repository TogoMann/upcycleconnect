package companies

import (
	db "backend/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetBySiret(siret string) (*Company, error) {
	var c Company
	err := r.db.QueryRow(db.Ctx, "SELECT id, siret, name, address, created_at FROM companies WHERE siret = $1", siret).Scan(&c.Id, &c.Siret, &c.Name, &c.Address, &c.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) Create(c Company) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO companies (siret, name, address) VALUES ($1, $2, $3) RETURNING id", c.Siret, c.Name, c.Address).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Company, error) {
	var c Company
	err := r.db.QueryRow(db.Ctx, "SELECT id, siret, name, address, created_at FROM companies WHERE id = $1", id).Scan(&c.Id, &c.Siret, &c.Name, &c.Address, &c.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("company not found")
		}
		return nil, err
	}
	return &c, nil
}
