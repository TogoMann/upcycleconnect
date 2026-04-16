package item

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]Item, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Item, error) {
	return s.repo.GetById(id)
}

func (s *Service) Create(item Item) (pgtype.Int8, error) {
	return s.repo.Create(item)
}

func (s *Service) Update(id pgtype.Int8, item Item) error {
	return s.repo.Update(id, item)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}

func (s *Service) RequestDeposit(dto Item) error {
	// Logic to save item with status 'deposited'
	// TODO: Create method
	// (Implementation depends on repository Create)
	return nil
}

func (s *Service) ValidateAndGenerateCode(itemId pgtype.Int8, userId pgtype.Int8) (string, error) {
	err := s.repo.UpdateStatus(itemId, Validated)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	err = s.repo.CreateAccessCode(itemId, userId, code)
	return code, err
}

func (s *Service) Collect(itemId pgtype.Int8, proId pgtype.Int8, code string) error {
	err := s.repo.MarkCodeUsed(code)
	if err != nil {
		return err
	}

	// L'objet appartient maintenant au pro
	return s.repo.Collect(itemId, proId)
}
