package cart

import (
	listingorder "backend/internal/modules/listing_order"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo         *Repository
	orderService *listingorder.Service
}

func NewService(repo *Repository, orderService *listingorder.Service) *Service {
	return &Service{repo: repo, orderService: orderService}
}

func (s *Service) GetByUserId(userId pgtype.Int8) ([]CartItemWithListing, error) {
	return s.repo.GetByUserId(userId)
}

func (s *Service) Add(userId, listingId pgtype.Int8) error {
	return s.repo.Add(userId, listingId)
}

func (s *Service) Remove(userId, listingId pgtype.Int8) error {
	return s.repo.Remove(userId, listingId)
}

func (s *Service) Checkout(userId pgtype.Int8) error {
	items, err := s.repo.GetByUserId(userId)
	if err != nil {
		return err
	}

	for _, item := range items {
		order := listingorder.ListingOrder{
			ListingId: item.ListingId,
			UserId:    userId,
			Price:     item.Listing.Price,
			Status:    listingorder.Paid, 
		}
		_, err := s.orderService.Create(order)
		if err != nil {
			return err
		}
	}

	return s.repo.Clear(userId)
}
