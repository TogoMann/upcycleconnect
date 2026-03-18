package news

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

func (s *Service) GetAll() ([]News, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id int64) (*News, error) {
	if id < 1 {
		return nil, fmt.Errorf("news/service Thread ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(newsDto News) (int64, error) {
	if strings.TrimSpace(newsDto.Content) == "" || strings.TrimSpace(newsDto.Title) == "" {
		return 0, fmt.Errorf("news/service Invalid string(s): Missing values.")
	}

	return s.repo.Create(newsDto)
}

func (s *Service) Delete(id int64) error {
	if id < 1 {
		return fmt.Errorf("news/service Thread ID invalide: %d", id)
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
