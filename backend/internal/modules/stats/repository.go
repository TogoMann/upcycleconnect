package stats

import (
	db "backend/internal/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPublicStats() (*PublicStats, error) {
	var s PublicStats

	err := r.db.QueryRow(db.Ctx, "SELECT COUNT(*) FROM users WHERE role != 'admin'").Scan(&s.ActiveMembers)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(db.Ctx, "SELECT COUNT(*) FROM item WHERE status = 'collected'").Scan(&s.ItemsRenewed)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(db.Ctx, "SELECT COUNT(*) FROM companies").Scan(&s.PartnerArtisans)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(db.Ctx, `
		SELECT COUNT(DISTINCT ci.id)
		FROM site s
		JOIN address a ON s.address_id = a.id
		JOIN city ci ON a.city_id = ci.id
	`).Scan(&s.RegionsCovered)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
