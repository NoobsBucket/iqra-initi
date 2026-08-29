package notifications

import "context"

type service interface {
	CreateNotification(ctx context.Context, userID, title, message, notifType string) (*Notification, error)
	GetUserNotifications(ctx context.Context, userID string) ([]*Notification, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, userID string) error
	DeleteNotification(ctx context.Context, id string) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateNotification(ctx context.Context, userID, title, message, notifType string) (*Notification, error) {
	return s.store.Create(ctx, userID, title, message, notifType)
}

func (s *Service) GetUserNotifications(ctx context.Context, userID string) ([]*Notification, error) {
	return s.store.GetByUser(ctx, userID)
}

func (s *Service) MarkRead(ctx context.Context, id string) error {
	return s.store.MarkRead(ctx, id)
}

func (s *Service) MarkAllRead(ctx context.Context, userID string) error {
	return s.store.MarkAllRead(ctx, userID)
}

func (s *Service) DeleteNotification(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}
