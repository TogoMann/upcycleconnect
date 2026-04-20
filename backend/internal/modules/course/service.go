package course

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllCatalogue() ([]OffreFrontend, error) {
	return s.repo.GetAllCatalogue()
}

func (s *Service) GetAllApprovedCatalogue() ([]OffreFrontend, error) {
	return s.repo.GetAllApprovedCatalogue()
}

func (s *Service) GetAll() ([]Course, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Course, error) {
	return s.repo.GetById(id)
}

func (s *Service) GetUserCourses(userId pgtype.Int8) ([]UserCourse, error) {
	return s.repo.GetUserCourses(userId)
}

func (s *Service) Create(c Course) (pgtype.Int8, error) {
	return s.repo.Create(c)
}

func (s *Service) Update(id pgtype.Int8, c Course) error {
	return s.repo.Update(id, c)
}

func (s *Service) Approve(id pgtype.Int8, approvedBy pgtype.Int8) error {
	return s.repo.Approve(id, approvedBy)
}

func (s *Service) Disapprove(id pgtype.Int8) error {
	return s.repo.Disapprove(id)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}
