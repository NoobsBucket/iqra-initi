package blog

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID              string     `json:"id"`
	CreatedBy       string     `json:"created_by"`
	CategoryID      string     `json:"category_id"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Excerpt         string     `json:"excerpt"`
	Content         string     `json:"content"`
	CoverImage      string     `json:"cover_image"`
	IsPublished     bool       `json:"is_published"`
	IsFeatured      bool       `json:"is_featured"`
	MetaTitle       string     `json:"meta_title"`
	MetaDescription string     `json:"meta_description"`
	MetaKeywords    string     `json:"meta_keywords"`
	Views           int        `json:"views"`
	ReadTime        int        `json:"read_time"`
	PublishedAt     *time.Time `json:"published_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Store interface {
	Create(ctx context.Context, post *Post) (*Post, error)
	GetAll(ctx context.Context) ([]*Post, error)
	GetByID(ctx context.Context, id string) (*Post, error)
	GetBySlug(ctx context.Context, slug string) (*Post, error)
	Update(ctx context.Context, post *Post) (*Post, error)
	Delete(ctx context.Context, id string) error
	Publish(ctx context.Context, id string) error
	IncrementViews(ctx context.Context, id string) error

	CreateCategory(ctx context.Context, name, slug, description string) (*BlogCategory, error)
	GetAllCategories(ctx context.Context) ([]*BlogCategory, error)
	GetCategoryByID(ctx context.Context, id string) (*BlogCategory, error)
	UpdateCategory(ctx context.Context, id, name, slug, description string) (*BlogCategory, error)
	DeleteCategory(ctx context.Context, id string) error
}

type store struct{ db *pgxpool.Pool }

func NewStore(db *pgxpool.Pool) Store { return &store{db: db} }

const postColumns = `id, created_by, category_id, title, slug, COALESCE(excerpt, ''), content, COALESCE(cover_image, ''), is_published, is_featured, COALESCE(meta_title, ''), COALESCE(meta_description, ''), COALESCE(meta_keywords, ''), views, read_time, published_at, created_at, updated_at`

func scanPost(row interface{ Scan(...any) error }, post *Post) error {
	return row.Scan(&post.ID, &post.CreatedBy, &post.CategoryID, &post.Title, &post.Slug, &post.Excerpt, &post.Content, &post.CoverImage, &post.IsPublished, &post.IsFeatured, &post.MetaTitle, &post.MetaDescription, &post.MetaKeywords, &post.Views, &post.ReadTime, &post.PublishedAt, &post.CreatedAt, &post.UpdatedAt)
}

func (s *store) Create(ctx context.Context, post *Post) (*Post, error) {
	err := scanPost(s.db.QueryRow(ctx, `INSERT INTO blog_posts (created_by, category_id, title, slug, excerpt, content, cover_image, meta_title, meta_description, meta_keywords, read_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING `+postColumns, post.CreatedBy, post.CategoryID, post.Title, post.Slug, post.Excerpt, post.Content, post.CoverImage, post.MetaTitle, post.MetaDescription, post.MetaKeywords, post.ReadTime), post)
	return post, err
}

func (s *store) GetAll(ctx context.Context) ([]*Post, error) {
	rows, err := s.db.Query(ctx, `SELECT `+postColumns+` FROM blog_posts ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*Post
	for rows.Next() {
		post := &Post{}
		if err := scanPost(rows, post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

func (s *store) GetByID(ctx context.Context, id string) (*Post, error) {
	post := &Post{}
	err := scanPost(s.db.QueryRow(ctx, `SELECT `+postColumns+` FROM blog_posts WHERE id = $1`, id), post)
	return post, err
}
func (s *store) GetBySlug(ctx context.Context, slug string) (*Post, error) {
	post := &Post{}
	err := scanPost(s.db.QueryRow(ctx, `SELECT `+postColumns+` FROM blog_posts WHERE slug = $1`, slug), post)
	return post, err
}

func (s *store) Update(ctx context.Context, post *Post) (*Post, error) {
	err := scanPost(s.db.QueryRow(ctx, `UPDATE blog_posts SET title = $1, slug = $2, excerpt = $3, content = $4, cover_image = $5, meta_title = $6, meta_description = $7, meta_keywords = $8, read_time = $9, updated_at = NOW() WHERE id = $10 RETURNING `+postColumns, post.Title, post.Slug, post.Excerpt, post.Content, post.CoverImage, post.MetaTitle, post.MetaDescription, post.MetaKeywords, post.ReadTime, post.ID), post)
	return post, err
}
func (s *store) Delete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM blog_posts WHERE id = $1`, id)
	return err
}
func (s *store) Publish(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `UPDATE blog_posts SET is_published = TRUE, published_at = NOW(), updated_at = NOW() WHERE id = $1`, id)
	return err
}
func (s *store) IncrementViews(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `UPDATE blog_posts SET views = views + 1 WHERE id = $1`, id)
	return err
}
