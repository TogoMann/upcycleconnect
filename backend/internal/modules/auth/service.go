package auth

import (
	"backend/internal/modules/companies"
	"backend/internal/modules/plans"
	"backend/internal/modules/settings"
	"backend/internal/modules/subscriptions"
	"backend/internal/modules/users"
	"backend/internal/utils"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func frontendURL() string {
	url := os.Getenv("FRONTEND_URL")
	if url == "" {
		url = "http://localhost:5173"
	}
	return url
}

type Service struct {
	userRepo     *users.Repository
	userService  *users.Service
	subRepo      *subscriptions.Repository
	planRepo     *plans.Repository
	compRepo     *companies.Repository
	settingsRepo *settings.Repository
}

func NewService(userRepo *users.Repository, userService *users.Service, subRepo *subscriptions.Repository, planRepo *plans.Repository, compRepo *companies.Repository, settingsRepo *settings.Repository) *Service {
	return &Service{userRepo: userRepo, userService: userService, subRepo: subRepo, planRepo: planRepo, compRepo: compRepo, settingsRepo: settingsRepo}
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
	open, err := s.settingsRepo.IsRegistrationOpen()
	if err == nil && !open {
		return nil, errors.New("les inscriptions sont temporairement fermées")
	}

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

	s.userService.AddScore(id, utils.ActionRegistration.Points, utils.ActionRegistration.Description)

	tier := "Free"
	var price float64 = 0.00
	until := time.Now().AddDate(1, 0, 0)

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

	if strings.EqualFold(tier, "Pro") {
		if req.Siret == "" {
			return nil, errors.New("le numéro SIRET est obligatoire pour le plan Pro")
		}
		if !utils.VerifySiret(req.Siret) {
			return nil, errors.New("numéro SIRET invalide (doit contenir 14 chiffres)")
		}

		company, err := s.compRepo.GetBySiret(req.Siret)
		if err != nil {
			return nil, err
		}

		var companyId pgtype.Int8
		if company == nil {
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

	go utils.SendConfirmationEmail(req.Email, req.Username)

	return &LoginResponse{
		Token: token,
		Role:  string(users.Client),
	}, nil
}

func (s *Service) AdminRequestPasswordReset(userId int64) error {
	user, err := s.userRepo.GetById(pgtype.Int8{Int64: userId, Valid: true})
	if err != nil {
		return err
	}

	token := utils.GenerateSecureToken()
	expiresAt := pgtype.Timestamp{Time: time.Now().Add(time.Hour * 2), Valid: true}

	err = s.userRepo.CreateResetToken(userId, token, expiresAt)
	if err != nil {
		return err
	}

	subject := "Réinitialisation de votre mot de passe - UpCycleConnect"
	resetUrl := frontendURL() + "/auth/reset-password?token=" + token
	body := fmt.Sprintf(`
		<h1>Bonjour %s !</h1>
		<p>Un administrateur a initié une réinitialisation de votre mot de passe.</p>
		<p>Veuillez cliquer sur le lien ci-dessous pour choisir un nouveau mot de passe :</p>
		<p><a href="%s">%s</a></p>
		<p>Ce lien expirera dans 2 heures.</p>
		<br/>
		<p>L'équipe UpCycleConnect</p>
	`, user.FirstName, resetUrl, resetUrl)

	return utils.SendEmail(user.Email, subject, body)
}

func (s *Service) RequestPasswordReset(email string) error {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	token := utils.GenerateSecureToken()
	expiresAt := pgtype.Timestamp{Time: time.Now().Add(time.Hour * 2), Valid: true}

	err = s.userRepo.CreateResetToken(user.Id.Int64, token, expiresAt)
	if err != nil {
		return err
	}

	subject := "Réinitialisation de votre mot de passe - UpCycleConnect"
	resetUrl := frontendURL() + "/auth/reset-password?token=" + token
	body := fmt.Sprintf(`
		<h1>Bonjour %s !</h1>
		<p>Vous avez demandé la réinitialisation de votre mot de passe.</p>
		<p>Veuillez cliquer sur le lien ci-dessous pour choisir un nouveau mot de passe :</p>
		<p><a href="%s">%s</a></p>
		<p>Ce lien expirera dans 2 heures. Si vous n'êtes pas à l'origine de cette demande, ignorez cet e-mail.</p>
		<br/>
		<p>L'équipe UpCycleConnect</p>
	`, user.FirstName, resetUrl, resetUrl)

	return utils.SendEmail(user.Email, subject, body)
}

func (s *Service) ResetPassword(token, newPassword string) error {
	userId, err := s.userRepo.GetUserIdByResetToken(token)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdatePassword(userId, string(hash))
	if err != nil {
		return err
	}

	return s.userRepo.DeleteResetToken(token)
}
