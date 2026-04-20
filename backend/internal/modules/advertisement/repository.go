package advertisement

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

func (r *Repository) GetAllPubs() ([]PubFrontend, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			a.id,
			'Campagne #' || a.id as titre,
			u.username as annonceur,
			a.ad_type as type,
			TO_CHAR(a.start_date, 'DD/MM/YYYY') as debut,
			TO_CHAR(a.end_date, 'DD/MM/YYYY') as fin,
			CASE WHEN a.status = 'validated' THEN 'active' ELSE 'inactive' END as statut,
			CAST(COALESCE(a.budget, 0) AS FLOAT8) as budget
		FROM advertisement a
		JOIN users u ON a.announcer_id = u.id
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[PubFrontend])
}

func (r *Repository) GetAll() ([]Advertisement, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT * FROM advertisement")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Advertisement])
}

func (r *Repository) GetById(id pgtype.Int8) (*Advertisement, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT * FROM advertisement WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	ad, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Advertisement])
	if err != nil {
		return nil, err
	}
	return &ad, nil
}

func (r *Repository) Create(ad Advertisement) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, `INSERT INTO advertisement (announcer_id, target_id, target_type, ad_type, budget, status, start_date, end_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, 
		ad.AnnouncerId, ad.TargetId, ad.TargetType, ad.AdType, ad.Budget, ad.Status, ad.StartDate, ad.EndDate).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) UpdateStatus(id pgtype.Int8, status string, approvedBy pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE advertisement SET status = $1, approved_by = $2 WHERE id = $3", status, approvedBy, id)
	return err
}

func (r *Repository) Reject(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE advertisement SET status = 'rejected', approved_by = NULL WHERE id = $1", id)
	return err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM advertisement WHERE id = $1", id)
	return err
}
