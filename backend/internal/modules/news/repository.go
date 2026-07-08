package news

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

func (r *Repository) GetAll(newsType string) ([]NewsFrontend, error) {
	if newsType == "" {
		rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, title, content, type, TO_CHAR(created_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"') as created_at, upvotes, downvotes, COALESCE(status, 'publie') as status, COALESCE(categorie, '') as categorie FROM news ORDER BY created_at DESC")
		if err != nil {
			return nil, fmt.Errorf("package news/repo GetAll query: %w", err)
		}
		return pgx.CollectRows(rows, pgx.RowToStructByName[NewsFrontend])
	}

	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, title, content, type, TO_CHAR(created_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"') as created_at, upvotes, downvotes, COALESCE(status, 'publie') as status, COALESCE(categorie, '') as categorie FROM news WHERE type = $1 ORDER BY created_at DESC", newsType)
	if err != nil {
		return nil, fmt.Errorf("package news/repo GetAll query: %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[NewsFrontend])
}

func (r *Repository) GetAllPublished() ([]NewsFrontend, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, title, content, type, TO_CHAR(created_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"') as created_at, upvotes, downvotes, COALESCE(status, 'publie') as status, COALESCE(categorie, '') as categorie FROM news WHERE status = 'publie' ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("package news/repo GetAllPublished query: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[NewsFrontend])
}

func (r *Repository) GetPublishedByType(newsType string) ([]NewsFrontend, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, title, content, type, TO_CHAR(created_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"') as created_at, upvotes, downvotes, COALESCE(status, 'publie') as status, COALESCE(categorie, '') as categorie FROM news WHERE type = $1 AND status = 'publie' ORDER BY created_at DESC", newsType)
	if err != nil {
		return nil, fmt.Errorf("package news/repo GetPublishedByType query: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[NewsFrontend])
}

func (r *Repository) GetById(id pgtype.Int8) (*News, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, title, content, type, created_at, upvotes, downvotes, COALESCE(status, 'publie') as status, COALESCE(categorie, '') as categorie FROM news WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package news/repo GetById query: %w", err)
	}

	news, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[News])

	if err != nil {
		return nil, fmt.Errorf("package news/repo GetById: %v", err.Error())
	}
	return &news, nil
}

func (r *Repository) Create(newsDto News) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO news (created_by, title, content, status, categorie, type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		newsDto.CreatedBy, newsDto.Title, newsDto.Content, newsDto.Status, newsDto.Categorie, newsDto.Type).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Update(id pgtype.Int8, newsDto News) error {
	tag, err := r.db.Exec(db.Ctx, "UPDATE news SET title = $1, content = $2 WHERE id = $3", newsDto.Title, newsDto.Content, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package news/repo Update: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM news WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package news/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
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
