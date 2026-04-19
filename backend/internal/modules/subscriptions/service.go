package subscriptions

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllAbonnements() ([]AbonnementFrontend, error) {
	return s.repo.GetAllAbonnements()
}

func (s *Service) GetAll() ([]Subscription, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Subscription, error) {
	return s.repo.GetById(id)
}

func (s *Service) Create(sub Subscription) (pgtype.Int8, error) {
	return s.repo.Create(sub)
}

func (s *Service) Update(id pgtype.Int8, sub Subscription) error {
	return s.repo.Update(id, sub)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}
