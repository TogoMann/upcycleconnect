package plans

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

func (r *Repository) GetAll() ([]Plan, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, price, billing_cycle, features, is_active, created_at FROM plans ORDER BY price ASC")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Plan])
}

func (r *Repository) GetById(id pgtype.Int8) (*Plan, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, price, billing_cycle, features, is_active, created_at FROM plans WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	plan, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Plan])
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *Repository) Create(p Plan) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO plans (name, description, price, billing_cycle, features, is_active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", p.Name, p.Description, p.Price, p.BillingCycle, p.Features, p.IsActive).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Update(id pgtype.Int8, p Plan) error {
	tag, err := r.db.Exec(db.Ctx, "UPDATE plans SET name = $1, description = $2, price = $3, billing_cycle = $4, features = $5, is_active = $6 WHERE id = $7", p.Name, p.Description, p.Price, p.BillingCycle, p.Features, p.IsActive, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("plan not found")
	}
	return nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM plans WHERE id = $1", id)
	return err
}
