package news

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

func (s *Service) GetAll() ([]News, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*News, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("news/service News ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(newsDto News) (pgtype.Int8, error) {
	if strings.TrimSpace(newsDto.Content) == "" || strings.TrimSpace(newsDto.Title) == "" {
		return pgtype.Int8{}, fmt.Errorf("news/service Invalid string(s): Missing values.")
	}

	return s.repo.Create(newsDto)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("news/service News ID invalide: %d", id)
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
