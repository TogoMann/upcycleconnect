package settings

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get() (*PlatformSettings, error) {
	return s.repo.Get()
}

func (s *Service) Update(settings PlatformSettings) error {
	return s.repo.Update(settings)
}

func (s *Service) GetPublic() (*PublicSettings, error) {
	return s.repo.GetPublic()
}

func (s *Service) IsRegistrationOpen() (bool, error) {
	return s.repo.IsRegistrationOpen()
}
