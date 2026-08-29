package blog

import (
	"context"
	"errors"
	"math"
	"strings"
)

var ErrNotFound = errors.New("blog post not found")

type service interface {
	CreatePost(ctx context.Context, post *Post) (*Post, error)
	GetPosts(ctx context.Context) ([]*Post, error)
	GetPost(ctx context.Context, id string) (*Post, error)
	GetPostBySlug(ctx context.Context, slug string) (*Post, error)
	UpdatePost(ctx context.Context, post *Post) (*Post, error)
	DeletePost(ctx context.Context, id string) error
	PublishPost(ctx context.Context, id string) error

	CreateBlogCategory(ctx context.Context, name, description string) (*BlogCategory, error)
	GetBlogCategories(ctx context.Context) ([]*BlogCategory, error)
	GetBlogCategory(ctx context.Context, id string) (*BlogCategory, error)
	UpdateBlogCategory(ctx context.Context, id, name, description string) (*BlogCategory, error)
	DeleteBlogCategory(ctx context.Context, id string) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }
func slugify(value string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(value), " ", "-"))
}
func readTime(content string) int {
	return int(math.Ceil(float64(len(strings.Fields(content))) / 200.0))
}

func (s *Service) CreatePost(ctx context.Context, post *Post) (*Post, error) {
	post.Slug = slugify(post.Title)
	post.ReadTime = readTime(post.Content)
	return s.store.Create(ctx, post)
}
func (s *Service) GetPosts(ctx context.Context) ([]*Post, error) { return s.store.GetAll(ctx) }
func (s *Service) GetPost(ctx context.Context, id string) (*Post, error) {
	return s.store.GetByID(ctx, id)
}
func (s *Service) GetPostBySlug(ctx context.Context, slug string) (*Post, error) {
	post, err := s.store.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	go func() { _ = s.store.IncrementViews(ctx, post.ID) }()
	return post, nil
}
func (s *Service) UpdatePost(ctx context.Context, post *Post) (*Post, error) {
	post.Slug = slugify(post.Title)
	post.ReadTime = readTime(post.Content)
	return s.store.Update(ctx, post)
}
func (s *Service) DeletePost(ctx context.Context, id string) error  { return s.store.Delete(ctx, id) }
func (s *Service) PublishPost(ctx context.Context, id string) error { return s.store.Publish(ctx, id) }

func (s *Service) CreateBlogCategory(ctx context.Context, name, description string) (*BlogCategory, error) {
	slug := slugify(name)
	return s.store.CreateCategory(ctx, name, slug, description)
}

func (s *Service) GetBlogCategories(ctx context.Context) ([]*BlogCategory, error) {
	return s.store.GetAllCategories(ctx)
}

func (s *Service) GetBlogCategory(ctx context.Context, id string) (*BlogCategory, error) {
	return s.store.GetCategoryByID(ctx, id)
}

func (s *Service) UpdateBlogCategory(ctx context.Context, id, name, description string) (*BlogCategory, error) {
	slug := slugify(name)
	return s.store.UpdateCategory(ctx, id, name, slug, description)
}

func (s *Service) DeleteBlogCategory(ctx context.Context, id string) error {
	return s.store.DeleteCategory(ctx, id)
}
