package users

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]User, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*User, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("users/service User ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(userDto User) (pgtype.Int8, error) {
	if strings.TrimSpace(userDto.Username) == "" || strings.TrimSpace(userDto.FirstName) == "" || strings.TrimSpace(userDto.LastName) == "" || strings.TrimSpace(userDto.Email) == "" {
		return pgtype.Int8{}, fmt.Errorf("users/service Invalid string(s): Missing values.")
	}
	if !IsValidRole(userDto.Role) {
		return pgtype.Int8{}, fmt.Errorf("invalid role: %s", userDto.Role)
	}

	return s.repo.Create(userDto)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("users/service User ID invalide: %d", id)
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
		return 0, fmt.Errorf("users/service User ID invalide: %d", userId)
	}
	return s.repo.GetScore(userId)
}
