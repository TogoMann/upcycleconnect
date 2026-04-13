package subscriptions

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

func (s *Service) GetAll() ([]Subscriptions, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Subscriptions, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("subscriptions/service ID invalide: %d", id)
	}
	return s.repo.GetById(id)
}

func (s *Service) Create(dto Subscriptions) (pgtype.Int8, error) {
	return s.repo.Create(dto)
}

func (s *Service) Update(id pgtype.Int8, sub Subscriptions) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("subscriptions/service ID invalide: %d", id)
	}

	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("subscriptions/service subscription not found")
	}

	return s.repo.Update(id, sub)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("subscriptions/service ID invalide: %d", id)
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
