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

func (r *Repository) GetProStats(userId int64) (*ProStats, error) {
	var s ProStats

	err := r.db.QueryRow(db.Ctx, "SELECT COUNT(*) FROM listing WHERE created_by = $1 AND status = 'active'", userId).Scan(&s.Annonces)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(db.Ctx, "SELECT COUNT(*) FROM project WHERE created_by = $1", userId).Scan(&s.Projets)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(db.Ctx, "SELECT CAST(COALESCE(SUM(points), 0) AS INTEGER) FROM score_history WHERE user_id = $1", userId).Scan(&s.Score)
	if err != nil {
		return nil, err
	}

	s.Vues = s.Annonces*42 + s.Projets*128 + 15

	return &s, nil
}
