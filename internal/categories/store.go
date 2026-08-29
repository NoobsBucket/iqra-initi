package categories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type Store interface {
	Create(ctx context.Context, name, slug, description, imageURL string) (*Category, error)
	GetAll(ctx context.Context) ([]*Category, error)
	GetByID(ctx context.Context, id string) (*Category, error)
	Update(ctx context.Context, id, name, slug, description, imageURL string) (*Category, error)
	Delete(ctx context.Context, id string) error
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{db: db}
}

func (s *store) Create(ctx context.Context, name, slug, description, imageURL string) (*Category, error) {
	category := &Category{}
	err := s.db.QueryRow(ctx, `
		INSERT INTO categories (name, slug, description, image_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, slug, description, image_url, created_at
	`, name, slug, description, imageURL).Scan(
		&category.ID, &category.Name, &category.Slug, &category.Description,
		&category.ImageURL, &category.CreatedAt,
	)
	return category, err
}

func (s *store) GetAll(ctx context.Context) ([]*Category, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, name, slug, description, image_url, created_at
		FROM categories ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		category := &Category{}
		if err := rows.Scan(
			&category.ID, &category.Name, &category.Slug, &category.Description,
			&category.ImageURL, &category.CreatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (s *store) GetByID(ctx context.Context, id string) (*Category, error) {
	category := &Category{}
	err := s.db.QueryRow(ctx, `
		SELECT id, name, slug, description, image_url, created_at
		FROM categories WHERE id = $1
	`, id).Scan(
		&category.ID, &category.Name, &category.Slug, &category.Description,
		&category.ImageURL, &category.CreatedAt,
	)
	return category, err
}

func (s *store) Update(ctx context.Context, id, name, slug, description, imageURL string) (*Category, error) {
	category := &Category{}
	err := s.db.QueryRow(ctx, `
		UPDATE categories
		SET name = $1, slug = $2, description = $3, image_url = $4
		WHERE id = $5
		RETURNING id, name, slug, description, image_url, created_at
	`, name, slug, description, imageURL, id).Scan(
		&category.ID, &category.Name, &category.Slug, &category.Description,
		&category.ImageURL, &category.CreatedAt,
	)
	return category, err
}

func (s *store) Delete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}
