package blog

import (
	"context"
	"time"
)

type BlogCategory struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *store) CreateCategory(ctx context.Context, name, slug, description string) (*BlogCategory, error) {
	c := &BlogCategory{}
	err := s.db.QueryRow(ctx, `
		INSERT INTO blog_categories (name, slug, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, slug, description, created_at
	`, name, slug, description).Scan(
		&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt,
	)
	return c, err
}

func (s *store) GetAllCategories(ctx context.Context) ([]*BlogCategory, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, name, slug, description, created_at
		FROM blog_categories ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*BlogCategory
	for rows.Next() {
		c := &BlogCategory{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func (s *store) GetCategoryByID(ctx context.Context, id string) (*BlogCategory, error) {
	c := &BlogCategory{}
	err := s.db.QueryRow(ctx, `
		SELECT id, name, slug, description, created_at
		FROM blog_categories WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt)
	return c, err
}

func (s *store) UpdateCategory(ctx context.Context, id, name, slug, description string) (*BlogCategory, error) {
	c := &BlogCategory{}
	err := s.db.QueryRow(ctx, `
		UPDATE blog_categories SET name = $1, slug = $2, description = $3
		WHERE id = $4
		RETURNING id, name, slug, description, created_at
	`, name, slug, description, id).Scan(
		&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt,
	)
	return c, err
}

func (s *store) DeleteCategory(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM blog_categories WHERE id = $1`, id)
	return err
}
