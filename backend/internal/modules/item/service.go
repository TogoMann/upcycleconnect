package item

import (
	"backend/internal/modules/users"
	"backend/internal/utils"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo        *Repository
	userService *users.Service
}

func NewService(repo *Repository, userService *users.Service) *Service {
	return &Service{repo: repo, userService: userService}
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
	item.PhysicalState = NormalizeState(string(item.PhysicalState))

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
	item.PhysicalState = NormalizeState(string(item.PhysicalState))
	return s.repo.Update(id, item)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}

func (s *Service) GetAdminDepots() ([]AdminDepot, error) {
	return s.repo.GetAdminDepots()
}

func (s *Service) AdminValidateDepot(id pgtype.Int8) error {
	itm, err := s.repo.GetById(id)
	if err != nil {
		return err
	}
	alreadyValidated := itm.Status == Validated || itm.Status == Collected

	if err := s.repo.UpdateStatus(id, Validated); err != nil {
		return err
	}

	if !alreadyValidated {
		s.userService.AddScore(itm.OwnerId, utils.ActionDepositValidated.Points, utils.ActionDepositValidated.Description)
	}

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

	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", n.Int64())

	err = s.repo.CreateAccessCode(item.LockerId, itemId, userId, code)
	return code, err
}

func (s *Service) Collect(itemId pgtype.Int8, proId pgtype.Int8, code string) error {
	if err := s.repo.MarkCodeUsed(code); err != nil {
		return err
	}

	itm, fetchErr := s.repo.GetById(itemId)
	if fetchErr == nil && itm.LockerId.Valid {
		_ = s.repo.UpdateLockerStatus(itm.LockerId, "Available")
	}

	if err := s.repo.Collect(itemId, proId); err != nil {
		return err
	}

	if fetchErr == nil && itm.OwnerId.Valid {
		s.userService.AddScore(itm.OwnerId, utils.ActionMaterialCollection.Points, utils.ActionMaterialCollection.Description)
	}

	return nil
}
