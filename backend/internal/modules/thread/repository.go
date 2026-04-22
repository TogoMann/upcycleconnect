package thread

import (
	db "backend/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]Thread, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			t.id, t.created_by, t.category, t.title, t.content, t.upvotes, t.downvotes, t.created_at, t.last_post_at,
			u.username, u.email
		FROM thread t
		LEFT JOIN users u ON t.created_by = u.id`)
	if err != nil {
		return nil, fmt.Errorf("package thread/repo GetAll query: %w", err)
	}

	threads, err := pgx.CollectRows(rows, pgx.RowToStructByName[Thread])

	if err != nil {
		return nil, fmt.Errorf("package thread/repo GetAll: %v", err.Error())
	}

	return threads, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Thread, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			t.id, t.created_by, t.category, t.title, t.content, t.upvotes, t.downvotes, t.created_at, t.last_post_at,
			u.username, u.email
		FROM thread t
		LEFT JOIN users u ON t.created_by = u.id
		WHERE t.id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("package thread/repo GetById query: %w", err)
	}

	thread, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Thread])

	if err != nil {
		return nil, fmt.Errorf("package thread/repo GetById: %v", err.Error())
	}
	return &thread, nil
}

func (r *Repository) Create(threadDto Thread) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO thread (created_by, category, title, content) VALUES ($1, $2, $3, $4) RETURNING id",
		threadDto.CreatedBy, threadDto.Category, threadDto.Title, threadDto.Content).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM thread WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package thread/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) Update(id pgtype.Int8, thread Thread) error {
	tag, err := r.db.Exec(db.Ctx,
		"UPDATE thread SET category=$1, title=$2, content=$3 WHERE id=$4",
		thread.Category, thread.Title, thread.Content, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package thread/repo Update: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
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

func (r *Repository) Upvote(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE thread SET upvotes = upvotes + 1 WHERE id = $1", id)
	return err
}

func (r *Repository) Downvote(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE thread SET downvotes = downvotes + 1 WHERE id = $1", id)
	return err
}
