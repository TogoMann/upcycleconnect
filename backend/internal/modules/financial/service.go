package financial

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCommissions() ([]Commission, error) {
	return s.repo.GetCommissions()
}

func (s *Service) GetFinancier() (*FinancierData, error) {
	return s.repo.GetFinancier()
}

func (s *Service) GetReport() (*FinancialReport, error) {
	return s.repo.GetReport()
}
