package comments

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]Comments, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Comments, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("comments/service ID invalide: %d", id)
	}
	return s.repo.GetById(id)
}

func (s *Service) Create(dto Comments) (pgtype.Int8, error) {
	return s.repo.Create(dto)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("comments/service ID invalide: %d", id)
	}

	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	hasReplies, err := s.repo.HasReplies(id)
	if err != nil {
		return err
	}

	if hasReplies {
		return s.repo.SoftDelete(id)
	}

	return s.repo.Delete(id)
}
