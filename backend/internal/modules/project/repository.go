package project

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

func (r *Repository) GetAll() ([]Project, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, listing_id, creator_id, title, description, final_score, status, created_at, completed_at FROM project")
	if err != nil {
		return nil, fmt.Errorf("package project/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Project])
	if err != nil {
		return nil, fmt.Errorf("package project/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Project, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, listing_id, creator_id, title, description, final_score, status, created_at, completed_at FROM project WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package project/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Project])
	if err != nil {
		return nil, fmt.Errorf("package project/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto Project) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO project (listing_id, creator_id, title, description, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		dto.ListingId, dto.CreatorId, dto.Title, dto.Description, dto.Status).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM project WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package project/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM project WHERE id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package project/repo ExistsById query: %w", err)
	}
	return true, nil
}

func (r *Repository) GetSteps(projectId pgtype.Int8) ([]ProjectStep, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, project_id, step_number, description, created_at FROM project_steps WHERE project_id = $1 ORDER BY step_number", projectId)
	if err != nil {
		return nil, fmt.Errorf("package project/repo GetSteps query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProjectStep])
	if err != nil {
		return nil, fmt.Errorf("package project/repo GetSteps: %v", err.Error())
	}

	return items, nil
}
