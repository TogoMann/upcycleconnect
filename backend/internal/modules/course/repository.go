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
			c.id,
			c.name as nom,
			'formation' as categorie,
			CAST(COALESCE(c.price, 0) AS FLOAT8) as prix,
			COALESCE(c.description, '') as description,
			c.approved as actif,
			c.type,
			c.date,
			c.end_date,
			COALESCE(u.first_name || ' ' || u.last_name, '') as organisateur
		FROM course c
		JOIN users u ON c.created_by = u.id
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[OffreFrontend])
}

func (r *Repository) GetAllApprovedCatalogue() ([]OffreFrontend, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT
			c.id,
			c.name as nom,
			'formation' as categorie,
			CAST(COALESCE(c.price, 0) AS FLOAT8) as prix,
			COALESCE(c.description, '') as description,
			c.approved as actif,
			c.type,
			c.date,
			c.end_date,
			COALESCE(u.first_name || ' ' || u.last_name, '') as organisateur
		FROM course c
		JOIN users u ON c.created_by = u.id
		WHERE c.approved = true
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[OffreFrontend])
}

func (r *Repository) GetAllForAdmin() ([]AdminCourseView, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT
			c.id,
			c.name,
			c.description,
			u.first_name || ' ' || u.last_name as creator_name,
			c.max_capacity,
			c.approved,
			c.status,
			c.correction_comment,
			c.price,
			c.date,
			c.end_date,
			c.type,
			c.session_link,
			(SELECT COUNT(*) FROM course_session WHERE course_id = c.id)::int as session_count
		FROM course c
		JOIN users u ON c.created_by = u.id
		ORDER BY c.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[AdminCourseView])
}

func (r *Repository) GetAll() ([]Course, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, max_capacity, created_by, created_at, approved, approved_by, approved_at, price, date, start_time, end_time, status, correction_comment, type, session_link, end_date FROM course")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Course])
}

func (r *Repository) GetById(id pgtype.Int8) (*Course, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT c.id, c.name, c.description, c.max_capacity, c.created_by, c.created_at, c.approved, c.approved_by, c.approved_at,
			c.price, c.date, c.start_time, c.end_time, c.status, c.correction_comment, c.type, c.session_link, c.end_date,
			COALESCE(u.first_name || ' ' || u.last_name, '') as created_by_name
		FROM course c
		JOIN users u ON c.created_by = u.id
		WHERE c.id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	course, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Course])
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *Repository) GetCoursesByCreator(userId pgtype.Int8) ([]Course, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, name, description, max_capacity, created_by, created_at, approved, approved_by, approved_at, price, date, start_time, end_time, status, correction_comment, type, session_link, end_date FROM course WHERE created_by = $1 ORDER BY created_at DESC", userId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[Course])
}

func (r *Repository) GetUserCourses(userId pgtype.Int8) ([]UserCourse, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT c.id, c.name, c.description, c.max_capacity, c.created_by, c.created_at, c.approved, c.approved_by, c.approved_at, c.price, c.date, c.start_time, c.end_time, c.status, c.correction_comment, c.type, c.session_link, c.end_date, co.buyer_id, co.booked_at
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
	err := r.db.QueryRow(db.Ctx, "INSERT INTO course (name, description, max_capacity, created_by, price, date, start_time, end_time, status, approved, type, session_link, end_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id", c.Name, c.Description, c.MaxCapacity, c.CreatedBy, c.Price, c.Date, c.StartTime, c.EndTime, c.Status, c.Approved, c.Type, c.SessionLink, c.EndDate).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Update(id pgtype.Int8, c Course) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE course SET name = $1, description = $2, max_capacity = $3, price = $4, date = $5, start_time = $6, end_time = $7, status = $8, correction_comment = $9, type = $10, session_link = $11, end_date = $12, approved = $13 WHERE id = $14", c.Name, c.Description, c.MaxCapacity, c.Price, c.Date, c.StartTime, c.EndTime, c.Status, c.CorrectionComment, c.Type, c.SessionLink, c.EndDate, c.Approved, id)
	return err
}

func (r *Repository) CreateSession(s CourseSession) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO course_session (course_id, session_date, start_time, end_time) VALUES ($1, $2, $3, $4) RETURNING id", s.CourseId, s.SessionDate, s.StartTime, s.EndTime).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) GetSessionsByCourseId(courseId pgtype.Int8) ([]CourseSession, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, course_id, session_date, start_time, end_time, created_at FROM course_session WHERE course_id = $1 ORDER BY session_date ASC, start_time ASC", courseId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[CourseSession])
}

func (r *Repository) DeleteSessionsByCourseId(courseId pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM course_session WHERE course_id = $1", courseId)
	return err
}

func (r *Repository) IsUserEnrolled(courseId pgtype.Int8, userId pgtype.Int8) (bool, error) {
	var exists bool
	err := r.db.QueryRow(db.Ctx, "SELECT EXISTS(SELECT 1 FROM course_order WHERE course_id = $1 AND buyer_id = $2)", courseId, userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) CreateDocument(doc CourseDocument) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(db.Ctx, "INSERT INTO course_document (course_id, filename, original_name) VALUES ($1, $2, $3) RETURNING id", doc.CourseId, doc.Filename, doc.OriginalName).Scan(&id)
	if err != nil {
		return pgtype.Int8{}, err
	}
	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) GetDocumentsByCourseId(courseId pgtype.Int8) ([]CourseDocument, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, course_id, filename, original_name, uploaded_at FROM course_document WHERE course_id = $1 ORDER BY uploaded_at DESC", courseId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[CourseDocument])
}

func (r *Repository) GetDocumentById(id pgtype.Int8) (*CourseDocument, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, course_id, filename, original_name, uploaded_at FROM course_document WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	doc, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[CourseDocument])
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *Repository) DeleteDocument(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM course_document WHERE id = $1", id)
	return err
}

func (r *Repository) Approve(id pgtype.Int8, approvedBy pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE course SET approved = true, status = 'approved', approved_by = $1, approved_at = NOW(), correction_comment = NULL WHERE id = $2", approvedBy, id)
	return err
}

func (r *Repository) Disapprove(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE course SET approved = false, status = 'rejected', approved_by = NULL, approved_at = NULL WHERE id = $1", id)
	return err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM course WHERE id = $1", id)
	return err
}
