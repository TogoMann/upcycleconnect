package listingorder

import (
	"backend/internal/modules/chat"
	"backend/internal/modules/financial"
	"backend/internal/modules/listing"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo           *Repository
	financialSvc   *financial.Service
	listingService *listing.Service
	chatRepo       *chat.Repository
}

func NewService(repo *Repository, financialSvc *financial.Service, listingService *listing.Service, chatRepo *chat.Repository) *Service {
	return &Service{repo: repo, financialSvc: financialSvc, listingService: listingService, chatRepo: chatRepo}
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
	var price pgtype.Numeric
	price.Scan(fmt.Sprintf("%.2f", req.Price))

	l, err := s.listingService.GetById(pgtype.Int8{Int64: req.ListingId, Valid: true})
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("listing not found: %w", err)
	}
	sellerId := l.CreatedBy.Int64

	lo := ListingOrder{
		ListingId: pgtype.Int8{Int64: req.ListingId, Valid: true},
		UserId:    pgtype.Int8{Int64: userId, Valid: true},
		Price:     price,
		Status:    Paid,
	}

	id, err := s.repo.Create(lo)
	if err != nil {
		return id, err
	}

	s.financialSvc.GenerateInvoiceForOrder(userId, &sellerId, id.Int64, "listing", req.Price)

	s.listingService.UpdateStatus(pgtype.Int8{Int64: req.ListingId, Valid: true}, "sold")
	s.chatRepo.CloseConversationsByListingId(req.ListingId)

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
