package auth

import (
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
}

func NewService(userRepo *users.Repository, subRepo *subscriptions.Repository) *Service {
	return &Service{userRepo: userRepo, subRepo: subRepo}
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

	// Auto-subscribe to Free plan
	until := time.Now().AddDate(1, 0, 0) // 1 year
	s.subRepo.Create(subscriptions.Subscription{
		SubscriberId: id,
		Tier:         "Free",
		Price:        0.00,
		Until:        pgtype.Date{Time: until, Valid: true},
	})

	token, err := utils.GenerateJWT(id.Int64, req.Username, string(users.Client))
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		Role:  string(users.Client),
	}, nil
}
