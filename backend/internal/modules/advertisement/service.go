package advertisement

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllPubs() ([]PubFrontend, error) {
	return s.repo.GetAllPubs()
}

func (s *Service) GetAll() ([]Advertisement, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Advertisement, error) {
	return s.repo.GetById(id)
}

func (s *Service) Create(ad Advertisement) (pgtype.Int8, error) {
	return s.repo.Create(ad)
}

func (s *Service) UpdateStatus(id pgtype.Int8, status string, approvedBy pgtype.Int8) error {
	return s.repo.UpdateStatus(id, status, approvedBy)
}

func (s *Service) Reject(id pgtype.Int8) error {
	return s.repo.Reject(id)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}
