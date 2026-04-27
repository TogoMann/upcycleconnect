package course

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

func (r *Repository) GetAllCatalogue() ([]OffreFrontend, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			id,
			name as nom,
			'formation' as categorie,
			CAST(COALESCE(price, 0) AS FLOAT8) as prix,
			COALESCE(description, '') as description,
			approved as actif
		FROM course
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[OffreFrontend])
}

func (r *Repository) GetAllApprovedCatalogue() ([]OffreFrontend, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			id,
			name as nom,
			'formation' as categorie,
			CAST(COALESCE(price, 0) AS FLOAT8) as prix,
			COALESCE(description, '') as description,
			approved as actif
		FROM course
		WHERE approved = true
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[OffreFrontend])
}

func (r *Repository) GetAll() ([]Course, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, max_capacity, created_by, created_at, approved, approved_by, approved_at, price, date, start_time, end_time FROM course")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Course])
}

func (r *Repository) GetById(id pgtype.Int8) (*Course, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, max_capacity, created_by, created_at, approved, approved_by, approved_at, price, date, start_time, end_time FROM course WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	course, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Course])
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *Repository) GetUserCourses(userId pgtype.Int8) ([]UserCourse, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT c.id, c.name, c.description, c.max_capacity, c.created_by, c.created_at, c.approved, c.approved_by, c.approved_at, c.price, c.date, c.start_time, c.end_time, co.buyer_id, co.booked_at 
		FROM course c 
		JOIN course_order co ON c.id = co.course_id 
		WHERE co.buyer_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[UserCourse])
}

func (r *Repository) Create(c Course) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO course (name, description, max_capacity, created_by, price, date, start_time, end_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", c.Name, c.Description, c.MaxCapacity, c.CreatedBy, c.Price, c.Date, c.StartTime, c.EndTime).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Update(id pgtype.Int8, c Course) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE course SET name = $1, description = $2, max_capacity = $3, price = $4, date = $5, start_time = $6, end_time = $7 WHERE id = $8", c.Name, c.Description, c.MaxCapacity, c.Price, c.Date, c.StartTime, c.EndTime, id)
	return err
}

func (r *Repository) Approve(id pgtype.Int8, approvedBy pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE course SET approved = true, approved_by = $1, approved_at = NOW() WHERE id = $2", approvedBy, id)
	return err
}

func (r *Repository) Disapprove(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE course SET approved = false, approved_by = NULL, approved_at = NULL WHERE id = $1", id)
	return err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM course WHERE id = $1", id)
	return err
}
