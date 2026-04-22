package users

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]UserFrontend, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*User, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("users/service User ID invalide: %d", id.Int64)
	}

	return s.repo.GetById(id)
}

func (s *Service) GetByUsername(username string) (*User, error) {
	if strings.TrimSpace(username) == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	return s.repo.GetByUsername(username)
}

func (s *Service) Create(userDto User) (pgtype.Int8, error) {
	if strings.TrimSpace(userDto.Username) == "" || strings.TrimSpace(userDto.FirstName) == "" || strings.TrimSpace(userDto.LastName) == "" || strings.TrimSpace(userDto.Email) == "" {
		return pgtype.Int8{}, fmt.Errorf("users/service Invalid string(s): Missing values.")
	}
	if !IsValidRole(userDto.Role) {
		return pgtype.Int8{}, fmt.Errorf("invalid role: %s", userDto.Role)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("failed to hash password: %w", err)
	}
	userDto.PasswordHash = string(hashedPassword)

	return s.repo.Create(userDto)
}

func (s *Service) Update(id pgtype.Int8, user User) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("users/service ID invalide: %d", id.Int64)
	}

	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("users/service user not found")
	}

	if !IsValidRole(user.Role) {
		return fmt.Errorf("invalid role: %s", user.Role)
	}

	return s.repo.Update(id, user)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("users/service User ID invalide: %d", id.Int64)
	}

	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	return s.repo.Delete(id)
}

func (s *Service) GetScore(userId pgtype.Int8) (int32, error) {
	if !userId.Valid || userId.Int64 < 1 {
		return 0, fmt.Errorf("users/service User ID invalide: %d", userId.Int64)
	}
	return s.repo.GetScore(userId)
}

func (s *Service) GetScoreHistory(userId pgtype.Int8) ([]ScoreHistory, error) {
	if !userId.Valid || userId.Int64 < 1 {
		return nil, fmt.Errorf("users/service User ID invalide: %d", userId.Int64)
	}
	return s.repo.GetScoreHistory(userId)
}

func (s *Service) AddScore(userId pgtype.Int8, points int32, description string) error {
	if !userId.Valid || userId.Int64 < 1 {
		return fmt.Errorf("users/service User ID invalide: %d", userId.Int64)
	}
	return s.repo.AddScore(userId, points, description)
}

func (s *Service) UpdateTutorialSeen(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("users/service User ID invalide")
	}
	return s.repo.UpdateTutorialSeen(id)
}
