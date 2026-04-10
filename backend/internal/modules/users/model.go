package users

import "github.com/jackc/pgx/v5/pgtype"

type UserRole string

const (
	Client  UserRole = "client"
	Pro     UserRole = "pro"
	Interne UserRole = "interne"
	Admin   UserRole = "admin"
)

type User struct {
	Id                 pgtype.Int8      `db:"id" json:"id"`
	Username           string           `db:"username" json:"username"`
	FirstName          string           `db:"first_name" json:"first_name"`
	LastName           string           `db:"last_name" json:"last_name"`
	Email              string           `db:"email" json:"email"`
	PasswordHash       string           `db:"password_hash" json:"password_hash"`
	Role               UserRole         `db:"role" json:"role"`
	LanguagePreference string           `db:"language_preference" json:"language_preference"`
	Score              pgtype.Int4      `db:"score" json:"score"`
	CreatedAt          pgtype.Timestamp `db:"created_at" json:"created_at"`
}
