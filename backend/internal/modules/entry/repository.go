package entry

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

func (r *Repository) GetAllDepots() ([]DepotFrontend, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			ep.entry_id as id,
			u.first_name || ' ' || u.last_name as utilisateur,
			COALESCE(i.material_type, 'Objet') as objet,
			TO_CHAR(ep.joined_at, 'DD/MM/YYYY') as date,
			CASE WHEN ep.status = 'accepted' THEN 'valide' ELSE 'en_attente' END as statut,
			FALSE as code_envoye
		FROM entry_participation ep
		JOIN users u ON ep.user_id = u.id
		LEFT JOIN item i ON i.owner_id = u.id
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[DepotFrontend])
}

func (r *Repository) GetAll() ([]Entry, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, created_at, schedule, start, ending FROM entry")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Entry])
}

func (r *Repository) GetById(id pgtype.Int8) (*Entry, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, created_by, created_at, schedule, start, ending FROM entry WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	entry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Entry])
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *Repository) Create(e Entry) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO entry (created_by, schedule, start, ending) VALUES ($1, $2, $3, $4) RETURNING id", e.CreatedBy, e.Schedule, e.Start, e.Ending).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM entry WHERE id = $1", id)
	return err
}

func (r *Repository) ValidateDepot(id int64) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE entry_participation SET status = 'accepted' WHERE entry_id = $1", id)
	return err
}
