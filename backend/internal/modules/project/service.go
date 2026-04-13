package project

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

func (s *Service) GetAll() ([]Project, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Project, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("project/service ID invalide: %d", id)
	}
	return s.repo.GetById(id)
}

func (s *Service) Create(dto Project) (pgtype.Int8, error) {
	if dto.Status == "" {
		dto.Status = InProgress
	}
	return s.repo.Create(dto)
}

func (s *Service) Update(id pgtype.Int8, p Project) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("project/service ID invalide: %d", id)
	}
	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("project/service project not found")
	}
	return s.repo.Update(id, p)
}

func (s *Service) CreateStep(dto ProjectStep) (pgtype.Int8, error) {
	return s.repo.CreateStep(dto)
}

func (s *Service) UpdateStep(id pgtype.Int8, dto ProjectStep) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("project/service Step ID invalide: %d", id)
	}
	return s.repo.UpdateStep(id, dto)
}

func (s *Service) DeleteStep(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("project/service Step ID invalide: %d", id)
	}
	return s.repo.DeleteStep(id)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("project/service ID invalide: %d", id)
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

func (s *Service) GetSteps(projectId pgtype.Int8) ([]ProjectStep, error) {
	if !projectId.Valid || projectId.Int64 < 1 {
		return nil, fmt.Errorf("project/service Project ID invalide: %d", projectId)
	}
	return s.repo.GetSteps(projectId)
}
