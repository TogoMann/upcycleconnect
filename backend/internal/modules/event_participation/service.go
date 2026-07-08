package eventparticipation

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]EventParticipation, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*EventParticipation, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("eventparticipation/service ID invalide: %d", id)
	}
	return s.repo.GetById(id)
}

func (s *Service) CreateFromRequest(userId int64, req CreateEventParticipationRequest) (pgtype.Int8, error) {
	premium, err := s.repo.IsEventPremium(req.EventId)
	if err != nil {
		return pgtype.Int8{}, err
	}

	if premium {
		eligible, err := s.repo.IsUserEligibleForPremiumEvent(userId)
		if err != nil {
			return pgtype.Int8{}, err
		}
		if !eligible {
			return pgtype.Int8{}, fmt.Errorf("upgrade_required")
		}
	}

	dto := EventParticipation{
		EventId: pgtype.Int8{Int64: req.EventId, Valid: true},
		UserId:  pgtype.Int8{Int64: userId, Valid: true},
	}
	return s.repo.Create(dto)
}

func (s *Service) Create(dto EventParticipation) (pgtype.Int8, error) {
	return s.repo.Create(dto)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("eventparticipation/service ID invalide: %d", id)
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
