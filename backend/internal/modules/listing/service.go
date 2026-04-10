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
		return nil, fmt.Errorf("listing/service Listing ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(loDto Listing) (pgtype.Int8, error) {
	val, err := loDto.Price.Value()

	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("listing/service Listing prix invalide: %v", err.Error())
	}

	if loDto.Price.Int.Sign() < 0 {
		return pgtype.Int8{}, fmt.Errorf("listing/service Listing prix négatif: %v", val)
	}
	return s.repo.Create(loDto)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("listing/service Listing ID invalide: %d", id)
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
