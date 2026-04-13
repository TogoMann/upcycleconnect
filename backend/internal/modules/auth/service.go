package auth

import (
	"backend/internal/modules/users"
	"backend/internal/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo *users.Repository
}

func NewService(userRepo *users.Repository) *Service {
	return &Service{userRepo: userRepo}
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

	token, err := utils.GenerateJWT(user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		Role:  string(user.Role),
	}, nil
}
