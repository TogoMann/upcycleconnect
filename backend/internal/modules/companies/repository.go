package companies

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

func (r *Repository) GetBySiret(siret string) (*Company, error) {
	var c Company
	query := `
		SELECT 
			co.id, co.siret, co.name, co.address_id, co.created_at,
			COALESCE(a.street_number || ' ' || a.street_name || ', ' || ci.zip_code || ' ' || ci.name, '') as address
		FROM companies co
		LEFT JOIN address a ON co.address_id = a.id
		LEFT JOIN city ci ON a.city_id = ci.id
		WHERE co.siret = $1`
	err := r.db.QueryRow(db.Ctx, query, siret).Scan(&c.Id, &c.Siret, &c.Name, &c.AddressId, &c.CreatedAt, &c.Address)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) Create(c Company) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO companies (siret, name, address_id) VALUES ($1, $2, $3) RETURNING id", c.Siret, c.Name, c.AddressId).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Company, error) {
	var c Company
	query := `
		SELECT 
			co.id, co.siret, co.name, co.address_id, co.created_at,
			COALESCE(a.street_number || ' ' || a.street_name || ', ' || ci.zip_code || ' ' || ci.name, '') as address
		FROM companies co
		LEFT JOIN address a ON co.address_id = a.id
		LEFT JOIN city ci ON a.city_id = ci.id
		WHERE co.id = $1`
	err := r.db.QueryRow(db.Ctx, query, id).Scan(&c.Id, &c.Siret, &c.Name, &c.AddressId, &c.CreatedAt, &c.Address)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("company not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) GetAll() ([]Company, error) {
	query := `
		SELECT 
			co.id, co.siret, co.name, co.address_id, co.created_at,
			COALESCE(a.street_number || ' ' || a.street_name || ', ' || ci.zip_code || ' ' || ci.name, '') as address
		FROM companies co
		LEFT JOIN address a ON co.address_id = a.id
		LEFT JOIN city ci ON a.city_id = ci.id
		ORDER BY co.id DESC`
	rows, err := r.db.Query(db.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var c Company
		if err := rows.Scan(&c.Id, &c.Siret, &c.Name, &c.AddressId, &c.CreatedAt, &c.Address); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}
	return companies, nil
}
