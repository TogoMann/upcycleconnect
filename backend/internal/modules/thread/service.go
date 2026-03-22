package thread

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

func (s *Service) GetAll() ([]Thread, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Thread, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("thread/service Thread ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(threadDto Thread) (pgtype.Int8, error) {
	if strings.TrimSpace(threadDto.Title) == "" || strings.TrimSpace(threadDto.Content) == "" {
		return pgtype.Int8{}, fmt.Errorf("thread/service Invalid string(s): Missing values.")
	}

	return s.repo.Create(threadDto)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
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
