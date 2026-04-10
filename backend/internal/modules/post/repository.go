package post

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

func (r *Repository) GetAll() ([]Post, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, thread_id, created_by, parent_id, content, upvotes, downvotes, created_at, edited_at FROM post")
	if err != nil {
		return nil, fmt.Errorf("package post/repo GetAll query: %w", err)
	}

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Post])

	if err != nil {
		return nil, fmt.Errorf("package post/repo GetAll: %v", err.Error())
	}

	return posts, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Post, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, thread_id, created_by, parent_id, content, upvotes, downvotes, created_at, edited_at FROM post WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package post/repo GetById query: %w", err)
	}

	post, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Post])

	if err != nil {
		return nil, fmt.Errorf("package post/repo GetById: %v", err.Error())
	}
	return &post, nil
}

func (r *Repository) GetThreadPosts(id pgtype.Int8) ([]ThreadPosts, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT
			p.id, p.thread_id, p.created_by, p.parent_id, p.content, p.upvotes, p.downvotes, p.created_at, p.edited_at,
			t.title, t.content AS thread_content
		FROM post AS p
		INNER JOIN thread AS t
		ON p.thread_id = t.id
		WHERE t.id = $1`, id)

	if err != nil {
		return nil, fmt.Errorf("package post/repo GetThreadPosts query: %w", err)
	}

	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[ThreadPosts])

	if err != nil {
		return nil, fmt.Errorf("package post/repo GetThreadPosts: %v", err.Error())
	}
	return posts, nil
}

func (r *Repository) Create(postDto Post) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO post (thread_id, created_by, parent_id, content) VALUES ($1, $2, $3, $4) RETURNING id",
		postDto.ThreadId, postDto.CreatedBy, postDto.ParentId, postDto.Content).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) UpdateContent(id pgtype.Int8, content string) error {
	tag, err := r.db.Exec(db.Ctx, "UPDATE post SET content = $1, edited_at = CURRENT_TIMESTAMP WHERE id = $2", content, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package post/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM post WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package post/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) SoftDelete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "UPDATE post SET content = '[Deleted post]', created_by = NULL WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package post/repo SoftDelete: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) HasReplies(id pgtype.Int8) (bool, error) {
	var exists bool
	err := r.db.QueryRow(db.Ctx, "SELECT EXISTS(SELECT 1 FROM post WHERE parent_id = $1)", id).Scan(&exists)
	return exists, err
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
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
