package courseorder

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

func (r *Repository) GetAll() ([]CourseOrder, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, course_id, buyer_id, stripe_payment_intent_id, price, booked_at FROM course_order")
	if err != nil {
		return nil, fmt.Errorf("package courseorder/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[CourseOrder])
	if err != nil {
		return nil, fmt.Errorf("package courseorder/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*CourseOrder, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, course_id, buyer_id, stripe_payment_intent_id, price, booked_at FROM course_order WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package courseorder/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CourseOrder])
	if err != nil {
		return nil, fmt.Errorf("package courseorder/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto CourseOrder) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO course_order (course_id, buyer_id, stripe_payment_intent_id, price) VALUES ($1, $2, $3, $4) RETURNING id",
		dto.CourseId, dto.BuyerId, dto.StripePaymentIntentId, dto.Price).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM course_order WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package courseorder/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM course_order WHERE id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package courseorder/repo ExistsById query: %w", err)
	}
	return true, nil
}
