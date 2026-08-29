package enrollment

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrAlreadyEnrolled = errors.New("already enrolled in this course")
	ErrNotFound        = errors.New("enrollment not found")
)

type service interface {
	EnrollUser(ctx context.Context, userID, courseID, fullName, whatsapp, notes string) (*Enrollment, error)
	GetUserEnrollments(ctx context.Context, userID string) ([]*Enrollment, error)
	GetEnrollment(ctx context.Context, id string) (*Enrollment, error)
	GetCourseEnrollments(ctx context.Context, courseID string) ([]*Enrollment, error)
	UpdateProgress(ctx context.Context, id string, progress int) error
	UpdatePaymentStatus(ctx context.Context, id, status string) error
	CompleteEnrollment(ctx context.Context, id string) error
	DeleteEnrollment(ctx context.Context, id string) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) EnrollUser(ctx context.Context, userID, courseID, fullName, whatsapp, notes string) (*Enrollment, error) {
	enrolled, err := s.store.IsEnrolled(ctx, userID, courseID)
	if err != nil {
		return nil, fmt.Errorf("check enrolled: %w", err)
	}
	if enrolled {
		return nil, ErrAlreadyEnrolled
	}

	paymentMethod := "manual"
	e := &Enrollment{
		UserID:        userID,
		CourseID:      courseID,
		FullName:      fullName,
		Whatsapp:      whatsapp,
		Notes:         notes,
		PaymentMethod: &paymentMethod,
	}

	result, err := s.store.Enroll(ctx, e)
	if err != nil {
		return nil, fmt.Errorf("enroll: %w", err)
	}
	return result, nil
}

func (s *Service) GetUserEnrollments(ctx context.Context, userID string) ([]*Enrollment, error) {
	return s.store.GetByUser(ctx, userID)
}

func (s *Service) GetEnrollment(ctx context.Context, id string) (*Enrollment, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) GetCourseEnrollments(ctx context.Context, courseID string) ([]*Enrollment, error) {
	return s.store.GetByCourse(ctx, courseID)
}

func (s *Service) UpdateProgress(ctx context.Context, id string, progress int) error {
	if progress < 0 || progress > 100 {
		return errors.New("progress must be between 0 and 100")
	}
	return s.store.UpdateProgress(ctx, id, progress)
}

func (s *Service) UpdatePaymentStatus(ctx context.Context, id, status string) error {
	return s.store.UpdatePaymentStatus(ctx, id, status)
}

func (s *Service) CompleteEnrollment(ctx context.Context, id string) error {
	return s.store.Complete(ctx, id)
}

func (s *Service) DeleteEnrollment(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}
