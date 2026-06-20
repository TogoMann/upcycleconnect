package notifications

import (
	"backend/internal/utils"
	"log"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]Notification, error) {
	return s.repo.GetAll()
}

func (s *Service) SendBulkEmail(req NotificationRequest) (Notification, error) {
	emails, err := s.repo.GetEmailsByRole(req.Cible)
	if err != nil {
		return Notification{}, err
	}

	envoyes := 0
	for _, email := range emails {
		err := utils.SendEmail(email, req.Titre, req.Message)
		if err == nil {
			envoyes++
		} else {
			log.Printf("Erreur envoi email notification à %s: %v", email, err)
		}
	}

	n := Notification{
		Titre:   req.Titre,
		Message: req.Message,
		Cible:   req.Cible,
		Envoyes: envoyes,
	}

	id, err := s.repo.Create(n)
	if err != nil {
		return n, err
	}
	n.Id = id
	n.Date = "A l'instant"
	
	return n, nil
}
