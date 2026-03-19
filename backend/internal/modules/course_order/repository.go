package courseorder

import (
	db "backend/internal/database"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]CourseOrder, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, course_id, buyer_id, price, booked_at FROM course_order")
	if err != nil {
		return nil, fmt.Errorf("package courseorder/repo GetAll query: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[CourseOrder])
	if err != nil {
		return nil, fmt.Errorf("package courseorder/repo GetAll: %v", err.Error())
	}

	return items, nil
}

func (r *Repository) GetById(id int64) (*CourseOrder, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, course_id, buyer_id, price, booked_at FROM course_order WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package courseorder/repo GetById query: %w", err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CourseOrder])
	if err != nil {
		return nil, fmt.Errorf("package courseorder/repo GetById: %v", err.Error())
	}
	return &item, nil
}

func (r *Repository) Create(dto CourseOrder) (int64, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO course_order (course_id, buyer_id, price) VALUES ($1, $2, $3)",
		dto.CourseId, dto.BuyerId, dto.Price)

	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), err
}

func (r *Repository) Delete(id int64) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM course_order WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package courseorder/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id int64) (bool, error) {
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
