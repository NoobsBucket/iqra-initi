package reviews

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Review struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CourseID   string    `json:"course_id"`
	ReviewText string    `json:"review_text"`
	Rating     float64   `json:"rating"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Store interface {
	Create(ctx context.Context, userID, courseID, reviewText string, rating float64) (*Review, error)
	GetByCourse(ctx context.Context, courseID string) ([]*Review, error)
	GetByID(ctx context.Context, id string) (*Review, error)
	Update(ctx context.Context, id, reviewText string, rating float64) (*Review, error)
	Delete(ctx context.Context, id string) error
}

type store struct{ db *pgxpool.Pool }

func NewStore(db *pgxpool.Pool) Store { return &store{db: db} }

const reviewColumns = `id, user_id, course_id, review_text, rating, created_at, updated_at`

func scanReview(row interface{ Scan(...any) error }, review *Review) error {
	return row.Scan(&review.ID, &review.UserID, &review.CourseID, &review.ReviewText, &review.Rating, &review.CreatedAt, &review.UpdatedAt)
}

func (s *store) Create(ctx context.Context, userID, courseID, reviewText string, rating float64) (*Review, error) {
	review := &Review{}
	err := scanReview(s.db.QueryRow(ctx, `INSERT INTO reviews (user_id, course_id, review_text, rating) VALUES ($1, $2, $3, $4) RETURNING `+reviewColumns, userID, courseID, reviewText, rating), review)
	return review, err
}

func (s *store) GetByCourse(ctx context.Context, courseID string) ([]*Review, error) {
	rows, err := s.db.Query(ctx, `SELECT `+reviewColumns+` FROM reviews WHERE course_id = $1 ORDER BY created_at DESC`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []*Review
	for rows.Next() {
		review := &Review{}
		if err := scanReview(rows, review); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, rows.Err()
}

func (s *store) GetByID(ctx context.Context, id string) (*Review, error) {
	review := &Review{}
	err := scanReview(s.db.QueryRow(ctx, `SELECT `+reviewColumns+` FROM reviews WHERE id = $1`, id), review)
	return review, err
}

func (s *store) Update(ctx context.Context, id, reviewText string, rating float64) (*Review, error) {
	review := &Review{}
	err := scanReview(s.db.QueryRow(ctx, `UPDATE reviews SET review_text = $1, rating = $2, updated_at = NOW() WHERE id = $3 RETURNING `+reviewColumns, reviewText, rating, id), review)
	return review, err
}

func (s *store) Delete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM reviews WHERE id = $1`, id)
	return err
}
