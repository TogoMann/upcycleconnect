package plans

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]Plan, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Plan, error) {
	return s.repo.GetById(id)
}

func (s *Service) Create(p Plan) (pgtype.Int8, error) {
	return s.repo.Create(p)
}

func (s *Service) Update(id pgtype.Int8, p Plan) error {
	return s.repo.Update(id, p)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}
