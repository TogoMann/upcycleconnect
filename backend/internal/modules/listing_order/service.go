package listingorder

import (
	"github.com/jackc/pgx/v5/pgtype"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]ListingOrder, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*ListingOrder, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("listing_order/service Listing order ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(loDto ListingOrder) (pgtype.Int8, error) {
	return s.repo.Create(loDto)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("listing_order/service Thread ID invalide: %d", id)
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
