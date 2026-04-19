package container

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

func (r *Repository) GetAll() ([]ConteneurFrontend, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			c.id, 
			'CONT-' || c.id as code_barres,
			a.street_number || ' ' || a.street_name || ', ' || ci.name as localisation,
			CASE 
				WHEN c.status = 'Available' THEN 'actif'
				WHEN c.status = 'Occupied' THEN 'plein'
				ELSE 'hs'
			END as etat,
			CASE 
				WHEN c.size = 'S' THEN 10
				WHEN c.size = 'M' THEN 20
				WHEN c.size = 'L' THEN 50
				ELSE 20
			END as capacite,
			CAST((SELECT COUNT(*) FROM item i WHERE i.container_id = c.id) AS INTEGER) as objets
		FROM container c
		JOIN site s ON c.site_id = s.id
		JOIN address a ON s.address_id = a.id
		JOIN city ci ON a.city_id = ci.id
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[ConteneurFrontend])
}

func (r *Repository) GetById(id pgtype.Int8) (*Container, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, site_id, status, size, created_at FROM container WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	container, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Container])
	if err != nil {
		return nil, err
	}
	return &container, nil
}

func (r *Repository) Create(c Container) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO container (site_id, status, size) VALUES ($1, $2, $3) RETURNING id", c.SiteId, c.Status, c.Size).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Update(id pgtype.Int8, c Container) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE container SET site_id = $1, status = $2, size = $3 WHERE id = $4", c.SiteId, c.Status, c.Size, id)
	return err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM container WHERE id = $1", id)
	return err
}
