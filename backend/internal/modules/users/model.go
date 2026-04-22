package users

import "github.com/jackc/pgx/v5/pgtype"

type UserRole string

const (
	Client  UserRole = "client"
	Pro     UserRole = "pro"
	Interne UserRole = "interne"
	Admin   UserRole = "admin"
)

type UserFrontend struct {
	Id                 int64    `db:"id" json:"id"`
	Username           string   `db:"username" json:"username"`
	FirstName          string   `db:"first_name" json:"first_name"`
	LastName           string   `db:"last_name" json:"last_name"`
	Email              string   `db:"email" json:"email"`
	Role               UserRole `db:"role" json:"role"`
	LanguagePreference string   `db:"language_preference" json:"language_preference"`
	HasSeenTutorial    bool     `db:"has_seen_tutorial" json:"has_seen_tutorial"`
	CreatedAt          string   `db:"created_at" json:"created_at"`
	Score              int32    `db:"score" json:"score"`
}

type User struct {
	Id                 pgtype.Int8      `db:"id" json:"id"`
	Username           string           `db:"username" json:"username"`
	FirstName          string           `db:"first_name" json:"first_name"`
	LastName           string           `db:"last_name" json:"last_name"`
	Email              string           `db:"email" json:"email"`
	PasswordHash       string           `db:"password_hash" json:"password_hash"`
	Role               UserRole         `db:"role" json:"role"`
	LanguagePreference string           `db:"language_preference" json:"language_preference"`
	HasSeenTutorial    bool             `db:"has_seen_tutorial" json:"has_seen_tutorial"`
	CreatedAt          pgtype.Timestamp `db:"created_at" json:"created_at"`
}

type ScoreHistory struct {
	Id          pgtype.Int8 `db:"id" json:"id"`
	UserId      pgtype.Int8 `db:"user_id" json:"user_id"`
	Points      int32       `db:"points" json:"points"`
	Description string      `db:"description" json:"description"`
	CreatedAt   string      `db:"created_at" json:"created_at"`
}
