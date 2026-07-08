package stats

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetPublicStats() (*PublicStats, error) {
	return s.repo.GetPublicStats()
}
