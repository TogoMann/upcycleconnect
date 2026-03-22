package users

import (
	"github.com/jackc/pgx/v5/pgtype"
	db "backend/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]User, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, first_name, last_name, email, password_hash, role, score, created_at FROM users")
	if err != nil {
		return nil, fmt.Errorf("package users/repo GetAllusers query: %w", err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])

	if err != nil {
		return nil, fmt.Errorf("package users/repo GetAllUsers: %v", err.Error())
	}

	return users, nil
}

func (r *Repository) GetById(id pgtype.Int8) (*User, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, first_name, last_name, email, password_hash, role, score, created_at FROM users WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package users/repo GetUserById query: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])

	if err != nil {
		return nil, fmt.Errorf("package users/repo GetUserById: %v", err.Error())
	}
	return &user, nil
}

func (r *Repository) Create(userDto User) (pgtype.Int8, error) {
	tag, err := r.db.Exec(
		db.Ctx,
		"INSERT INTO users (first_name, last_name, email, password_hash, role) VALUES ($1, $2, $3, $4, $5)",
		userDto.FirstName, userDto.LastName, userDto.Email, userDto.PasswordHash, userDto.Role)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: tag.RowsAffected(), Valid: true}, err
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE users WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package users/repo: Id invalide: %d", id)
	}
	return nil
}

func (r *Repository) ExistsById(id pgtype.Int8) (bool, error) {
	var idFound int64

	err := r.db.QueryRow(db.Ctx, "SELECT 1 FROM users WHERE id = $1", id).Scan(&idFound)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("package users/repo ExistsById query: %w", err)
	}

	return true, nil
}
