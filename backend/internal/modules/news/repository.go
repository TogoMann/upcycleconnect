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

func orderClauseForSort(sortBy string) string {
	switch sortBy {
	case "top":
		return "ORDER BY (n.upvotes - n.downvotes) DESC, n.created_at DESC"
	case "bottom":
		return "ORDER BY (n.upvotes - n.downvotes) ASC, n.created_at DESC"
	default:
		return "ORDER BY n.created_at DESC"
	}
}

func (r *Repository) GetAll(newsType string, userId pgtype.Int8, sortBy string) ([]NewsFrontend, error) {
	baseQuery := `
		SELECT n.id, n.created_by, n.title, n.content, n.type,
			TO_CHAR(n.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
			n.upvotes, n.downvotes, nv.vote_type as my_vote
		FROM news n
		LEFT JOIN news_vote nv ON nv.news_id = n.id AND nv.user_id = $1
	`
	order := orderClauseForSort(sortBy)

	if newsType == "" {
		rows, err := r.db.Query(db.Ctx, baseQuery+order, userId)
		if err != nil {
			return nil, fmt.Errorf("package news/repo GetAll query: %w", err)
		}
		return pgx.CollectRows(rows, pgx.RowToStructByName[NewsFrontend])
	}

	rows, err := r.db.Query(db.Ctx, baseQuery+" WHERE n.type = $2 "+order, userId, newsType)
	if err != nil {
		return nil, fmt.Errorf("package news/repo GetAll query: %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[NewsFrontend])
}

func (r *Repository) GetById(id pgtype.Int8) (*News, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, title, content, type, created_at, upvotes, downvotes FROM news WHERE id = $1", id)
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
		"INSERT INTO news (created_by, title, content, type) VALUES ($1, $2, $3, $4) RETURNING id",
		newsDto.CreatedBy, newsDto.Title, newsDto.Content, newsDto.Type).Scan(&id)

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

func (r *Repository) Vote(newsId pgtype.Int8, userId pgtype.Int8, voteType string) (*NewsFrontend, error) {
	tx, err := r.db.Begin(db.Ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(db.Ctx)

	var existingType string
	err = tx.QueryRow(db.Ctx, "SELECT vote_type FROM news_vote WHERE news_id = $1 AND user_id = $2 FOR UPDATE", newsId, userId).Scan(&existingType)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	hasExisting := err == nil

	switch {
	case !hasExisting:
		if _, err := tx.Exec(db.Ctx, "INSERT INTO news_vote (news_id, user_id, vote_type) VALUES ($1, $2, $3)", newsId, userId, voteType); err != nil {
			return nil, err
		}
		if voteType == "up" {
			_, err = tx.Exec(db.Ctx, "UPDATE news SET upvotes = upvotes + 1 WHERE id = $1", newsId)
		} else {
			_, err = tx.Exec(db.Ctx, "UPDATE news SET downvotes = downvotes + 1 WHERE id = $1", newsId)
		}
		if err != nil {
			return nil, err
		}
	case existingType == voteType:
		if _, err := tx.Exec(db.Ctx, "DELETE FROM news_vote WHERE news_id = $1 AND user_id = $2", newsId, userId); err != nil {
			return nil, err
		}
		if voteType == "up" {
			_, err = tx.Exec(db.Ctx, "UPDATE news SET upvotes = upvotes - 1 WHERE id = $1", newsId)
		} else {
			_, err = tx.Exec(db.Ctx, "UPDATE news SET downvotes = downvotes - 1 WHERE id = $1", newsId)
		}
		if err != nil {
			return nil, err
		}
	default:
		if _, err := tx.Exec(db.Ctx, "UPDATE news_vote SET vote_type = $3, created_at = CURRENT_TIMESTAMP WHERE news_id = $1 AND user_id = $2", newsId, userId, voteType); err != nil {
			return nil, err
		}
		if voteType == "up" {
			_, err = tx.Exec(db.Ctx, "UPDATE news SET upvotes = upvotes + 1, downvotes = downvotes - 1 WHERE id = $1", newsId)
		} else {
			_, err = tx.Exec(db.Ctx, "UPDATE news SET downvotes = downvotes + 1, upvotes = upvotes - 1 WHERE id = $1", newsId)
		}
		if err != nil {
			return nil, err
		}
	}

	rows, err := tx.Query(db.Ctx, `
		SELECT n.id, n.created_by, n.title, n.content, n.type,
			TO_CHAR(n.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
			n.upvotes, n.downvotes, nv.vote_type as my_vote
		FROM news n
		LEFT JOIN news_vote nv ON nv.news_id = n.id AND nv.user_id = $2
		WHERE n.id = $1
	`, newsId, userId)
	if err != nil {
		return nil, err
	}
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[NewsFrontend])
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(db.Ctx); err != nil {
		return nil, err
	}

	return &updated, nil
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
