package settings

import "context"

type service interface {
	GetSettings(ctx context.Context) (*Settings, error)
	UpdateSettings(ctx context.Context, s *Settings) (*Settings, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetSettings(ctx context.Context) (*Settings, error) {
	return s.store.Get(ctx)
}

func (s *Service) UpdateSettings(ctx context.Context, settings *Settings) (*Settings, error) {
	return s.store.Update(ctx, settings)
}
