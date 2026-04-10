package comments

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

func (r *Repository) GetAll() ([]Comments, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, news_id, created_by, parent_id, content, created_at, upvotes, downvotes FROM comments")
	if err != nil {
		return nil, fmt.Errorf("package comments/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Comments])
	if err != nil {
		return nil, fmt.Errorf("package comments/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Comments, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, news_id, created_by, parent_id, content, created_at, upvotes, downvotes FROM comments WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package comments/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Comments])
	if err != nil {
		return nil, fmt.Errorf("package comments/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto Comments) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO comments (news_id, created_by, parent_id, content) VALUES ($1, $2, $3, $4) RETURNING id",
		dto.NewsId, dto.CreatedBy, dto.ParentId, dto.Content).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package comments/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) SoftDelete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "UPDATE comments SET content = '[Deleted comment]', created_by = NULL WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package comments/repo SoftDelete: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) HasReplies(id pgtype.Int8) (bool, error) {
	var exists bool
	err := r.db.QueryRow(db.Ctx, "SELECT EXISTS(SELECT 1 FROM comments WHERE parent_id = $1)", id).Scan(&exists)
	return exists, err
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM comments WHERE id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package comments/repo ExistsById query: %w", err)
	}
	return true, nil
}
