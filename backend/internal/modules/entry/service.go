package entry

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllDepots() ([]DepotFrontend, error) {
	return s.repo.GetAllDepots()
}

func (s *Service) GetAll() ([]Entry, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Entry, error) {
	return s.repo.GetById(id)
}

func (s *Service) Create(e Entry) (pgtype.Int8, error) {
	return s.repo.Create(e)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}

func (s *Service) ValidateDepot(id int64) error {
	return s.repo.ValidateDepot(id)
}
