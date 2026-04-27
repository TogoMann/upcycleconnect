package planning

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

func (r *Repository) GetUserPlanning(userId pgtype.Int8) ([]PlanningItem, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT id, 'depot' as type, 'Dépôt d''objet' as title, '' as description, TO_CHAR(schedule, 'YYYY-MM-DD') as date, TO_CHAR(start, 'HH24:MI') as start_time, TO_CHAR(ending, 'HH24:MI') as end_time, '' as location 
		FROM entry 
		WHERE created_by = $1

		UNION ALL

		SELECT c.id, 'workshop' as type, c.name as title, COALESCE(c.description, '') as description, TO_CHAR(c.date, 'YYYY-MM-DD') as date, TO_CHAR(c.start_time, 'HH24:MI') as start_time, TO_CHAR(c.end_time, 'HH24:MI') as end_time, '' as location 
		FROM course c 
		JOIN course_order co ON c.id = co.course_id 
		WHERE co.buyer_id = $1

		UNION ALL

		SELECT e.id, 'event' as type, 'Événement' as title, '' as description, TO_CHAR(e.date, 'YYYY-MM-DD') as date, TO_CHAR(e.start_time, 'HH24:MI') as start_time, TO_CHAR(e.end_time, 'HH24:MI') as end_time, COALESCE(e.location, '') as location 
		FROM event e 
		JOIN event_participation ep ON e.id = ep.event_id 
		WHERE ep.user_id = $1

		UNION ALL

		SELECT id, 'personal' as type, title, COALESCE(description, '') as description, TO_CHAR(date, 'YYYY-MM-DD') as date, TO_CHAR(start_time, 'HH24:MI') as start_time, TO_CHAR(end_time, 'HH24:MI') as end_time, '' as location 
		FROM personal_event 
		WHERE user_id = $1
		
		ORDER BY date ASC, start_time ASC
	`, userId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[PlanningItem])
}

func (r *Repository) CreatePersonalEvent(e PersonalEvent) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO personal_event (user_id, title, description, date, start_time, end_time) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", e.UserId, e.Title, e.Description, e.Date, e.StartTime, e.EndTime).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) DeletePersonalEvent(id pgtype.Int8, userId pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM personal_event WHERE id = $1 AND user_id = $2", id, userId)
	return err
}
