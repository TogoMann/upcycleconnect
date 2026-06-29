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

	n := Notification{
		Titre:   req.Titre,
		Message: req.Message,
		Cible:   req.Cible,
		Envoyes: 0,
	}

	id, err := s.repo.Create(n)
	if err != nil {
		return n, err
	}
	n.Id = id
	n.Date = "A l'instant"

	// Run sending loop in a background goroutine to prevent blocking the HTTP response
	go func(notificationId int64, targets []string, title, msg string) {
		successCount := 0
		for _, email := range targets {
			err := utils.SendEmail(email, title, msg)
			if err == nil {
				successCount++
			} else {
				log.Printf("Erreur envoi email notification à %s: %v", email, err)
			}
		}
		
		// Update the final success count in the database
		_ = s.repo.UpdateEnvoyes(notificationId, successCount)
	}(id, emails, req.Titre, req.Message)
	
	return n, nil
}
