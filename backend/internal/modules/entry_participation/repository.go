package entryparticipation

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

func (r *Repository) GetAll() ([]EntryParticipation, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT entry_id, user_id, status, joined_at FROM entry_participation")
	if err != nil {
		return nil, fmt.Errorf("package entryparticipation/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[EntryParticipation])
	if err != nil {
		return nil, fmt.Errorf("package entryparticipation/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id int64) (*EntryParticipation, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT entry_id, user_id, status, joined_at FROM entry_participation WHERE entry_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package entryparticipation/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[EntryParticipation])
	if err != nil {
		return nil, fmt.Errorf("package entryparticipation/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto EntryParticipation) (int64, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO entry_participation (entry_id, user_id, status) VALUES ($1, $2, $3)",
		dto.EntryId, dto.UserId, dto.Status)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), err
}

func (r *Repository) Delete(id int64) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM entry_participation WHERE entry_id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package entryparticipation/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id int64) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM entry_participation WHERE entry_id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package entryparticipation/repo ExistsById query: %w", err)
	}
	return true, nil
}
