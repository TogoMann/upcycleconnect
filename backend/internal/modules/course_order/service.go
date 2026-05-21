package courseorder

import (
	"backend/internal/modules/course"
	"backend/internal/modules/financial"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo         *Repository
	financialSvc *financial.Service
	courseSvc    *course.Service
}

func NewService(repo *Repository, financialSvc *financial.Service, courseSvc *course.Service) *Service {
	return &Service{repo: repo, financialSvc: financialSvc, courseSvc: courseSvc}
}

func (s *Service) GetAll() ([]CourseOrder, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*CourseOrder, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("courseorder/service ID invalide: %d", id)
	}
	return s.repo.GetById(id)
}

func (s *Service) GetByUserId(userId pgtype.Int8) ([]CourseOrderWithCourse, error) {
	if !userId.Valid || userId.Int64 < 1 {
		return nil, fmt.Errorf("courseorder/service User ID invalide: %d", userId.Int64)
	}
	return s.repo.GetByUserId(userId)
}

func (s *Service) CreateFromRequest(userId int64, req CreateCourseOrderRequest) (pgtype.Int8, error) {
	var price pgtype.Numeric
	price.Scan(fmt.Sprintf("%.2f", req.Price))

	c, err := s.courseSvc.GetById(pgtype.Int8{Int64: req.CourseId, Valid: true})
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("course not found: %w", err)
	}
	sellerId := c.CreatedBy.Int64

	lo := CourseOrder{
		CourseId: pgtype.Int8{Int64: req.CourseId, Valid: true},
		BuyerId:  pgtype.Int8{Int64: userId, Valid: true},
		Price:    price,
	}

	id, err := s.repo.Create(lo)
	if err != nil {
		return id, err
	}

	s.financialSvc.GenerateInvoiceForOrder(userId, &sellerId, id.Int64, "course", req.Price)

	return id, nil
}

func (s *Service) Create(dto CourseOrder) (pgtype.Int8, error) {
	return s.repo.Create(dto)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("courseorder/service ID invalide: %d", id)
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
