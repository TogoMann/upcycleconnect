package event

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

func (s *Service) GetAll() ([]Event, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Event, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("event/service ID invalide: %d", id.Int64)
	}
	return s.repo.GetById(id)
}

func (s *Service) Create(dto Event) (pgtype.Int8, error) {
	return s.repo.Create(dto)
}

func (s *Service) Update(id pgtype.Int8, e Event) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("event/service ID invalide: %d", id.Int64)
	}

	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("event/service event not found")
	}

	return s.repo.Update(id, e)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("event/service ID invalide: %d", id.Int64)
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

func (s *Service) Approve(id pgtype.Int8, adminId pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("event/service Event ID invalide")
	}
	if !adminId.Valid || adminId.Int64 < 1 {
		return fmt.Errorf("event/service Admin ID invalide")
	}
	return s.repo.Approve(id, adminId)
}
