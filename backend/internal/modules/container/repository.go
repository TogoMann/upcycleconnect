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
				WHEN c.status = 'HS' THEN 'hs'
				WHEN NOT EXISTS (SELECT 1 FROM locker l WHERE l.container_id = c.id AND l.status = 'Available') THEN 'plein'
				ELSE 'actif'
			END as etat,
			CAST((
				SELECT COUNT(*)
				FROM locker l 
				WHERE l.container_id = c.id
			) AS INTEGER) as capacite,
			CAST((
				SELECT COUNT(*) 
				FROM locker l 
				WHERE l.container_id = c.id AND l.status != 'Available'
			) AS INTEGER) as objets
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
	rows, err := r.db.Query(db.Ctx, "SELECT id, site_id, status, created_at FROM container WHERE id = $1", id)
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
	err := r.db.QueryRow(db.Ctx, "INSERT INTO container (site_id, status) VALUES ($1, $2) RETURNING id", c.SiteId, c.Status).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Update(id pgtype.Int8, c Container) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE container SET site_id = $1, status = $2 WHERE id = $3", c.SiteId, c.Status, id)
	return err
}

func (r *Repository) GetLockersByContainerId(containerId pgtype.Int8) ([]Locker, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, container_id, label, status, size, created_at FROM locker WHERE container_id = $1 ORDER BY label", containerId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Locker])
}

func (r *Repository) CreateLocker(l Locker) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO locker (container_id, label, status, size) VALUES ($1, $2, $3, $4) RETURNING id", l.ContainerId, l.Label, l.Status, l.Size).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) UpdateLocker(id pgtype.Int8, l Locker) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE locker SET label = $1, status = $2, size = $3 WHERE id = $4", l.Label, l.Status, l.Size, id)
	return err
}

func (r *Repository) DeleteLocker(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM locker WHERE id = $1", id)
	return err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM container WHERE id = $1", id)
	return err
}
