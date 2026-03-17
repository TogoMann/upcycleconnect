package post

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

func (s *Service) GetAll() ([]Post, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id int64) (*Post, error) {
	if id < 1 {
		return nil, fmt.Errorf("post/service Thread ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(postDto Post) (int64, error) {
	if strings.TrimSpace(postDto.Content) == "" {
		return 0, fmt.Errorf("post/service Invalid string(s): Missing values.")
	}

	return s.repo.Create(postDto)
}

func (s *Service) Delete(id int64) error {
	if id < 1 {
		return fmt.Errorf("post/service Thread ID invalide: %d", id)
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
