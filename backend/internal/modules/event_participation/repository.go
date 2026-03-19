package eventparticipation

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

func (r *Repository) GetAll() ([]Event, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT event_id, user_id FROM event_participation")
	if err != nil {
		return nil, fmt.Errorf("package event_participation/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Event])
	if err != nil {
		return nil, fmt.Errorf("package event_participation/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id int64) (*Event, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT event_id, user_id FROM event_participation WHERE event_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package event_participation/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Event])
	if err != nil {
		return nil, fmt.Errorf("package event_participation/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto Event) (int64, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO event_participation (event_id, user_id) VALUES ($1, $2)",
		dto.EventId, dto.UserId)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), err
}

func (r *Repository) Delete(id int64) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM event_participation WHERE event_id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package event_participation/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id int64) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM event_participation WHERE event_id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package event_participation/repo ExistsById query: %w", err)
	}
	return true, nil
}
