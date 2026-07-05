package listingorder

import (
	"backend/internal/modules/chat"
	"backend/internal/modules/container"
	"backend/internal/modules/financial"
	"backend/internal/modules/item"
	"backend/internal/modules/listing"
	"backend/internal/modules/users"
	"backend/internal/utils"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo             *Repository
	financialSvc     *financial.Service
	listingService   *listing.Service
	chatRepo         *chat.Repository
	containerService *container.Service
	itemService      *item.Service
	userService      *users.Service
}

func NewService(repo *Repository, financialSvc *financial.Service, listingService *listing.Service, chatRepo *chat.Repository, containerService *container.Service, itemService *item.Service, userService *users.Service) *Service {
	return &Service{repo: repo, financialSvc: financialSvc, listingService: listingService, chatRepo: chatRepo, containerService: containerService, itemService: itemService, userService: userService}
}

func (s *Service) GetAll() ([]ListingOrder, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*ListingOrder, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("listing_order/service Listing order ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) GetByUserId(userId pgtype.Int8) ([]ListingOrderWithListing, error) {
	if !userId.Valid || userId.Int64 < 1 {
		return nil, fmt.Errorf("listing_order/service User ID invalide: %d", userId.Int64)
	}
	return s.repo.GetByUserId(userId)
}

func (s *Service) CreateFromRequest(userId int64, req CreateListingOrderRequest) (pgtype.Int8, error) {
	l, err := s.listingService.GetById(pgtype.Int8{Int64: req.ListingId, Valid: true})
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("listing not found: %w", err)
	}

	realPrice, err := l.Price.Float64Value()
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("prix de l'annonce invalide")
	}

	if realPrice.Float64 > 0 {
		return pgtype.Int8{}, fmt.Errorf("cette annonce est payante, utilisez /listing-order/checkout")
	}

	return s.createConfirmedOrder(userId, l)
}

func (s *Service) CreatePaidOrder(userId int64, listingId int64) (pgtype.Int8, error) {
	l, err := s.listingService.GetById(pgtype.Int8{Int64: listingId, Valid: true})
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("listing not found: %w", err)
	}

	return s.createConfirmedOrder(userId, l)
}

func (s *Service) createConfirmedOrder(userId int64, l *listing.Listing) (pgtype.Int8, error) {
	sellerId := l.CreatedBy.Int64

	realPrice, err := l.Price.Float64Value()
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("prix de l'annonce invalide")
	}

	lo := ListingOrder{
		ListingId: l.Id,
		UserId:    pgtype.Int8{Int64: userId, Valid: true},
		Price:     l.Price,
		Status:    Paid,
	}

	id, err := s.repo.Create(lo)
	if err != nil {
		return id, err
	}

	s.financialSvc.GenerateInvoiceForOrder(userId, &sellerId, id.Int64, "listing", realPrice.Float64)

	s.listingService.UpdateStatus(l.Id, "sold")
	s.chatRepo.CloseConversationsByListingId(l.Id.Int64)

	if l.HandoffMode == listing.Locker {
		lockerId, err := s.containerService.FindAvailableLocker(l.CityId)
		if err == nil {
			_, err = s.containerService.CreateLockerAccess(lockerId, l.ItemId, pgtype.Int8{Int64: userId, Valid: true})
			if err == nil {
				s.containerService.UpdateLockerStatus(lockerId, "Occupied")
			}
		}
	}

	points := utils.ActionSaleCompleted.Points
	if l.ItemId.Valid {
		if itm, err := s.itemService.GetById(l.ItemId); err == nil && itm.Weight.Valid {
			if w, err := itm.Weight.Float64Value(); err == nil {
				points += int32(w.Float64 * utils.PointsPerKg)
			}
		}
	}
	s.userService.AddScore(pgtype.Int8{Int64: sellerId, Valid: true}, points, utils.ActionSaleCompleted.Description)

	return id, nil
}

func (s *Service) Create(loDto ListingOrder) (pgtype.Int8, error) {
	id, err := s.repo.Create(loDto)
	if err != nil {
		return id, err
	}

	if loDto.Status == Paid {
		l, _ := s.listingService.GetById(loDto.ListingId)
		sellerId := l.CreatedBy.Int64
		p, _ := loDto.Price.Float64Value()
		s.financialSvc.GenerateInvoiceForOrder(loDto.UserId.Int64, &sellerId, id.Int64, "listing", p.Float64)
		s.listingService.UpdateStatus(loDto.ListingId, "sold")
		s.chatRepo.CloseConversationsByListingId(loDto.ListingId.Int64)
	}

	return id, nil
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("listing_order/service Thread ID invalide: %d", id)
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
