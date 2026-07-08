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
	Id                 pgtype.Int8 `db:"id" json:"id"`
	Username           string      `db:"username" json:"username"`
	FirstName          string      `db:"first_name" json:"first_name"`
	LastName           string      `db:"last_name" json:"last_name"`
	Email              string      `db:"email" json:"email"`
	Role               string      `db:"role" json:"role"`
	LanguagePreference string      `db:"language_preference" json:"language_preference"`
	HasSeenTutorial    bool        `db:"has_seen_tutorial" json:"has_seen_tutorial"`
	CreatedAt          string      `db:"created_at" json:"created_at"`
	Score              int32       `db:"score" json:"score"`
	Siret              string      `db:"siret" json:"siret"`
	CompanyId          int64       `db:"company_id" json:"company_id"`
	Plan               string      `db:"plan" json:"plan"`
	PredictedService   string      `db:"predicted_service" json:"predicted_service"`
	Probability        float64     `db:"probability" json:"probability"`
}

type PaginatedUsers struct {
	Data       []UserFrontend `json:"data"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
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
	Siret              pgtype.Text      `db:"-" json:"-"`
	CompanyId          pgtype.Int8      `db:"company_id" json:"company_id"`
	IsBanned           bool             `db:"is_banned" json:"is_banned"`
	BanExpiresAt       pgtype.Timestamp `db:"ban_expires_at" json:"ban_expires_at"`
}

type ScoreHistory struct {
	Id          pgtype.Int8 `db:"id" json:"id"`
	UserId      pgtype.Int8 `db:"user_id" json:"user_id"`
	Points      int32       `db:"points" json:"points"`
	Description string      `db:"description" json:"description"`
	CreatedAt   string      `db:"created_at" json:"created_at"`
}

type QuestProgress struct {
	Description      string `json:"description"`
	Current          int64  `json:"current"`
	Threshold        int    `json:"threshold"`
	WindowDays       int    `json:"window_days"`
	BonusPoints      int32  `json:"bonus_points"`
	BonusDescription string `json:"bonus_description"`
}
