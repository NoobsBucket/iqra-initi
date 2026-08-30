package users

import (
	"context"
	"errors"
)

var ErrInvalidRole = errors.New("invalid role")

var validRoles = map[string]bool{
	"user":        true,
	"admin":       true,
	"instructor":  true,
	"blog_writer": true,
}

type service interface {
	GetUsers(ctx context.Context, search string) ([]*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	AssignRole(ctx context.Context, id, role string) (*User, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetUsers(ctx context.Context, search string) ([]*User, error) {
	return s.store.GetAll(ctx, search)
}

func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) AssignRole(ctx context.Context, id, role string) (*User, error) {
	if !validRoles[role] {
		return nil, ErrInvalidRole
	}
	return s.store.UpdateRole(ctx, id, role)
}
