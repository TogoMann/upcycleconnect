package course

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

func (s *Service) GetAll() ([]Course, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Course, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("course/service ID invalide: %d", id.Int64)
	}
	return s.repo.GetById(id)
}

func (s *Service) GetUserCourses(id pgtype.Int8) ([]UserCourse, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("course/service User ID invalide: %d", id.Int64)
	}
	return s.repo.GetUserCourses(id)
}

func (s *Service) Create(dto Course) (pgtype.Int8, error) {
	return s.repo.Create(dto)
}

func (s *Service) Update(id pgtype.Int8, c Course) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("course/service ID invalide: %d", id.Int64)
	}

	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("course/service course not found")
	}

	return s.repo.Update(id, c)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("course/service ID invalide: %d", id.Int64)
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

func (s *Service) Approve(id pgtype.Int8, adminId pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("course/service Course ID invalide")
	}
	if !adminId.Valid || adminId.Int64 < 1 {
		return fmt.Errorf("course/service Admin ID invalide")
	}
	return s.repo.Approve(id, adminId)
}
