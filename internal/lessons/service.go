package lessons

import (
	"context"
	"errors"
)

var ErrNotFound = newErr("lesson not found")

func newErr(msg string) error {
	return errors.New(msg)
}

type service interface {
	CreateLesson(ctx context.Context, lesson *Lesson) (*Lesson, error)
	GetCourseLessons(ctx context.Context, courseID string) ([]*Lesson, error)
	GetLesson(ctx context.Context, id string) (*Lesson, error)
	UpdateLesson(ctx context.Context, lesson *Lesson) (*Lesson, error)
	DeleteLesson(ctx context.Context, id string) error
	ReorderLesson(ctx context.Context, id string, order int) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateLesson(ctx context.Context, lesson *Lesson) (*Lesson, error) {
	return s.store.Create(ctx, lesson)
}

func (s *Service) GetCourseLessons(ctx context.Context, courseID string) ([]*Lesson, error) {
	return s.store.GetByCourse(ctx, courseID)
}

func (s *Service) GetLesson(ctx context.Context, id string) (*Lesson, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) UpdateLesson(ctx context.Context, lesson *Lesson) (*Lesson, error) {
	return s.store.Update(ctx, lesson)
}

func (s *Service) DeleteLesson(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func (s *Service) ReorderLesson(ctx context.Context, id string, order int) error {
	return s.store.Reorder(ctx, id, order)
}
