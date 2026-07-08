package course

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type OffreFrontend struct {
	Id           int64       `db:"id" json:"id"`
	Nom          string      `db:"nom" json:"nom"`
	Categorie    string      `db:"categorie" json:"categorie"`
	Prix         float64     `db:"prix" json:"prix"`
	Description  string      `db:"description" json:"description"`
	Actif        bool        `db:"actif" json:"actif"`
	Type         string      `db:"type" json:"type"`
	Date         pgtype.Date `db:"date" json:"date"`
	EndDate      pgtype.Date `db:"end_date" json:"end_date"`
	Organisateur string      `db:"organisateur" json:"organisateur"`
}

type CourseType string

const (
	Presentiel CourseType = "presentiel"
	EnLigne    CourseType = "en_ligne"
)

type Course struct {
	Id                pgtype.Int8      `db:"id" json:"id"`
	Name              string           `db:"name" json:"name"`
	Description       string           `db:"description" json:"description"`
	MaxCapacity       pgtype.Int4      `db:"max_capacity" json:"max_capacity"`
	CreatedBy         pgtype.Int8      `db:"created_by" json:"created_by"`
	CreatedAt         pgtype.Timestamp `db:"created_at" json:"created_at"`
	Approved          bool             `db:"approved" json:"approved"`
	ApprovedBy        pgtype.Int8      `db:"approved_by" json:"approved_by"`
	ApprovedAt        pgtype.Timestamp `db:"approved_at" json:"approved_at"`
	Price             pgtype.Numeric   `db:"price" json:"price"`
	Date              pgtype.Date      `db:"date" json:"date"`
	StartTime         pgtype.Time      `db:"start_time" json:"start_time"`
	EndTime           pgtype.Time      `db:"end_time" json:"end_time"`
	Status            string           `db:"status" json:"status"`
	CorrectionComment pgtype.Text      `db:"correction_comment" json:"correction_comment"`
	Type              CourseType       `db:"type" json:"type"`
	SessionLink       pgtype.Text      `db:"session_link" json:"session_link"`
	EndDate           pgtype.Date      `db:"end_date" json:"end_date"`
	CreatedByName     string           `db:"created_by_name" json:"created_by_name"`
}

type CourseDocument struct {
	Id           pgtype.Int8      `db:"id" json:"id"`
	CourseId     pgtype.Int8      `db:"course_id" json:"course_id"`
	Filename     string           `db:"filename" json:"filename"`
	OriginalName string           `db:"original_name" json:"original_name"`
	UploadedAt   pgtype.Timestamp `db:"uploaded_at" json:"uploaded_at"`
}

type CourseSession struct {
	Id          pgtype.Int8      `db:"id" json:"id"`
	CourseId    pgtype.Int8      `db:"course_id" json:"course_id"`
	SessionDate pgtype.Date      `db:"session_date" json:"session_date"`
	StartTime   pgtype.Time      `db:"start_time" json:"start_time"`
	EndTime     pgtype.Time      `db:"end_time" json:"end_time"`
	CreatedAt   pgtype.Timestamp `db:"created_at" json:"created_at"`
}

type UserCourse struct {
	Course
	BuyerID  pgtype.Int8      `db:"buyer_id" json:"buyer_id"`
	BookedAt pgtype.Timestamp `db:"booked_at" json:"booked_at"`
}

type FormationListItem struct {
	Course
	Categorie string `db:"categorie" json:"categorie"`
	Duree     string `db:"duree" json:"duree"`
	Inscrits  int    `db:"inscrits" json:"inscrits"`
}

type AdminCourseView struct {
	Id                pgtype.Int8    `db:"id" json:"id"`
	Name              string         `db:"name" json:"name"`
	Description       string         `db:"description" json:"description"`
	CreatorName       string         `db:"creator_name" json:"creator_name"`
	MaxCapacity       pgtype.Int4    `db:"max_capacity" json:"max_capacity"`
	Approved          bool           `db:"approved" json:"approved"`
	Status            string         `db:"status" json:"status"`
	CorrectionComment pgtype.Text    `db:"correction_comment" json:"correction_comment"`
	Price             pgtype.Numeric `db:"price" json:"price"`
	Date              pgtype.Date    `db:"date" json:"date"`
	EndDate           pgtype.Date    `db:"end_date" json:"end_date"`
	Type              CourseType     `db:"type" json:"type"`
	SessionLink       pgtype.Text    `db:"session_link" json:"session_link"`
	SessionCount      int            `db:"session_count" json:"session_count"`
}
