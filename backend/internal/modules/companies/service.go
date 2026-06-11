package companies

import (
	"backend/internal/utils"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCompany(c Company) (pgtype.Int8, error) {
	if !utils.VerifySiret(c.Siret) {
		return pgtype.Int8{}, fmt.Errorf("invalid SIRET number")
	}

	existing, _ := s.repo.GetBySiret(c.Siret)
	if existing != nil {
		return pgtype.Int8{}, fmt.Errorf("company already exists")
	}

	return s.repo.Create(c)
}

func (s *Service) GetCompanyById(id int64) (*Company, error) {
	return s.repo.GetById(pgtype.Int8{Int64: id, Valid: true})
}

func (s *Service) GetCompanyBySiret(siret string) (*Company, error) {
	return s.repo.GetBySiret(siret)
}

func (s *Service) GetAllCompanies() ([]Company, error) {
	return s.repo.GetAll()
}
