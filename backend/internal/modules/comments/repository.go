package comments

import (
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *Repository) GetAll() ([]Comments, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, news_id, created_by, content, created_at, upvotes, downvotes FROM comments")
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
	rows, err := r.db.Query(db.Ctx, "SELECT id, news_id, created_by, content, created_at, upvotes, downvotes FROM comments WHERE id = $1", id)
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
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO comments (news_id, created_by, content) VALUES ($1, $2, $3)",
		dto.NewsId, dto.CreatedBy, dto.Content)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: tag.RowsAffected(), Valid: true}, err
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
