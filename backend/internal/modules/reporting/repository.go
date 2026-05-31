package reporting

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDepositedItemsCount(ctx context.Context, start, end time.Time) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM item WHERE created_at >= $1 AND created_at <= $2`
	err := r.db.QueryRow(ctx, query, start, end).Scan(&count)
	return count, err
}

func (r *Repository) GetActorStats(ctx context.Context) ([]ActorStats, error) {
	query := `SELECT role, COUNT(*) FROM users GROUP BY role`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []ActorStats
	for rows.Next() {
		var s ActorStats
		if err := rows.Scan(&s.Role, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func (r *Repository) GetPrestationStats(ctx context.Context) ([]PrestationStats, error) {
	query := `
		SELECT 'event' as type, COUNT(*), COALESCE(SUM(price), 0) FROM event e JOIN event_participation ep ON e.id = ep.event_id GROUP BY type
		UNION ALL
		SELECT 'course' as type, COUNT(*), COALESCE(SUM(price), 0) FROM course_order GROUP BY type
		UNION ALL
		SELECT 'listing' as type, COUNT(*), COALESCE(SUM(price), 0) FROM listing_order GROUP BY type
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []PrestationStats
	for rows.Next() {
		var s PrestationStats
		if err := rows.Scan(&s.Type, &s.Count, &s.Total); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func (r *Repository) GetUserPredictions(ctx context.Context, page, limit int) (*PaginatedPredictions, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM user_predictions").Scan(&total)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
			up.user_id, u.username, u.email, up.predicted_service_type, up.probability, up.calculated_at 
		FROM user_predictions up
		JOIN users u ON up.user_id = u.id
		ORDER BY up.probability DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predictions []UserPrediction
	for rows.Next() {
		var p UserPrediction
		if err := rows.Scan(&p.UserID, &p.Username, &p.Email, &p.PredictedServiceType, &p.Probability, &p.CalculatedAt); err != nil {
			return nil, err
		}
		predictions = append(predictions, p)
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return &PaginatedPredictions{
		Data:       predictions,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
