package payment_methods

import "github.com/jackc/pgx/v5/pgtype"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]PaymentMethod, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*PaymentMethod, error) {
	return s.repo.GetById(id)
}

func (s *Service) GetByUserId(userId pgtype.Int8) ([]PaymentMethod, error) {
	return s.repo.GetByUserId(userId)
}

func (s *Service) HasPaymentMethod(userId pgtype.Int8) (bool, error) {
	return s.repo.HasPaymentMethod(userId)
}

func (s *Service) Create(pm PaymentMethod) (pgtype.Int8, error) {
	return s.repo.Create(pm)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}

func (s *Service) DeleteByIdAndUserId(id pgtype.Int8, userId pgtype.Int8) error {
	return s.repo.DeleteByIdAndUserId(id, userId)
}

func (s *Service) SetDefault(id pgtype.Int8, userId pgtype.Int8) error {
	return s.repo.SetDefault(id, userId)
}
