package post

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

func (r *Repository) GetAll() ([]Post, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, thread_id, created_by, content, upvotes, downvotes, created_at, edited_at FROM post")
	if err != nil {
		return nil, fmt.Errorf("package post/repo GetAll query: %w", err)
	}

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Post])

	if err != nil {
		return nil, fmt.Errorf("package post/repo GetAll: %v", err.Error())
	}

	return posts, nil
}

func (r *Repository) GetById(id int64) (*Post, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, thread_id, created_by, content, upvotes, downvotes, created_at, edited_at FROM post WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package post/repo GetById query: %w", err)
	}

	post, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Post])

	if err != nil {
		return nil, fmt.Errorf("package post/repo GetById: %v", err.Error())
	}
	return &post, nil
}

func (r *Repository) Create(postDto Post) (int64, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO post (thread_id, created_by, content) VALUES ($1, $2, $3, $4)",
		postDto.CreatedBy, postDto.CreatedBy, postDto.Content)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), err
}

func (r *Repository) Delete(id int64) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE post WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package post/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id int64) (bool, error) {
	var idFound int64

	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM post WHERE id = $1", id).Scan(&idFound)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package post/repo ExistsById query: %w", err)
	}

	return true, nil
}
