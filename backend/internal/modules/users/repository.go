package users

import (
	db "backend/internal/database"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]UserFrontend, error) {
	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			id, username, first_name, last_name, email, role, language_preference, has_seen_tutorial, 
			TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
			CAST((SELECT COALESCE(SUM(points), 0) FROM score_history WHERE user_id = users.id) AS INTEGER) as score
		FROM users
	`)
	if err != nil {
		return nil, fmt.Errorf("package users/repo GetAllusers query: %w", err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[UserFrontend])
}

func (r *Repository) GetById(id pgtype.Int8) (*User, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, username, first_name, last_name, email, password_hash, role, language_preference, has_seen_tutorial, created_at FROM users WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("package users/repo GetUserById query: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])

	if err != nil {
		return nil, fmt.Errorf("package users/repo GetUserById: %v", err.Error())
	}
	return &user, nil
}

func (r *Repository) GetByUsername(username string) (*User, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, username, first_name, last_name, email, password_hash, role, language_preference, has_seen_tutorial, created_at FROM users WHERE username = $1", username)
	if err != nil {
		return nil, fmt.Errorf("package users/repo GetByUsername query: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("package users/repo GetByUsername: %v", err.Error())
	}
	return &user, nil
}

func (r *Repository) UpdateTutorialSeen(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "UPDATE users SET has_seen_tutorial = TRUE WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *Repository) Create(userDto User) (pgtype.Int8, error) {
	var id int64
	err := r.db.QueryRow(
		db.Ctx,
		"INSERT INTO users (username, first_name, last_name, email, password_hash, role, language_preference) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		userDto.Username, userDto.FirstName, userDto.LastName, userDto.Email, userDto.PasswordHash, userDto.Role, userDto.LanguagePreference).Scan(&id)

	if err != nil {
		return pgtype.Int8{}, err
	}

	return pgtype.Int8{Int64: id, Valid: true}, nil
}

func (r *Repository) Delete(id pgtype.Int8) error {
	tag, err := r.db.Exec(db.Ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package users/repo: Id invalide: %d", id.Int64)
	}
	return nil
}

func (r *Repository) Update(id pgtype.Int8, user User) error {
	tag, err := r.db.Exec(db.Ctx,
		"UPDATE users SET username=$1, first_name=$2, last_name=$3, email=$4, role=$5, language_preference=$6 WHERE id=$7",
		user.Username, user.FirstName, user.LastName, user.Email, user.Role, user.LanguagePreference, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("package users/repo Update: Id invalide: %d", id.Int64)
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

func (r *Repository) GetScore(userId pgtype.Int8) (int32, error) {
	var totalScore int32
	err := r.db.QueryRow(db.Ctx, "SELECT COALESCE(SUM(points), 0) FROM score_history WHERE user_id = $1", userId).Scan(&totalScore)
	if err != nil {
		return 0, fmt.Errorf("package users/repo GetScore query: %w", err)
	}
	return totalScore, nil
}
