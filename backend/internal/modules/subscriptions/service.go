package subscriptions

import (
	"backend/internal/modules/plans"
	"backend/internal/modules/users"
	"backend/internal/utils"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo     *Repository
	userRepo *users.Repository
	planRepo *plans.Repository
}

func NewService(repo *Repository, userRepo *users.Repository, planRepo *plans.Repository) *Service {
	return &Service{repo: repo, userRepo: userRepo, planRepo: planRepo}
}

func (s *Service) ChoosePlan(userId int64, planId int64, siret string) error {
	// 1. Fetch plan
	p, err := s.planRepo.GetById(pgtype.Int8{Int64: planId, Valid: true})
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	// 2. Verification logic
	role := users.Client
	if p.Name == "Pro" {
		if !utils.VerifySiret(siret) {
			return fmt.Errorf("numéro SIRET invalide")
		}
		role = users.Pro
	}

	// 3. Update user
	u, err := s.userRepo.GetById(pgtype.Int8{Int64: userId, Valid: true})
	if err != nil {
		return err
	}

	u.Role = role
	if role == users.Pro {
		u.Siret = pgtype.Text{String: siret, Valid: true}
	} else {
		u.Siret = pgtype.Text{Valid: false}
	}

	err = s.userRepo.Update(u.Id, *u)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// 4. Create subscription
	price, _ := p.Price.Float64Value()
	sub := Subscription{
		SubscriberId: pgtype.Int8{Int64: userId, Valid: true},
		Price:        price.Float64,
		Tier:         p.Name,
		Until:        pgtype.Date{Time: time.Now().AddDate(0, 1, 0), Valid: true}, // 1 month
	}

	_, err = s.repo.Create(sub)
	return err
}

func (s *Service) GetAllAbonnements() ([]AbonnementFrontend, error) {
	return s.repo.GetAllAbonnements()
}

func (s *Service) GetAll() ([]Subscription, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Subscription, error) {
	return s.repo.GetById(id)
}

func (s *Service) GetActiveByUserId(userId pgtype.Int8) (*Subscription, error) {
	return s.repo.GetActiveByUserId(userId)
}

func (s *Service) Create(sub Subscription) (pgtype.Int8, error) {
	return s.repo.Create(sub)
}

func (s *Service) Update(id pgtype.Int8, sub Subscription) error {
	return s.repo.Update(id, sub)
}

func (s *Service) Delete(id pgtype.Int8) error {
	return s.repo.Delete(id)
}
