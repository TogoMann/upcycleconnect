package eventparticipation

import (
	db "backend/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]EventParticipation, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT event_id, user_id, stripe_payment_intent_id FROM event_participation")
	if err != nil {
		return nil, fmt.Errorf("package event_participation/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[EventParticipation])
	if err != nil {
		return nil, fmt.Errorf("package event_participation/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*EventParticipation, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT event_id, user_id, stripe_payment_intent_id FROM event_participation WHERE event_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package event_participation/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[EventParticipation])
	if err != nil {
		return nil, fmt.Errorf("package event_participation/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto EventParticipation) (pgtype.Int8, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO event_participation (event_id, user_id, stripe_payment_intent_id) VALUES ($1, $2, $3)",
		dto.EventId, dto.UserId, dto.StripePaymentIntentId)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: tag.RowsAffected(), Valid: true}, err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM event_participation WHERE event_id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package event_participation/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM event_participation WHERE event_id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package event_participation/repo ExistsById query: %w", err)
	}
	return true, nil
}

func (r *Repository) IsEventPremium(eventId int64) (bool, error) {
	var premium bool
	err := r.db.QueryRow(db.Ctx, "SELECT premium FROM event WHERE id = $1", eventId).Scan(&premium)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return premium, nil
}

func (r *Repository) IsUserEligibleForPremiumEvent(userId int64) (bool, error) {
	var plan string
	var role string
	err := r.db.QueryRow(db.Ctx, `
		SELECT COALESCE(p.name, ''), u.role 
		FROM users u 
		LEFT JOIN subscriptions s ON s.user_id = u.id AND s.status = 'active'
		LEFT JOIN plans p ON p.id = s.plan_id
		WHERE u.id = $1
	`, userId).Scan(&plan, &role)
	if err != nil {
		return false, err
	}

	isPro := role == "pro"
	isAdmin := role == "admin"
	isPremiumPlan := plan == "Premium" || plan == "premium"

	return isPro || isAdmin || isPremiumPlan, nil
}
