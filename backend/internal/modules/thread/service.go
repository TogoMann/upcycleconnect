package thread

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

func (s *Service) GetAll() ([]Thread, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id int64) (*Thread, error) {
	if id < 1 {
		return nil, fmt.Errorf("thread/service Thread ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(threadDto Thread) (int64, error) {
	if strings.TrimSpace(threadDto.Title) == "" || strings.TrimSpace(threadDto.Content) == "" {
		return 0, fmt.Errorf("thread/service Invalid string(s): Missing values.")
	}

	return s.repo.Create(threadDto)
}

func (s *Service) Delete(id int64) error {
	if id < 1 {
		return fmt.Errorf("thread/service Thread ID invalide: %d", id)
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
