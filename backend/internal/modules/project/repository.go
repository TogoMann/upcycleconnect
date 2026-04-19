package project

import (
	db "backend/internal/database"
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

func (r *Repository) GetAllProjets() ([]ProjetFrontend, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			p.id,
			p.title as nom,
			'Upcycling' as type,
			u.username as auteur,
			CASE WHEN p.status = 'done' THEN 'termine' ELSE 'en_cours' END as statut,
			TO_CHAR(p.created_at, 'DD/MM/YYYY') as date,
			CAST((CASE WHEN p.status = 'featured' THEN true ELSE false END) AS BOOLEAN) as mis_en_avant
		FROM project p
		JOIN users u ON p.creator_id = u.id
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[ProjetFrontend])
}

func (r *Repository) GetAll() ([]Project, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT * FROM project")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Project])
}

func (r *Repository) GetById(id pgtype.Int8) (*Project, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT * FROM project WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	p, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Project])
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) Create(p Project) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO project (listing_id, creator_id, title, description) VALUES ($1, $2, $3, $4) RETURNING id", p.ListingId, p.CreatorId, p.Title, p.Description).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Update(id pgtype.Int8, p Project) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE project SET title = $1, description = $2, status = $3 WHERE id = $4", p.Title, p.Description, p.Status, id)
	return err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM project WHERE id = $1", id)
	return err
}

func (r *Repository) UpdateFeatured(id int64, featured bool) error {
	status := "in progress"
	if featured {
		status = "featured"
	}
	_, err := r.db.Exec(db.Ctx, "UPDATE project SET status = $1 WHERE id = $2", status, id)
	return err
}

func (r *Repository) GetSteps(projectId pgtype.Int8) ([]Step, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT * FROM project_steps WHERE project_id = $1 ORDER BY step_number", projectId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Step])
}

func (r *Repository) CreateStep(s Step) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO project_steps (project_id, step_number, description) VALUES ($1, $2, $3) RETURNING id", s.ProjectId, s.StepNumber, s.Description).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) UpdateStep(id int64, s Step) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE project_steps SET step_number = $1, description = $2 WHERE id = $3", s.StepNumber, s.Description, id)
	return err
}

func (r *Repository) DeleteStep(id int64) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM project_steps WHERE id = $1", id)
	return err
}
