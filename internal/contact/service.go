package contact

import "context"

type service interface {
	SendMessage(ctx context.Context, name, email, phone, subject, message string) (*Message, error)
	GetMessages(ctx context.Context) ([]*Message, error)
	MarkRead(ctx context.Context, id string) error
	DeleteMessage(ctx context.Context, id string) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) SendMessage(ctx context.Context, name, email, phone, subject, message string) (*Message, error) {
	m := &Message{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Subject: subject,
		Message: message,
	}
	return s.store.Create(ctx, m)
}

func (s *Service) GetMessages(ctx context.Context) ([]*Message, error) {
	return s.store.GetAll(ctx)
}

func (s *Service) MarkRead(ctx context.Context, id string) error {
	return s.store.MarkRead(ctx, id)
}

func (s *Service) DeleteMessage(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}
