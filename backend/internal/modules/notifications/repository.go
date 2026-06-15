package notifications

import (
	db "backend/internal/database"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(n Notification) (int64, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, 
		"INSERT INTO notifications (titre, message, cible, envoyes) VALUES ($1, $2, $3, $4) RETURNING id",
		n.Titre, n.Message, n.Cible, n.Envoyes).Scan(&id)
	return id, err
}

func (r *Repository) GetAll() ([]Notification, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, titre, message, cible, TO_CHAR(created_at, 'DD/MM/YYYY HH24:MI') as date, envoyes FROM notifications ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Notification])
}

func (r *Repository) GetEmailsByRole(role string) ([]string, error) {
	var query string
	var args []interface{}

	if role == "tous" {
		query = "SELECT email FROM users WHERE email IS NOT NULL"
	} else {
		query = "SELECT email FROM users WHERE role = $1 AND email IS NOT NULL"
		args = append(args, role)
	}

	rows, err := r.db.Query(db.Ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}
