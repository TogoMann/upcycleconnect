package course

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

func (r *Repository) GetAll() ([]Course, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, max_capacity, created_by, created_at, approved, approved_by, approved_at, price FROM course")
	if err != nil {
		return nil, fmt.Errorf("package course/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Course])
	if err != nil {
		return nil, fmt.Errorf("package course/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*Course, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, max_capacity, created_by, created_at, approved, approved_by, approved_at, price FROM course WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package course/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Course])
	if err != nil {
		return nil, fmt.Errorf("package course/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) GetUserCourses(id pgtype.Int8) ([]UserCourse, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT
			c.id, c.name, c.description, c.max_capacity, c.created_by, c.created_at, c.approved, c.approved_by, c.approved_at, c.price,
		    co.buyer_id, co.booked_at
			FROM course AS c
			INNER JOIN course_order AS co
			ON co.course_id = c.id
			WHERE buyer_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("package course/repo GetUserCourses query: %w", err)
	}

	courses, err := pgx.CollectRows(rows, pgx.RowToStructByName[UserCourse])
	if err != nil {
		return nil, fmt.Errorf("package course/repo GetUserCourses: %v", err.Error())
	}
	return courses, nil
}

func (r *Repository) Create(dto Course) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO course (name, description, max_capacity, created_by, approved, price) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		dto.Name, dto.Description, dto.MaxCapacity, dto.CreatedBy, dto.Approved, dto.Price).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM course WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package course/repo: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) Update(id pgtype.Int8, c Course) error {
	tag, err := r.db.Exec(db.Ctx,
		"UPDATE course SET name=$1, description=$2, max_capacity=$3, approved=$4, approved_by=$5, approved_at=$6, price=$7 WHERE id=$8",
		c.Name, c.Description, c.MaxCapacity, c.Approved, c.ApprovedBy, c.ApprovedAt, c.Price, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package course/repo Update: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) Approve(id pgtype.Int8, adminId pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx,
		"UPDATE course SET approved = TRUE, approved_by = $1, approved_at = NOW() WHERE id = $2",
		adminId, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package course/repo Approve: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64
	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM course WHERE id = $1", id).Scan(&idFound)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package course/repo ExistsById query: %w", err)
	}
	return true, nil
}
