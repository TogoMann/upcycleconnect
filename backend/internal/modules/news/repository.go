package news

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

func (r *Repository) GetAll() ([]News, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, title, content, created_at, upvotes, downvotes FROM news")
	if err != nil {
		return nil, fmt.Errorf("package news/repo GetAll query: %w", err)
	}

	news, err := pgx.CollectRows(rows, pgx.RowToStructByName[News])

	if err != nil {
		return nil, fmt.Errorf("package news/repo GetAll: %v", err.Error())
	}

	return news, nil
}

func (r *Repository) GetById(id int64) (*News, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, title, content, created_at, upvotes, downvotes FROM news WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package news/repo GetById query: %w", err)
	}

	news, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[News])

	if err != nil {
		return nil, fmt.Errorf("package news/repo GetById: %v", err.Error())
	}
	return &news, nil
}

func (r *Repository) Create(newsDto News) (int64, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO news (created_by, title, content) VALUES ($1, $2, $3)",
		newsDto.CreatedBy, newsDto.Title, newsDto.Content)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), err
}

func (r *Repository) Delete(id int64) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE news WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package news/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id int64) (bool, error) {
	var idFound int64

	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM news WHERE id = $1", id).Scan(&idFound)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package news/repo ExistsById query: %w", err)
	}

	return true, nil
}
