package contract

import (
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *Repository) GetAll() ([]Contract, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, created_by, content, created_at, until FROM contract")
	if err != nil {
		return nil, fmt.Errorf("package contract/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Contract])
	if err != nil {
		return nil, fmt.Errorf("package contract/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Contract, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, created_by, content, created_at, until FROM contract WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package contract/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Contract])
	if err != nil {
		return nil, fmt.Errorf("package contract/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto Contract) (pgtype.Int8, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO contract (name, created_by, content, until) VALUES ($1, $2, $3, $4)",
		dto.Name, dto.CreatedBy, dto.Content, dto.Until)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: tag.RowsAffected(), Valid: true}, err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM contract WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package contract/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM contract WHERE id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package contract/repo ExistsById query: %w", err)
	}
	return true, nil
}
