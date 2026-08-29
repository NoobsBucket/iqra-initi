package enrollment

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Enrollment struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	CourseID      string     `json:"course_id"`
	FullName      string     `json:"full_name"`
	Whatsapp      string     `json:"whatsapp"`
	Notes         string     `json:"notes"`
	Status        *string    `json:"status"`
	Progress      int        `json:"progress"`
	PaymentStatus *string    `json:"payment_status"`
	PaymentMethod *string    `json:"payment_method"`
	EnrolledAt    time.Time  `json:"enrolled_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type Store interface {
	Enroll(ctx context.Context, e *Enrollment) (*Enrollment, error)
	GetByUser(ctx context.Context, userID string) ([]*Enrollment, error)
	GetByID(ctx context.Context, id string) (*Enrollment, error)
	GetByCourse(ctx context.Context, courseID string) ([]*Enrollment, error)
	UpdateProgress(ctx context.Context, id string, progress int) error
	UpdatePaymentStatus(ctx context.Context, id, status string) error
	Complete(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	IsEnrolled(ctx context.Context, userID, courseID string) (bool, error)
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{db: db}
}

func (s *store) Enroll(ctx context.Context, e *Enrollment) (*Enrollment, error) {
	err := s.db.QueryRow(ctx, `
		INSERT INTO user_courses (user_id, course_id, full_name, whatsapp, notes, payment_method)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, course_id, full_name, whatsapp, notes, status, progress, payment_status, payment_method, enrolled_at, completed_at, updated_at
	`, e.UserID, e.CourseID, e.FullName, e.Whatsapp, e.Notes, e.PaymentMethod).Scan(
		&e.ID, &e.UserID, &e.CourseID, &e.FullName, &e.Whatsapp, &e.Notes,
		&e.Status, &e.Progress, &e.PaymentStatus, &e.PaymentMethod, &e.EnrolledAt, &e.CompletedAt, &e.UpdatedAt,
	)
	return e, err
}

func (s *store) GetByUser(ctx context.Context, userID string) ([]*Enrollment, error) {
	rows, err := s.db.Query(ctx, `
		SELECT uc.id, uc.user_id, uc.course_id, uc.full_name, uc.whatsapp, uc.notes,
		       uc.status, uc.progress, uc.payment_status, uc.payment_method, uc.enrolled_at, uc.completed_at, uc.updated_at
		FROM user_courses uc
		WHERE uc.user_id = $1
		ORDER BY uc.enrolled_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []*Enrollment
	for rows.Next() {
		e := &Enrollment{}
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.CourseID, &e.FullName, &e.Whatsapp, &e.Notes,
			&e.Status, &e.Progress, &e.PaymentStatus, &e.PaymentMethod, &e.EnrolledAt, &e.CompletedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, e)
	}
	return enrollments, rows.Err()
}

func (s *store) GetByID(ctx context.Context, id string) (*Enrollment, error) {
	e := &Enrollment{}
	err := s.db.QueryRow(ctx, `
		SELECT id, user_id, course_id, full_name, whatsapp, notes, status, progress, payment_status, payment_method, enrolled_at, completed_at, updated_at
		FROM user_courses WHERE id = $1
	`, id).Scan(
		&e.ID, &e.UserID, &e.CourseID, &e.FullName, &e.Whatsapp, &e.Notes,
		&e.Status, &e.Progress, &e.PaymentStatus, &e.PaymentMethod, &e.EnrolledAt, &e.CompletedAt, &e.UpdatedAt,
	)
	return e, err
}

func (s *store) GetByCourse(ctx context.Context, courseID string) ([]*Enrollment, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, user_id, course_id, full_name, whatsapp, notes, status, progress, payment_status, payment_method, enrolled_at, completed_at, updated_at
		FROM user_courses WHERE course_id = $1
		ORDER BY enrolled_at DESC
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []*Enrollment
	for rows.Next() {
		e := &Enrollment{}
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.CourseID, &e.FullName, &e.Whatsapp, &e.Notes,
			&e.Status, &e.Progress, &e.PaymentStatus, &e.PaymentMethod, &e.EnrolledAt, &e.CompletedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, e)
	}
	return enrollments, rows.Err()
}

func (s *store) UpdateProgress(ctx context.Context, id string, progress int) error {
	_, err := s.db.Exec(ctx, `
		UPDATE user_courses SET progress = $1, updated_at = NOW() WHERE id = $2
	`, progress, id)
	return err
}

func (s *store) UpdatePaymentStatus(ctx context.Context, id, status string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE user_courses 
		SET payment_status = $1, 
		    status = CASE WHEN $1 = 'paid' THEN 'active' ELSE status END,
		    updated_at = NOW()
		WHERE id = $2
	`, status, id)
	return err
}

func (s *store) Complete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE user_courses 
		SET status = 'completed', completed_at = NOW(), progress = 100, updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

func (s *store) Delete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM user_courses WHERE id = $1`, id)
	return err
}

func (s *store) IsEnrolled(ctx context.Context, userID, courseID string) (bool, error) {
	var count int
	err := s.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM user_courses 
		WHERE user_id = $1 AND course_id = $2
	`, userID, courseID).Scan(&count)
	return count > 0, err
}
