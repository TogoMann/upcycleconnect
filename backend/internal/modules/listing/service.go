package listing

import (
	"fmt"
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

func (s *Service) GetById(id int64) (*Listing, error) {
	if id < 1 {
		return nil, fmt.Errorf("listing/service Listing order ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(loDto Listing) (int64, error) {
	val, err := loDto.Price.Value()

	if err != nil {
		return 0, fmt.Errorf("listing/service Listing prix invalide: %v", err.Error())
	}

	if loDto.Price.Int.Sign() < 0 {
		return 0, fmt.Errorf("listing/service Listing prix négatif: %v", val)
	}
	return s.repo.Create(loDto)
}

func (s *Service) Delete(id int64) error {
	if id < 1 {
		return fmt.Errorf("listing/service Thread ID invalide: %d", id)
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
