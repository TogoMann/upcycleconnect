package brands

import "github.com/jackc/pgx/v5/pgtype"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]Brand, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Brand, error) {
	return s.repo.GetById(id)
}

func (s *Service) GetByCreator(userId pgtype.Int8) ([]Brand, error) {
	return s.repo.GetByCreator(userId)
}

func (s *Service) Create(b Brand) (pgtype.Int8, error) {
	return s.repo.Create(b)
}

func (s *Service) Update(id pgtype.Int8, b Brand) error {
	return s.repo.Update(id, b)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}
