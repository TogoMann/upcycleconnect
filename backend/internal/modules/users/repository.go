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

func (r *Repository) GetAll(page, limit int) (*PaginatedUsers, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int
	err := r.db.QueryRow(db.Ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("package users/repo GetAll count query: %w", err)
	}

	rows, err := r.db.Query(db.Ctx, `
		SELECT 
			u.id, u.username, u.first_name, u.last_name, u.email, u.role, u.language_preference, u.has_seen_tutorial,
			TO_CHAR(u.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
			CAST((SELECT COALESCE(SUM(points), 0) FROM score_history WHERE user_id = u.id) AS INTEGER) as score,
			COALESCE(c.siret, '') as siret,
			COALESCE(u.company_id, 0) as company_id,
			COALESCE((SELECT tier FROM subscriptions WHERE subscriber_id = u.id AND until >= CURRENT_DATE ORDER BY until DESC LIMIT 1), 'Free') as plan,
			COALESCE(up.predicted_service_type, '') as predicted_service,
			COALESCE(up.probability, 0) as probability
		FROM users u
		LEFT JOIN companies c ON u.company_id = c.id
		LEFT JOIN user_predictions up ON u.id = up.user_id
		ORDER BY u.id DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("package users/repo GetAll users query: %w", err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[UserFrontend])
	if err != nil {
		return nil, fmt.Errorf("package users/repo CollectRows: %w", err)
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}

	return &PaginatedUsers{
		Data:       users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
func (r *Repository) GetById(id pgtype.Int8) (*User, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, username, first_name, last_name, email, password_hash, role, language_preference, has_seen_tutorial, created_at, company_id FROM users WHERE id = $1", id)
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
	rows, err := r.db.Query(db.Ctx, "SELECT id, username, first_name, last_name, email, password_hash, role, language_preference, has_seen_tutorial, created_at, company_id FROM users WHERE username = $1", username)
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

func (r *Repository) GetByEmail(email string) (*User, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, username, first_name, last_name, email, password_hash, role, language_preference, has_seen_tutorial, created_at, company_id FROM users WHERE email = $1", email)
	if err != nil {
		return nil, fmt.Errorf("package users/repo GetByEmail query: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("package users/repo GetByEmail: %v", err.Error())
	}
	return &user, nil
}

func (r *Repository) GetMe(username string) (*UserFrontend, error) {
	var u UserFrontend
	err := r.db.QueryRow(db.Ctx, `
		SELECT 
			u.id, u.username, u.first_name, u.last_name, u.email, u.role, u.language_preference, u.has_seen_tutorial, 
			TO_CHAR(u.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
			CAST((SELECT COALESCE(SUM(points), 0) FROM score_history WHERE user_id = u.id) AS INTEGER) as score,
			COALESCE(c.siret, '') as siret,
			COALESCE(u.company_id, 0) as company_id,
			COALESCE((SELECT tier FROM subscriptions WHERE subscriber_id = u.id AND until >= CURRENT_DATE ORDER BY until DESC LIMIT 1), 'Free') as plan,
			COALESCE(up.predicted_service_type, '') as predicted_service,
			COALESCE(up.probability, 0) as probability
		FROM users u
		LEFT JOIN companies c ON u.company_id = c.id
		LEFT JOIN user_predictions up ON u.id = up.user_id
		WHERE u.username = $1
	`, username).Scan(&u.Id, &u.Username, &u.FirstName, &u.LastName, &u.Email, &u.Role, &u.LanguagePreference, &u.HasSeenTutorial, &u.CreatedAt, &u.Score, &u.Siret, &u.CompanyId, &u.Plan, &u.PredictedService, &u.Probability)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("package users/repo GetMe: %w", err)
	}
	return &u, nil
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
		"INSERT INTO users (username, first_name, last_name, email, password_hash, role, language_preference, company_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		userDto.Username, userDto.FirstName, userDto.LastName, userDto.Email, userDto.PasswordHash, userDto.Role, userDto.LanguagePreference, userDto.CompanyId).Scan(&id)

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
		"UPDATE users SET username=$1, first_name=$2, last_name=$3, email=$4, role=$5, language_preference=$6, company_id=$7 WHERE id=$8",
		user.Username, user.FirstName, user.LastName, user.Email, user.Role, user.LanguagePreference, user.CompanyId, id)
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

func (r *Repository) GetScoreHistory(userId pgtype.Int8) ([]ScoreHistory, error) {
	rows, err := r.db.Query(db.Ctx, "SELECT id, user_id, points, description, TO_CHAR(created_at, 'YYYY-MM-DD\"T\"HH24:MI:SS\"Z\"') as created_at FROM score_history WHERE user_id = $1 ORDER BY created_at DESC", userId)
	if err != nil {
		return nil, fmt.Errorf("package users/repo GetScoreHistory query: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[ScoreHistory])
}

func (r *Repository) AddScore(userId pgtype.Int8, points int32, description string) error {
	_, err := r.db.Exec(db.Ctx, "INSERT INTO score_history (user_id, points, description) VALUES ($1, $2, $3)", userId, points, description)
	if err != nil {
		return fmt.Errorf("package users/repo AddScore: %w", err)
	}
	return nil
}

func (r *Repository) CountScoreEventsSince(userId pgtype.Int8, description string, windowDays int) (int64, error) {
	var count int64
	err := r.db.QueryRow(db.Ctx,
		"SELECT COUNT(*) FROM score_history WHERE user_id = $1 AND description = $2 AND created_at >= NOW() - make_interval(days => $3)",
		userId, description, windowDays).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("package users/repo CountScoreEventsSince: %w", err)
	}
	return count, nil
}

func (r *Repository) CreateResetToken(userId int64, token string, expiresAt pgtype.Timestamp) error {
	_, err := r.db.Exec(db.Ctx, "INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", userId, token, expiresAt)
	if err != nil {
		return fmt.Errorf("package users/repo CreateResetToken: %w", err)
	}
	return nil
}

func (r *Repository) GetUserIdByResetToken(token string) (int64, error) {
	var userId int64
	err := r.db.QueryRow(db.Ctx, "SELECT user_id FROM password_reset_tokens WHERE token = $1 AND expires_at > NOW()", token).Scan(&userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("invalid or expired token")
		}
		return 0, fmt.Errorf("package users/repo GetUserIdByResetToken: %w", err)
	}
	return userId, nil
}

func (r *Repository) DeleteResetToken(token string) error {
	_, err := r.db.Exec(db.Ctx, "DELETE FROM password_reset_tokens WHERE token = $1", token)
	if err != nil {
		return fmt.Errorf("package users/repo DeleteResetToken: %w", err)
	}
	return nil
}

func (r *Repository) UpdatePassword(userId int64, passwordHash string) error {
	_, err := r.db.Exec(db.Ctx, "UPDATE users SET password_hash = $1 WHERE id = $2", passwordHash, userId)
	if err != nil {
		return fmt.Errorf("package users/repo UpdatePassword: %w", err)
	}
	return nil
}
