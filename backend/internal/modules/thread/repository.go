package thread

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

func (r *Repository) GetAll() ([]Thread, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, title, content, upvotes, downvotes, created_at, last_post_at FROM thread")
	if err != nil {
		return nil, fmt.Errorf("package thread/repo GetAll query: %w", err)
	}

	threads, err := pgx.CollectRows(rows, pgx.RowToStructByName[Thread])

	if err != nil {
		return nil, fmt.Errorf("package thread/repo GetAll: %v", err.Error())
	}

	return threads, nil
}

func (r *Repository) GetById(id int64) (*Thread, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, title, content, upvotes, downvotes, created_at, last_post_at FROM thread WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package thread/repo GetById query: %w", err)
	}

	thread, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Thread])

	if err != nil {
		return nil, fmt.Errorf("package thread/repo GetById: %v", err.Error())
	}
	return &thread, nil
}

func (r *Repository) Create(threadDto Thread) (int64, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO thread (created_by, title, content, created_at) VALUES ($1, $2, $3, $4)",
		threadDto.CreatedBy, threadDto.Title, threadDto.Content, threadDto.CreatedAt)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), err
}

func (r *Repository) Delete(id int64) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE thread WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package thread/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id int64) (bool, error) {
	var idFound int64

	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM thread WHERE id = $1", id).Scan(&idFound)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package thread/repo ExistsById query: %w", err)
	}

	return true, nil
}
