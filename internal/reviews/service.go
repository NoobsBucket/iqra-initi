package reviews

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("review not found")

type service interface {
	CreateReview(ctx context.Context, userID, courseID, reviewText string, rating float64) (*Review, error)
	GetCourseReviews(ctx context.Context, courseID string) ([]*Review, error)
	UpdateReview(ctx context.Context, id, reviewText string, rating float64) (*Review, error)
	DeleteReview(ctx context.Context, id string) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) CreateReview(ctx context.Context, userID, courseID, reviewText string, rating float64) (*Review, error) {
	if rating < 1 || rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}
	return s.store.Create(ctx, userID, courseID, reviewText, rating)
}
func (s *Service) GetCourseReviews(ctx context.Context, courseID string) ([]*Review, error) {
	return s.store.GetByCourse(ctx, courseID)
}
func (s *Service) UpdateReview(ctx context.Context, id, reviewText string, rating float64) (*Review, error) {
	return s.store.Update(ctx, id, reviewText, rating)
}
func (s *Service) DeleteReview(ctx context.Context, id string) error { return s.store.Delete(ctx, id) }
