package event

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

func (r *Repository) GetAll() ([]Event, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, approved, approved_by, approved_at, price, date, created_by, created_at FROM event")
	if err != nil {
		return nil, fmt.Errorf("package event/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Event])
	if err != nil {
		return nil, fmt.Errorf("package event/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Event, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, approved, approved_by, approved_at, price, date, created_by, created_at FROM event WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package event/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Event])
	if err != nil {
		return nil, fmt.Errorf("package event/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto Event) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO event (approved, price, date, created_by) VALUES ($1, $2, $3, $4) RETURNING id",
		dto.Approved, dto.Price, dto.Date, dto.CreatedBy).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM event WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package event/repo: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) Update(id pgtype.Int8, e Event) error {
	tag, err := r.db.Exec(db.Ctx,
		"UPDATE event SET approved=$1, approved_by=$2, approved_at=$3, price=$4, date=$5, created_by=$6 WHERE id=$7",
		e.Approved, e.ApprovedBy, e.ApprovedAt, e.Price, e.Date, e.CreatedBy, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package event/repo Update: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) Approve(id pgtype.Int8, adminId pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx,
		"UPDATE event SET approved = TRUE, approved_by = $1, approved_at = NOW() WHERE id = $2",
		adminId, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package event/repo Approve: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM event WHERE id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package event/repo ExistsById query: %w", err)
	}
	return true, nil
}
