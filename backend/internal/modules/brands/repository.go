package brands

import (
	db "backend/internal/database"

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

func (r *Repository) GetAll() ([]Brand, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, logo_url, website, created_by, is_active, created_at FROM brands")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Brand])
}

func (r *Repository) GetById(id pgtype.Int8) (*Brand, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, logo_url, website, created_by, is_active, created_at FROM brands WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	b, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Brand])
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) GetByCreator(userId pgtype.Int8) ([]Brand, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, logo_url, website, created_by, is_active, created_at FROM brands WHERE created_by = $1", userId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Brand])
}

func (r *Repository) Create(b Brand) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO brands (name, description, logo_url, website, created_by, is_active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		b.Name, b.Description, b.LogoUrl, b.Website, b.CreatedBy, b.IsActive).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Update(id pgtype.Int8, b Brand) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE brands SET name = $1, description = $2, logo_url = $3, website = $4, is_active = $5 WHERE id = $6",
		b.Name, b.Description, b.LogoUrl, b.Website, b.IsActive, id)
	return err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM brands WHERE id = $1", id)
	return err
}
