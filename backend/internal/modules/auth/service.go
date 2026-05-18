package auth

import (
	"backend/internal/modules/companies"
	"backend/internal/modules/plans"
	"backend/internal/modules/subscriptions"
	"backend/internal/modules/users"
	"backend/internal/utils"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo *users.Repository
	subRepo  *subscriptions.Repository
	planRepo *plans.Repository
	compRepo *companies.Repository
}

func NewService(userRepo *users.Repository, subRepo *subscriptions.Repository, planRepo *plans.Repository, compRepo *companies.Repository) *Service {
	return &Service{userRepo: userRepo, subRepo: subRepo, planRepo: planRepo, compRepo: compRepo}
}

func (s *Service) Login(username, password string) (*LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.Id.Int64, user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		Role:  string(user.Role),
	}, nil
}

func (s *Service) Register(req RegisterRequest) (*LoginResponse, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" || req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("tous les champs sont obligatoires")
	}

	existing, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("ce nom d'utilisateur est déjà pris")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("erreur lors du chiffrement du mot de passe")
	}

	newUser := users.User{
		Username:           req.Username,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		Email:              req.Email,
		PasswordHash:       string(hash),
		Role:               users.Client,
		LanguagePreference: "fr",
	}

	id, err := s.userRepo.Create(newUser)
	if err != nil {
		return nil, errors.New("erreur lors de la création du compte")
	}

	// Determine plan to subscribe to
	tier := "Free"
	var price float64 = 0.00
	until := time.Now().AddDate(1, 0, 0) // 1 year default for free

	if req.PlanId > 0 {
		plan, err := s.planRepo.GetById(pgtype.Int8{Int64: req.PlanId, Valid: true})
		if err == nil && plan != nil {
			tier = plan.Name
			p, _ := plan.Price.Float64Value()
			price = p.Float64
			if plan.BillingCycle == "monthly" {
				until = time.Now().AddDate(0, 1, 0)
			}
		}
	}

	// If plan is Pro, update user role to 'pro' and require SIRET
	if strings.EqualFold(tier, "Pro") {
		if req.Siret == "" {
			return nil, errors.New("le numéro SIRET est obligatoire pour le plan Pro")
		}
		if !utils.VerifySiret(req.Siret) {
			return nil, errors.New("numéro SIRET invalide (doit contenir 14 chiffres)")
		}

		// Handle company creation or linking
		company, err := s.compRepo.GetBySiret(req.Siret)
		if err != nil {
			return nil, err
		}

		var companyId pgtype.Int8
		if company == nil {
			// Create new company
			newComp := companies.Company{
				Siret: req.Siret,
			}
			companyId, err = s.compRepo.Create(newComp)
			if err != nil {
				return nil, err
			}
		} else {
			companyId = company.Id
		}

		newUser.Id = id
		newUser.Role = users.Pro
		newUser.CompanyId = companyId
		s.userRepo.Update(id, newUser)
	}

	s.subRepo.Create(subscriptions.Subscription{
		SubscriberId: id,
		Tier:         tier,
		Price:        price,
		Until:        pgtype.Date{Time: until, Valid: true},
	})

	token, err := utils.GenerateJWT(id.Int64, req.Username, string(newUser.Role))
	if err != nil {
		return nil, err
	}

	// Send confirmation email asynchronously
	go utils.SendConfirmationEmail(req.Email, req.Username)

	return &LoginResponse{
		Token: token,
		Role:  string(users.Client),
	}, nil
}
