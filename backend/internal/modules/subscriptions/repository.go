package subscriptions

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

func (r *Repository) GetAllAbonnements() ([]AbonnementFrontend, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			s.id,
			'Upcycle Corp' as entreprise,
			u.first_name || ' ' || u.last_name as utilisateur,
			LOWER(s.tier) as plan,
			TO_CHAR(s.created_at, 'DD/MM/YYYY') as debut,
			TO_CHAR(s.until, 'DD/MM/YYYY') as fin,
			CASE WHEN s.until >= CURRENT_DATE THEN 'actif' ELSE 'expiré' END as statut,
			CAST(s.price AS FLOAT8) as montant
		FROM subscriptions s
		JOIN users u ON s.subscriber_id = u.id
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[AbonnementFrontend])
}

func (r *Repository) GetAll() ([]Subscription, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, subscriber_id, price, tier, created_at, until FROM subscriptions")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Subscription])
}

func (r *Repository) GetById(id pgtype.Int8) (*Subscription, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, subscriber_id, price, tier, created_at, until FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	sub, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Subscription])
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *Repository) Create(s Subscription) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO subscriptions (subscriber_id, price, tier, until) VALUES ($1, $2, $3, $4) RETURNING id", s.SubscriberId, s.Price, s.Tier, s.Until).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Update(id pgtype.Int8, s Subscription) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE subscriptions SET price = $1, tier = $2, until = $3 WHERE id = $4", s.Price, s.Tier, s.Until, id)
	return err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM subscriptions WHERE id = $1", id)
	return err
}
