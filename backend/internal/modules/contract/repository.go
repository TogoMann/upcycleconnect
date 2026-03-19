package contract

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

func (r *Repository) GetById(id int64) (*Contract, error) {
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

func (r *Repository) Create(dto Contract) (int64, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO contract (name, created_by, content, until) VALUES ($1, $2, $3, $4)",
		dto.Name, dto.CreatedBy, dto.Content, dto.Until)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), err
}

func (r *Repository) Delete(id int64) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM contract WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package contract/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id int64) (bool, error) {
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
