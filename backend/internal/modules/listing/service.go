package listing

import (
	"backend/internal/modules/container"
	"backend/internal/modules/item"
	"backend/internal/modules/subscriptions"
	"backend/internal/modules/users"
	"backend/internal/utils"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo             *Repository
	subRepo          *subscriptions.Repository
	containerService *container.Service
	itemService      *item.Service
	userService      *users.Service
}

func NewService(repo *Repository, subRepo *subscriptions.Repository, containerService *container.Service, itemService *item.Service, userService *users.Service) *Service {
	return &Service{repo: repo, subRepo: subRepo, containerService: containerService, itemService: itemService, userService: userService}
}

func (s *Service) GetAll(page, limit int) (*PaginatedListings, error) {
	return s.repo.GetAll(page, limit)
}

func (s *Service) GetAllApproved(excludeUserId pgtype.Int8, minLevel int32) ([]Listing, error) {
	minScore := int32(0)
	if minLevel > 1 {
		minScore = (minLevel - 1) * 100
	}
	return s.repo.GetAllApproved(excludeUserId, minScore)
}

func (s *Service) GetByUserId(userId pgtype.Int8) ([]Listing, error) {
	if !userId.Valid || userId.Int64 < 1 {
		return nil, fmt.Errorf("listing/service User ID invalide: %d", userId.Int64)
	}
	return s.repo.GetByUserId(userId)
}

func (s *Service) GetById(id pgtype.Int8) (*Listing, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("listing/service Listing ID invalide: %d", id.Int64)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(loDto Listing, lockerId pgtype.Int8, physicalState string, itemSize string) (pgtype.Int8, error) {
	if loDto.Category == "" {
		return pgtype.Int8{}, fmt.Errorf("listing/service Listing category manquante")
	}

	count, err := s.repo.CountByUserId(loDto.CreatedBy)
	if err != nil {
		return pgtype.Int8{}, err
	}

	u, err := s.userService.GetById(loDto.CreatedBy)
	if err != nil {
		return pgtype.Int8{}, err
	}

	tier, err := s.subRepo.GetActiveTierByUserId(loDto.CreatedBy)
	if err != nil {
		return pgtype.Int8{}, err
	}

	limit := 5
	if tier == "Premium" {
		limit = 20
	} else if tier == "Pro" || (u != nil && u.Role == users.Pro) {
		limit = 1000000
	}

	if count >= int64(limit) {
		return pgtype.Int8{}, fmt.Errorf("limite d'annonces atteinte pour votre plan (%s: %d max)", tier, limit)
	}

	val, err := loDto.Price.Value()

	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("listing/service Listing prix invalide: %v", err.Error())
	}

	if loDto.Price.Int.Sign() < 0 {
		return pgtype.Int8{}, fmt.Errorf("listing/service Listing prix négatif: %v", val)
	}

	if loDto.HandoffMode == Locker {
		itemId, err := s.depositIntoLocker(loDto, lockerId, physicalState, itemSize)
		if err != nil {
			return pgtype.Int8{}, err
		}
		loDto.ItemId = itemId
	}

	id, err := s.repo.Create(loDto)
	if err != nil {
		if loDto.ItemId.Valid {
			_ = s.containerService.UpdateLockerStatus(lockerId, "Available")
			_ = s.itemService.Delete(loDto.ItemId)
		}
		return id, err
	}

	s.userService.AddScore(loDto.CreatedBy, utils.ActionListingCreated.Points, utils.ActionListingCreated.Description)

	return id, nil
}

func (s *Service) depositIntoLocker(loDto Listing, lockerId pgtype.Int8, physicalState string, itemSize string) (pgtype.Int8, error) {
	if !lockerId.Valid {
		return pgtype.Int8{}, fmt.Errorf("listing/service un casier doit être sélectionné pour ce mode de remise")
	}

	size := item.ItemSize(itemSize)
	if _, ok := item.SizeRank[size]; !ok {
		size = item.SizeM
	}

	locker, err := s.containerService.GetLockerById(lockerId)
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("listing/service casier introuvable")
	}
	if locker.Status != "Available" {
		return pgtype.Int8{}, fmt.Errorf("listing/service ce casier n'est plus disponible")
	}
	if item.SizeRank[item.ItemSize(locker.Size)] < item.SizeRank[size] {
		return pgtype.Int8{}, fmt.Errorf("listing/service ce casier est trop petit pour cet objet")
	}

	containerObj, err := s.containerService.GetById(locker.ContainerId)
	if err != nil {
		return pgtype.Int8{}, err
	}

	state := item.NormalizeState(physicalState)

	claimed, err := s.containerService.ClaimLocker(lockerId)
	if err != nil {
		return pgtype.Int8{}, err
	}
	if !claimed {
		return pgtype.Int8{}, fmt.Errorf("listing/service ce casier n'est plus disponible")
	}

	itemId, err := s.itemService.Create(item.Item{
		OwnerId:       loDto.CreatedBy,
		LockerId:      lockerId,
		SiteId:        containerObj.SiteId,
		MaterialType:  string(loDto.Category),
		PhysicalState: state,
		Size:          size,
		Status:        item.Deposited,
		Weight:        loDto.Weight,
	})
	if err != nil {
		_ = s.containerService.UpdateLockerStatus(lockerId, "Available")
		return pgtype.Int8{}, err
	}

	return itemId, nil
}

func (s *Service) Update(id pgtype.Int8, l Listing) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("listing/service ID invalide: %d", id.Int64)
	}

	if l.Category == "" {
		return fmt.Errorf("listing/service Listing category manquante")
	}

	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("listing/service listing not found")
	}

	return s.repo.Update(id, l)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("listing/service ID invalide: %d", id.Int64)
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

func (s *Service) UpdateStatus(id pgtype.Int8, status ListingStatus) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("listing/service ID invalide")
	}
	return s.repo.UpdateStatus(id, status)
}

func (s *Service) Approve(id pgtype.Int8, adminId pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("listing/service Listing ID invalide")
	}
	if !adminId.Valid || adminId.Int64 < 1 {
		return fmt.Errorf("listing/service Admin ID invalide")
	}
	return s.repo.Approve(id, adminId)
}

func (s *Service) Disapprove(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("listing/service Listing ID invalide")
	}
	return s.repo.Disapprove(id)
}
