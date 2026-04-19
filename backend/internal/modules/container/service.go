package container

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]ConteneurFrontend, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Container, error) {
	return s.repo.GetById(id)
}

func (s *Service) Create(c Container) (pgtype.Int8, error) {
	return s.repo.Create(c)
}

func (s *Service) Update(id pgtype.Int8, c Container) error {
	return s.repo.Update(id, c)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}
