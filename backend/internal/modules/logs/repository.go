package logs

import (
	db "backend/internal/database"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(utilisateur, action, ressource, ip, niveau string) error {
	_, err := r.db.Exec(db.Ctx, `
		INSERT INTO logs (utilisateur, action, ressource, ip, niveau)
		VALUES ($1, $2, $3, $4, $5)
	`, utilisateur, action, ressource, ip, niveau)
	return err
}

func (r *Repository) GetAll() ([]Log, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT id, utilisateur, action, ressource, ip, TO_CHAR(created_at, 'DD/MM/YYYY HH24:MI') as date, niveau
		FROM logs
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("GetAll logs query: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Log])
}
