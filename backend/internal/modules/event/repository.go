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
	rows, err := r.db.Query(db.Ctx, "SELECT id, approved, approved_by, approved_at, price, date, start_time, end_time, COALESCE(location, '') as location, created_by, created_at FROM event WHERE approved = true")
	if err != nil {
		return nil, fmt.Errorf("package event/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Event])
	if err != nil {
		return nil, fmt.Errorf("package event/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetAllFull() ([]EventFull, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT
			e.id, e.approved, e.approved_by, e.approved_at, e.price, e.date, e.start_time, e.end_time,
			COALESCE(e.location, '') AS location, e.created_by, e.created_at,
			c.username AS creator_username, c.email AS creator_email,
			a.username AS approver_username, a.email AS approver_email
		FROM event e
		LEFT JOIN users c ON c.id = e.created_by
		LEFT JOIN users a ON a.id = e.approved_by
	`)
	if err != nil {
		return nil, fmt.Errorf("package event/repo GetAllFull query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[EventFull])
	if err != nil {
		return nil, fmt.Errorf("package event/repo GetAllFull: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Event, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, approved, approved_by, approved_at, price, date, start_time, end_time, COALESCE(location, '') as location, created_by, created_at FROM event WHERE id = $1", id)
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
		"INSERT INTO event (approved, approved_by, approved_at, price, date, start_time, end_time, location, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		dto.Approved, dto.ApprovedBy, dto.ApprovedAt, dto.Price, dto.Date, dto.StartTime, dto.EndTime, dto.Location, dto.CreatedBy).Scan(&id)

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
		"UPDATE event SET approved=$1, approved_by=$2, approved_at=$3, price=$4, date=$5, start_time=$6, end_time=$7, location=$8, created_by=$9 WHERE id=$10",
		e.Approved, e.ApprovedBy, e.ApprovedAt, e.Price, e.Date, e.StartTime, e.EndTime, e.Location, e.CreatedBy, id)
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

func (r *Repository) Disapprove(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx,
		"UPDATE event SET approved = FALSE, approved_by = NULL, approved_at = NULL WHERE id = $1",
		id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package event/repo Disapprove: Id invalide: %d", id.Int64)
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
