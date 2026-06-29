package payment_methods

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

func (r *Repository) GetAll() ([]PaymentMethod, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, user_id, stripe_payment_method_id, card_last4, card_brand, card_exp_month, card_exp_year, is_default, created_at FROM payment_methods")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[PaymentMethod])
}

func (r *Repository) GetById(id pgtype.Int8) (*PaymentMethod, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, user_id, stripe_payment_method_id, card_last4, card_brand, card_exp_month, card_exp_year, is_default, created_at FROM payment_methods WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	pm, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[PaymentMethod])
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

func (r *Repository) GetByUserId(userId pgtype.Int8) ([]PaymentMethod, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, user_id, stripe_payment_method_id, card_last4, card_brand, card_exp_month, card_exp_year, is_default, created_at FROM payment_methods WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[PaymentMethod])
}

func (r *Repository) HasPaymentMethod(userId pgtype.Int8) (bool, error) {
	var count int
	err := r.db.QueryRow(db.Ctx, "SELECT COUNT(id) FROM payment_methods WHERE user_id = $1", userId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) Create(pm PaymentMethod) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO payment_methods (user_id, stripe_payment_method_id, card_last4, card_brand, card_exp_month, card_exp_year, is_default) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		pm.UserId, pm.StripePaymentMethodId, pm.CardLast4, pm.CardBrand, pm.CardExpMonth, pm.CardExpYear, pm.IsDefault).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM payment_methods WHERE id = $1", id)
	return err
}

func (r *Repository) DeleteByIdAndUserId(id pgtype.Int8, userId pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM payment_methods WHERE id = $1 AND user_id = $2", id, userId)
	return err
}

func (r *Repository) SetDefault(id pgtype.Int8, userId pgtype.Int8) error {
	tx, err := r.db.Begin(db.Ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(db.Ctx)

	_, err = tx.Exec(db.Ctx, "UPDATE payment_methods SET is_default = FALSE WHERE user_id = $1", userId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(db.Ctx, "UPDATE payment_methods SET is_default = TRUE WHERE id = $1 AND user_id = $2", id, userId)
	if err != nil {
		return err
	}

	return tx.Commit(db.Ctx)
}
