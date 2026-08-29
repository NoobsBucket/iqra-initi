package courses

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("course not found")

type service interface {
	CreateCourse(ctx context.Context, createdBy, title, description, imageURL, currency string, price, discountPrice float64, categoryIDs []string) (*Course, error)
	GetCourses(ctx context.Context) ([]*Course, error)
	GetCourse(ctx context.Context, id string) (*Course, error)
	UpdateCourse(ctx context.Context, id, title, description, imageURL, currency string, price, discountPrice float64) (*Course, error)
	DeleteCourse(ctx context.Context, id string) error
	PublishCourse(ctx context.Context, id string) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreateCourse(ctx context.Context, createdBy, title, description, imageURL, currency string, price, discountPrice float64, categoryIDs []string) (*Course, error) {
	course := &Course{
		CreatedBy:     createdBy,
		Title:         title,
		Description:   description,
		ImageURL:      imageURL,
		Price:         price,
		DiscountPrice: discountPrice,
		Currency:      currency,
	}
	created, err := s.store.Create(ctx, course)
	if err != nil {
		return nil, err
	}
	for _, catID := range categoryIDs {
		if catID == "" {
			continue
		}
		if err := s.store.AddCategory(ctx, created.ID, catID); err != nil {
			return nil, err
		}
	}
	created.CategoryIDs = categoryIDs
	return created, nil
}

func (s *Service) GetCourses(ctx context.Context) ([]*Course, error) {
	return s.store.GetAll(ctx)
}

func (s *Service) GetCourse(ctx context.Context, id string) (*Course, error) {
	c, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	ids, err := s.store.GetCategoryIDs(ctx, id)
	if err != nil {
		return nil, err
	}
	c.CategoryIDs = ids
	return c, nil
}

func (s *Service) UpdateCourse(ctx context.Context, id, title, description, imageURL, currency string, price, discountPrice float64) (*Course, error) {
	course := &Course{
		ID:            id,
		Title:         title,
		Description:   description,
		ImageURL:      imageURL,
		Price:         price,
		DiscountPrice: discountPrice,
		Currency:      currency,
	}
	return s.store.Update(ctx, course)
}

func (s *Service) DeleteCourse(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func (s *Service) PublishCourse(ctx context.Context, id string) error {
	return s.store.Publish(ctx, id)
}