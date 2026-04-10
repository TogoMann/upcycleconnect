package subscriptions

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

func (r *Repository) GetAll() ([]Subscriptions, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, subscriber_id, price, tier, created_at, until FROM subscriptions")
	if err != nil {
		return nil, fmt.Errorf("package subscriptions/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Subscriptions])
	if err != nil {
		return nil, fmt.Errorf("package subscriptions/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Subscriptions, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, subscriber_id, price, tier, created_at, until FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package subscriptions/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Subscriptions])
	if err != nil {
		return nil, fmt.Errorf("package subscriptions/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto Subscriptions) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO subscriptions (subscriber_id, price, tier, until) VALUES ($1, $2, $3, $4) RETURNING id",
		dto.SubscriberId, dto.Price, dto.Tier, dto.Until).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package subscriptions/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM subscriptions WHERE id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package subscriptions/repo ExistsById query: %w", err)
	}
	return true, nil
}
