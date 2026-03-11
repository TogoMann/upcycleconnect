package users

import (
	"fmt"
	"strings"
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

func (s *Service) GetById(id int64) (*User, error) {
	if id < 1 {
		return nil, fmt.Errorf("users/service User ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(userDto User) (int64, error) {
	if strings.TrimSpace(userDto.FirstName) == "" || strings.TrimSpace(userDto.LastName) == "" || strings.TrimSpace(userDto.Email) == "" {
		return 0, fmt.Errorf("users/service Invalid string(s): Missing values.")
	}
	if !IsValidRole(userDto.Role) {
		return 0, fmt.Errorf("invalid role: %s", userDto.Role)
	}

	return s.repo.Create(userDto)
}

func (s *Service) Delete(id int64) error {
	if id < 1 {
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
