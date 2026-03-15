package users

import (
	"time"
)

type UserRole string

const (
	Client UserRole = "client"
	Pro    UserRole = "pro"
	Intra  UserRole = "intra"
	Admin  UserRole = "admin"
)

type User struct {
	Id           int64     `db:"id" json:"id"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	Role         UserRole  `db:"role" json:"role"`
	Score        int64     `db:"score" json:"score"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
