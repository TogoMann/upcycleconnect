package cart

import (
	courseorder "backend/internal/modules/course_order"
	eventparticipation "backend/internal/modules/event_participation"
	"backend/internal/modules/financial"
	listingorder "backend/internal/modules/listing_order"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo          *Repository
	orderService  *listingorder.Service
	eventPartRepo *eventparticipation.Repository
	courseOrdRepo *courseorder.Repository
	financialSvc  *financial.Service
}

func NewService(repo *Repository, orderService *listingorder.Service, eventPartRepo *eventparticipation.Repository, courseOrdRepo *courseorder.Repository, financialSvc *financial.Service) *Service {
	return &Service{
		repo:          repo,
		orderService:  orderService,
		eventPartRepo: eventPartRepo,
		courseOrdRepo: courseOrdRepo,
		financialSvc:  financialSvc,
	}
}

func (s *Service) GetByUserId(userId pgtype.Int8) ([]CartItemDetailed, error) {
	return s.repo.GetByUserId(userId)
}

func (s *Service) Add(userId, listingId, eventId, courseId pgtype.Int8) error {
	return s.repo.Add(userId, listingId, eventId, courseId)
}

func (s *Service) Remove(userId, listingId, eventId, courseId pgtype.Int8) error {
	return s.repo.Remove(userId, listingId, eventId, courseId)
}

func (s *Service) Checkout(userId pgtype.Int8) error {
	items, err := s.repo.GetByUserId(userId)
	if err != nil {
		return fmt.Errorf("checkout: get cart failed: %w", err)
	}

	fmt.Printf("Starting checkout for user %d with %d items\n", userId.Int64, len(items))

	for _, item := range items {
		if item.EventId.Valid {
			fmt.Printf("Processing event %d\n", item.EventId.Int64)
			participation := eventparticipation.EventParticipation{
				EventId: item.EventId,
				UserId:  userId,
			}
			_, err := s.eventPartRepo.Create(participation)
			if err != nil {
				fmt.Printf("Warning: event participation for event %d might already exist or failed: %v\n", item.EventId.Int64, err)
			} else if item.Event != nil {
				price, _ := item.Event.Price.Float64Value()
				sellerId := item.Event.CreatedBy.Int64
				s.financialSvc.GenerateInvoiceForOrder(userId.Int64, &sellerId, item.EventId.Int64, "event", price.Float64)
			}
		} else if item.CourseId.Valid {
			fmt.Printf("Processing course %d\n", item.CourseId.Int64)
			order := courseorder.CourseOrder{
				CourseId: item.CourseId,
				BuyerId:  userId,
				Price:    item.Course.Price,
			}
			id, err := s.courseOrdRepo.Create(order)
			if err != nil {
				fmt.Printf("Warning: course order for course %d might already exist or failed: %v\n", item.CourseId.Int64, err)
			} else if item.Course != nil {
				price, _ := item.Course.Price.Float64Value()
				sellerId := item.Course.CreatedBy.Int64
				s.financialSvc.GenerateInvoiceForOrder(userId.Int64, &sellerId, id.Int64, "course", price.Float64)
			}
		}
	}

	fmt.Printf("Clearing direct pay items for user %d\n", userId.Int64)
	if err := s.repo.ClearDirectPay(userId); err != nil {
		return fmt.Errorf("checkout: clear cart failed: %w", err)
	}

	return nil
}
