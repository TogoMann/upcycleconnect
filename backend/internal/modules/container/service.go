package container

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) generateAccessCode() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *Service) CreateLockerAccess(lockerId, itemId, userId pgtype.Int8) (LockerAccess, error) {
	code := s.generateAccessCode()
	expiresAt := time.Now().Add(72 * time.Hour)

	access := LockerAccess{
		LockerId:   lockerId,
		ItemId:     itemId,
		UserId:     userId,
		AccessCode: code,
		ExpiresAt:  pgtype.Timestamp{Time: expiresAt, Valid: true},
	}

	id, err := s.repo.CreateLockerAccess(access)
	if err != nil {
		return LockerAccess{}, err
	}
	access.Id = id
	return access, nil
}

func (s *Service) GetUserAccesses(userId pgtype.Int8) ([]LockerAccessDetails, error) {
	return s.repo.GetAccessesByUserId(userId)
}

func (s *Service) GetAll() ([]ConteneurFrontend, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Container, error) {
	return s.repo.GetById(id)
}

func (s *Service) GetLockersByContainerId(containerId pgtype.Int8) ([]Locker, error) {
	return s.repo.GetLockersByContainerId(containerId)
}

func (s *Service) CreateLocker(l Locker) (pgtype.Int8, error) {
	return s.repo.CreateLocker(l)
}

func (s *Service) UpdateLocker(id pgtype.Int8, l Locker) error {
	return s.repo.UpdateLocker(id, l)
}

func (s *Service) DeleteLocker(id pgtype.Int8) error {
	return s.repo.DeleteLocker(id)
}

func (s *Service) Create(c Container) (pgtype.Int8, error) {
	return s.repo.Create(c)
}

func (s *Service) Update(id pgtype.Int8, c Container) error {
	return s.repo.Update(id, c)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}
