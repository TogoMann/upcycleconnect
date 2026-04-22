package thread

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]Thread, error) {
	return s.repo.GetAll()
}

func (s *Service) GetById(id pgtype.Int8) (*Thread, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("thread/service Thread ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) Create(threadDto Thread) (pgtype.Int8, error) {
	if strings.TrimSpace(threadDto.Title) == "" || strings.TrimSpace(threadDto.Content) == "" || strings.TrimSpace(string(threadDto.Category)) == "" {
		return pgtype.Int8{}, fmt.Errorf("thread/service Invalid string(s): Missing values.")
	}

	return s.repo.Create(threadDto)
}

func (s *Service) Update(id pgtype.Int8, thread Thread) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("thread/service ID invalide: %d", id)
	}

	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("thread/service thread not found")
	}

	if strings.TrimSpace(thread.Title) == "" || strings.TrimSpace(thread.Content) == "" || strings.TrimSpace(string(thread.Category)) == "" {
		return fmt.Errorf("thread/service Invalid string(s): Missing values.")
	}

	return s.repo.Update(id, thread)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("thread/service Thread ID invalide: %d", id)
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

func (s *Service) Upvote(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("thread/service Thread ID invalide: %d", id)
	}
	return s.repo.Upvote(id)
}

func (s *Service) Downvote(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("thread/service Thread ID invalide: %d", id)
	}
	return s.repo.Downvote(id)
}

func (s *Service) IncrementViews(threadId pgtype.Int8, userId pgtype.Int8) error {
	return s.repo.IncrementViews(threadId, userId)
}
