package courses

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Course struct {
	ID            string    `json:"id"`
	CreatedBy     string    `json:"created_by"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ImageURL      string    `json:"image_url"`
	Price         float64   `json:"price"`
	DiscountPrice float64   `json:"discount_price"`
	Currency      string    `json:"currency"`
	IsPublished   bool      `json:"is_published"`
	IsFeatured    bool      `json:"is_featured"`
	AverageRating float64   `json:"average_rating"`
	TotalReviews  int       `json:"total_reviews"`
	TotalStudents int       `json:"total_students"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CategoryIDs   []string  `json:"category_ids,omitempty"`
}

type Store interface {
	Create(ctx context.Context, course *Course) (*Course, error)
	GetAll(ctx context.Context) ([]*Course, error)
	GetByID(ctx context.Context, id string) (*Course, error)
	Update(ctx context.Context, course *Course) (*Course, error)
	Delete(ctx context.Context, id string) error
	Publish(ctx context.Context, id string) error
	AddCategory(ctx context.Context, courseID, categoryID string) error
	ClearCategories(ctx context.Context, courseID string) error
	GetCategoryIDs(ctx context.Context, courseID string) ([]string, error)
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{db: db}
}

func (s *store) Create(ctx context.Context, c *Course) (*Course, error) {
	err := s.db.QueryRow(ctx, `
		INSERT INTO courses (created_by, title, description, image_url, price, discount_price, currency)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_by, title, description, image_url, price, discount_price, currency, is_published, is_featured, average_rating, total_reviews, total_students, created_at, updated_at
	`, c.CreatedBy, c.Title, c.Description, c.ImageURL, c.Price, c.DiscountPrice, c.Currency).Scan(
		&c.ID, &c.CreatedBy, &c.Title, &c.Description, &c.ImageURL,
		&c.Price, &c.DiscountPrice, &c.Currency, &c.IsPublished, &c.IsFeatured,
		&c.AverageRating, &c.TotalReviews, &c.TotalStudents, &c.CreatedAt, &c.UpdatedAt,
	)
	return c, err
}

func (s *store) GetAll(ctx context.Context) ([]*Course, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, created_by, title, description, image_url, price, discount_price, currency, is_published, is_featured, average_rating, total_reviews, total_students, created_at, updated_at
		FROM courses ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*Course
	for rows.Next() {
		c := &Course{}
		if err := rows.Scan(
			&c.ID, &c.CreatedBy, &c.Title, &c.Description, &c.ImageURL,
			&c.Price, &c.DiscountPrice, &c.Currency, &c.IsPublished, &c.IsFeatured,
			&c.AverageRating, &c.TotalReviews, &c.TotalStudents, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, nil
}

func (s *store) GetByID(ctx context.Context, id string) (*Course, error) {
	c := &Course{}
	err := s.db.QueryRow(ctx, `
		SELECT id, created_by, title, description, image_url, price, discount_price, currency, is_published, is_featured, average_rating, total_reviews, total_students, created_at, updated_at
		FROM courses WHERE id = $1
	`, id).Scan(
		&c.ID, &c.CreatedBy, &c.Title, &c.Description, &c.ImageURL,
		&c.Price, &c.DiscountPrice, &c.Currency, &c.IsPublished, &c.IsFeatured,
		&c.AverageRating, &c.TotalReviews, &c.TotalStudents, &c.CreatedAt, &c.UpdatedAt,
	)
	return c, err
}

func (s *store) Update(ctx context.Context, c *Course) (*Course, error) {
	err := s.db.QueryRow(ctx, `
		UPDATE courses
		SET title = $1, description = $2, image_url = $3, price = $4, discount_price = $5, currency = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING id, created_by, title, description, image_url, price, discount_price, currency, is_published, is_featured, average_rating, total_reviews, total_students, created_at, updated_at
	`, c.Title, c.Description, c.ImageURL, c.Price, c.DiscountPrice, c.Currency, c.ID).Scan(
		&c.ID, &c.CreatedBy, &c.Title, &c.Description, &c.ImageURL,
		&c.Price, &c.DiscountPrice, &c.Currency, &c.IsPublished, &c.IsFeatured,
		&c.AverageRating, &c.TotalReviews, &c.TotalStudents, &c.CreatedAt, &c.UpdatedAt,
	)
	return c, err
}

func (s *store) Delete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM courses WHERE id = $1`, id)
	return err
}

func (s *store) Publish(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `UPDATE courses SET is_published = TRUE, updated_at = NOW() WHERE id = $1`, id)
	return err
}

func (s *store) AddCategory(ctx context.Context, courseID, categoryID string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO course_categories (course_id, category_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, courseID, categoryID)
	return err
}

func (s *store) ClearCategories(ctx context.Context, courseID string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM course_categories WHERE course_id = $1`, courseID)
	return err
}

func (s *store) GetCategoryIDs(ctx context.Context, courseID string) ([]string, error) {
	rows, err := s.db.Query(ctx, `SELECT category_id FROM course_categories WHERE course_id = $1`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}