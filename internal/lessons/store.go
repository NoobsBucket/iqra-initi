package lessons

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Lesson struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	VideoURL    string    `json:"video_url"`
	OrderIndex  int       `json:"order_index"`
	IsFree      bool      `json:"is_free"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Store interface {
	Create(ctx context.Context, lesson *Lesson) (*Lesson, error)
	GetByCourse(ctx context.Context, courseID string) ([]*Lesson, error)
	GetByID(ctx context.Context, id string) (*Lesson, error)
	Update(ctx context.Context, lesson *Lesson) (*Lesson, error)
	Delete(ctx context.Context, id string) error
	Reorder(ctx context.Context, id string, order int) error
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{db: db}
}

func (s *store) Create(ctx context.Context, l *Lesson) (*Lesson, error) {
	err := s.db.QueryRow(ctx, `
		INSERT INTO lessons (course_id, title, description, video_url, order_index, is_free)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, course_id, title, description, video_url, order_index, is_free, created_at, updated_at
	`, l.CourseID, l.Title, l.Description, l.VideoURL, l.OrderIndex, l.IsFree).Scan(
		&l.ID, &l.CourseID, &l.Title, &l.Description, &l.VideoURL,
		&l.OrderIndex, &l.IsFree, &l.CreatedAt, &l.UpdatedAt,
	)
	return l, err
}

func (s *store) GetByCourse(ctx context.Context, courseID string) ([]*Lesson, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, course_id, title, description, video_url, order_index, is_free, created_at, updated_at
		FROM lessons WHERE course_id = $1 ORDER BY order_index ASC
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []*Lesson
	for rows.Next() {
		l := &Lesson{}
		if err := rows.Scan(
			&l.ID, &l.CourseID, &l.Title, &l.Description, &l.VideoURL,
			&l.OrderIndex, &l.IsFree, &l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}
	return lessons, rows.Err()
}

func (s *store) GetByID(ctx context.Context, id string) (*Lesson, error) {
	l := &Lesson{}
	err := s.db.QueryRow(ctx, `
		SELECT id, course_id, title, description, video_url, order_index, is_free, created_at, updated_at
		FROM lessons WHERE id = $1
	`, id).Scan(
		&l.ID, &l.CourseID, &l.Title, &l.Description, &l.VideoURL,
		&l.OrderIndex, &l.IsFree, &l.CreatedAt, &l.UpdatedAt,
	)
	return l, err
}

func (s *store) Update(ctx context.Context, l *Lesson) (*Lesson, error) {
	err := s.db.QueryRow(ctx, `
		UPDATE lessons
		SET title = $1, description = $2, video_url = $3, order_index = $4, is_free = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING id, course_id, title, description, video_url, order_index, is_free, created_at, updated_at
	`, l.Title, l.Description, l.VideoURL, l.OrderIndex, l.IsFree, l.ID).Scan(
		&l.ID, &l.CourseID, &l.Title, &l.Description, &l.VideoURL,
		&l.OrderIndex, &l.IsFree, &l.CreatedAt, &l.UpdatedAt,
	)
	return l, err
}

func (s *store) Delete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM lessons WHERE id = $1`, id)
	return err
}

func (s *store) Reorder(ctx context.Context, id string, order int) error {
	_, err := s.db.Exec(ctx, `UPDATE lessons SET order_index = $1, updated_at = NOW() WHERE id = $2`, order, id)
	return err
}
