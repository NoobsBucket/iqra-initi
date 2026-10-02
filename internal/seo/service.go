package seo

import "context"

type service interface {
	GetPage(ctx context.Context, path string) (*Page, error)
	GetPages(ctx context.Context) ([]*Page, error)
	CreatePage(ctx context.Context, page *Page) (*Page, error)
	UpdatePage(ctx context.Context, page *Page) (*Page, error)
	DeletePage(ctx context.Context, id string) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) GetPage(ctx context.Context, path string) (*Page, error) {
	return s.store.GetByPath(ctx, path)
}

func (s *Service) GetPages(ctx context.Context) ([]*Page, error) {
	return s.store.GetAll(ctx)
}

func (s *Service) CreatePage(ctx context.Context, page *Page) (*Page, error) {
	return s.store.Create(ctx, page)
}

func (s *Service) UpdatePage(ctx context.Context, page *Page) (*Page, error) {
	return s.store.Update(ctx, page)
}

func (s *Service) DeletePage(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}
