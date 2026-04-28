package planning

import "github.com/jackc/pgx/v5/pgtype"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserPlanning(userId pgtype.Int8) ([]PlanningItem, error) {
	return s.repo.GetUserPlanning(userId)
}

func (s *Service) GetAllPlannings() ([]AdminPlanningItem, error) {
	return s.repo.GetAllPlannings()
}

func (s *Service) CreatePersonalEvent(e PersonalEvent) (pgtype.Int8, error) {
	return s.repo.CreatePersonalEvent(e)
}

func (s *Service) DeletePersonalEvent(id pgtype.Int8, userId pgtype.Int8) error {
	return s.repo.DeletePersonalEvent(id, userId)
}
