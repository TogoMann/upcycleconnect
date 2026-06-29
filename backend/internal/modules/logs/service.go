package logs

import (
	"backend/internal/middlewares"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	repo *Repository
}

var GlobalService *Service

func NewService(repo *Repository) *Service {
	s := &Service{repo: repo}
	GlobalService = s
	return s
}

func (s *Service) Log(utilisateur, action, ressource, ip, niveau string) error {
	if ip == "" {
		ip = "127.0.0.1"
	}
	if niveau == "" {
		niveau = "info"
	}
	return s.repo.Create(utilisateur, action, ressource, ip, niveau)
}

func (s *Service) GetAll() ([]Log, error) {
	return s.repo.GetAll()
}

func Add(utilisateur, action, ressource, ip, niveau string) {
	if GlobalService != nil {
		_ = GlobalService.Log(utilisateur, action, ressource, ip, niveau)
	}
}

func AddFromRequest(r *http.Request, action, ressource, niveau string) {
	if GlobalService == nil {
		return
	}

	username := "Visiteur"
	claims, ok := r.Context().Value(middlewares.ClaimsKey).(jwt.MapClaims)
	if ok {
		if u, exists := claims["username"].(string); exists && u != "" {
			username = u
		} else if email, exists := claims["email"].(string); exists && email != "" {
			username = email
		} else if sub, exists := claims["sub"].(float64); exists {
			username = fmt.Sprintf("User #%d", int64(sub))
		}
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	_ = GlobalService.Log(username, action, ressource, ip, niveau)
}
