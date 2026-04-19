package project

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllProjets() ([]ProjetFrontend, error) {
	return s.repo.GetAllProjets()
}

func (s *Service) GetAll() ([]Project, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Project, error) {
	return s.repo.GetById(id)
}

func (s *Service) Create(p Project) (pgtype.Int8, error) {
	return s.repo.Create(p)
}

func (s *Service) Update(id pgtype.Int8, p Project) error {
	return s.repo.Update(id, p)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}

func (s *Service) UpdateFeatured(id int64, featured bool) error {
	return s.repo.UpdateFeatured(id, featured)
}

func (s *Service) GetSteps(projectId pgtype.Int8) ([]Step, error) {
	return s.repo.GetSteps(projectId)
}

func (s *Service) CreateStep(st Step) (pgtype.Int8, error) {
	return s.repo.CreateStep(st)
}

func (s *Service) UpdateStep(id int64, st Step) error {
	return s.repo.UpdateStep(id, st)
}

func (s *Service) DeleteStep(id int64) error {
	return s.repo.DeleteStep(id)
}
