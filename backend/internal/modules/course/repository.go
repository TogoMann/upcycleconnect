package course

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

func (r *Repository) GetAll() ([]Course, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, created_at, approved, approved_by, approved_at, price FROM course")
	if err != nil {
		return nil, fmt.Errorf("package course/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Course])
	if err != nil {
		return nil, fmt.Errorf("package course/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Course, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, created_at, approved, approved_by, approved_at, price FROM course WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package course/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Course])
	if err != nil {
		return nil, fmt.Errorf("package course/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto Course) (pgtype.Int8, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO course (created_by, approved, price) VALUES ($1, $2, $3)",
		dto.CreatedBy, dto.Approved, dto.Price)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: tag.RowsAffected(), Valid: true}, err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM course WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package course/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM course WHERE id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package course/repo ExistsById query: %w", err)
	}
	return true, nil
}
