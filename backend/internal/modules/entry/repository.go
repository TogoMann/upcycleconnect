package entry

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

func (r *Repository) GetAll() ([]Entry, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, created_at, schedule, start, ending FROM entry")
	if err != nil {
		return nil, fmt.Errorf("package entry/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Entry])
	if err != nil {
		return nil, fmt.Errorf("package entry/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Entry, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, created_at, schedule, start, ending FROM entry WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package entry/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Entry])
	if err != nil {
		return nil, fmt.Errorf("package entry/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto Entry) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO entry (created_by, schedule, start, ending) VALUES ($1, $2, $3, $4) RETURNING id",
		dto.CreatedBy, dto.Schedule, dto.Start, dto.Ending).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM entry WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package entry/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM entry WHERE id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package entry/repo ExistsById query: %w", err)
	}
	return true, nil
}
