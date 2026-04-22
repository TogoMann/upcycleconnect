package post

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

func (s *Service) GetAll() ([]Post, error) {
	return s.repo.GetAll()
}
func (s *Service) GetById(id pgtype.Int8) (*Post, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("post/service Post ID invalide: %d", id)
	}

	return s.repo.GetById(id)
}

func (s *Service) GetThreadPosts(id pgtype.Int8) ([]ThreadPosts, error) {
	if !id.Valid || id.Int64 < 1 {
		return nil, fmt.Errorf("post/service Thread ID invalide: %d", id)
	}

	return s.repo.GetThreadPosts(id)
}

func (s *Service) Create(postDto Post) (pgtype.Int8, error) {
	if strings.TrimSpace(postDto.Content) == "" {
		return pgtype.Int8{}, fmt.Errorf("post/service Invalid string(s): Missing values.")
	}

	return s.repo.Create(postDto)
}

func (s *Service) UpdateContent(id pgtype.Int8, content string) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("post/service Post ID invalide: %d", id)
	}
	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("post/service Post not found")
	}

	return s.repo.UpdateContent(id, content)
}

func (s *Service) Delete(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("post/service Post ID invalide: %d", id)
	}

	exists, err := s.repo.ExistsById(id)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	hasReplies, err := s.repo.HasReplies(id)
	if err != nil {
		return err
	}

	if hasReplies {
		return s.repo.SoftDelete(id)
	}

	return s.repo.Delete(id)
}

func (s *Service) Upvote(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("post/service Post ID invalide: %d", id)
	}
	return s.repo.Upvote(id)
}

func (s *Service) Downvote(id pgtype.Int8) error {
	if !id.Valid || id.Int64 < 1 {
		return fmt.Errorf("post/service Post ID invalide: %d", id)
	}
	return s.repo.Downvote(id)
}
