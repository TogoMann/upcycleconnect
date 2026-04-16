package listing

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

func (s *Service) GetAll() ([]Listing, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Listing, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("listing/service Listing ID invalide: %d", id.Int64)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(loDto Listing) (pgtype.Int8, error) {
	if loDto.Category == "" {
		return pgtype.Int8{}, fmt.Errorf("listing/service Listing category manquante")
	}

	val, err := loDto.Price.Value()

	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("listing/service Listing prix invalide: %v", err.Error())
	}

	if loDto.Price.Int.Sign() < 0 {
		return pgtype.Int8{}, fmt.Errorf("listing/service Listing prix négatif: %v", val)
	}
	return s.repo.Create(loDto)
}

func (s *Service) Update(id pgtype.Int8, l Listing) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("listing/service ID invalide: %d", id.Int64)
	}

	if l.Category == "" {
		return fmt.Errorf("listing/service Listing category manquante")
	}

	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("listing/service listing not found")
	}

	return s.repo.Update(id, l)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("listing/service ID invalide: %d", id.Int64)
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
		return fmt.Errorf("listing/service Listing ID invalide")
	}
	if !adminId.Valid || adminId.Int64 < 1 {
		return fmt.Errorf("listing/service Admin ID invalide")
	}
	return s.repo.Approve(id, adminId)
}

