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

func (s *Service) GetByUserId(userId pgtype.Int8) ([]Item, error) {
	return s.repo.GetByUserId(userId)
}

func (s *Service) Create(item Item) (pgtype.Int8, error) {
	id, err := s.repo.Create(item)
	if err != nil {
		return id, err
	}
	if item.LockerId.Valid {
		_ = s.repo.UpdateLockerStatus(item.LockerId, "Occupied")
	}
	return id, nil
}

func (s *Service) Update(id pgtype.Int8, item Item) error {
	return s.repo.Update(id, item)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}

func (s *Service) RequestDeposit(dto Item) error {
	
	
	
	return nil
}

func (s *Service) ValidateAndGenerateCode(itemId pgtype.Int8, userId pgtype.Int8) (string, error) {
	item, err := s.repo.GetById(itemId)
	if err != nil {
		return "", err
	}

	err = s.repo.UpdateStatus(itemId, Validated)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	err = s.repo.CreateAccessCode(item.LockerId, itemId, userId, code)
	return code, err
}

func (s *Service) Collect(itemId pgtype.Int8, proId pgtype.Int8, code string) error {
	err := s.repo.MarkCodeUsed(code)
	if err != nil {
		return err
	}

	item, err := s.repo.GetById(itemId)
	if err == nil && item.LockerId.Valid {
		_ = s.repo.UpdateLockerStatus(item.LockerId, "Available")
	}

	return s.repo.Collect(itemId, proId)
}
