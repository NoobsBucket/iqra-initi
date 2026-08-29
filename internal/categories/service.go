package categories

import (
	"context"
	"errors"
	"strings"
)

var ErrNotFound = errors.New("category not found")

type service interface {
	CreateCategory(ctx context.Context, name, description, imageURL string) (*Category, error)
	GetCategories(ctx context.Context) ([]*Category, error)
	GetCategory(ctx context.Context, id string) (*Category, error)
	UpdateCategory(ctx context.Context, id, name, description, imageURL string) (*Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func slugify(value string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(value), " ", "-"))
}

func (s *Service) CreateCategory(ctx context.Context, name, description, imageURL string) (*Category, error) {
	return s.store.Create(ctx, name, slugify(name), description, imageURL)
}

func (s *Service) GetCategories(ctx context.Context) ([]*Category, error) {
	return s.store.GetAll(ctx)
}

func (s *Service) GetCategory(ctx context.Context, id string) (*Category, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) UpdateCategory(ctx context.Context, id, name, description, imageURL string) (*Category, error) {
	return s.store.Update(ctx, id, name, slugify(name), description, imageURL)
}

func (s *Service) DeleteCategory(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}
