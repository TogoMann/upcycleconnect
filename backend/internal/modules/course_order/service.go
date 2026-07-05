package courseorder

import (
	"backend/internal/modules/course"
	"backend/internal/modules/financial"
	"backend/internal/modules/users"
	"backend/internal/utils"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo         *Repository
	financialSvc *financial.Service
	courseSvc    *course.Service
	userService  *users.Service
}

func NewService(repo *Repository, financialSvc *financial.Service, courseSvc *course.Service, userService *users.Service) *Service {
	return &Service{repo: repo, financialSvc: financialSvc, courseSvc: courseSvc, userService: userService}
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

func (s *Service) checkCapacity(courseId pgtype.Int8, maxCapacity pgtype.Int4) error {
	if !maxCapacity.Valid {
		return nil
	}
	count, err := s.repo.CountByCourseId(courseId)
	if err != nil {
		return err
	}
	if count >= int64(maxCapacity.Int32) {
		return fmt.Errorf("cette formation a atteint sa capacité maximale")
	}
	return nil
}

func (s *Service) CreateFromRequest(userId int64, req CreateCourseOrderRequest) (pgtype.Int8, error) {
	c, err := s.courseSvc.GetById(pgtype.Int8{Int64: req.CourseId, Valid: true})
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("course not found: %w", err)
	}
	sellerId := c.CreatedBy.Int64

	if err := s.checkCapacity(pgtype.Int8{Int64: req.CourseId, Valid: true}, c.MaxCapacity); err != nil {
		return pgtype.Int8{}, err
	}

	realPrice, err := c.Price.Float64Value()
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("prix de la formation invalide")
	}

	lo := CourseOrder{
		CourseId: pgtype.Int8{Int64: req.CourseId, Valid: true},
		BuyerId:  pgtype.Int8{Int64: userId, Valid: true},
		Price:    c.Price,
	}

	id, err := s.repo.Create(lo)
	if err != nil {
		return id, err
	}

	s.financialSvc.GenerateInvoiceForOrder(userId, &sellerId, id.Int64, "course", realPrice.Float64)
	s.userService.AddScore(lo.BuyerId, utils.ActionAtelierParticipation.Points, utils.ActionAtelierParticipation.Description)

	return id, nil
}

func (s *Service) Create(dto CourseOrder) (pgtype.Int8, error) {
	c, err := s.courseSvc.GetById(dto.CourseId)
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("course not found: %w", err)
	}

	if err := s.checkCapacity(dto.CourseId, c.MaxCapacity); err != nil {
		return pgtype.Int8{}, err
	}

	id, err := s.repo.Create(dto)
	if err != nil {
		return id, err
	}

	s.userService.AddScore(dto.BuyerId, utils.ActionAtelierParticipation.Points, utils.ActionAtelierParticipation.Description)

	return id, nil
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
