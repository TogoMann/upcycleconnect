package entryparticipation

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

func (r *Repository) GetById(id pgtype.Int8) (*EntryParticipation, error) {
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

func (r *Repository) Create(dto EntryParticipation) (pgtype.Int8, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO entry_participation (entry_id, user_id, status) VALUES ($1, $2, $3)",
		dto.EntryId, dto.UserId, dto.Status)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: tag.RowsAffected(), Valid: true}, err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM entry_participation WHERE entry_id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package entryparticipation/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) GetByUserId(userId pgtype.Int8) ([]EntryParticipation, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT entry_id, user_id, status, joined_at FROM entry_participation WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[EntryParticipation])
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
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
