package notifications

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type Store interface {
	Create(ctx context.Context, userID, title, message, notifType string) (*Notification, error)
	GetByUser(ctx context.Context, userID string) ([]*Notification, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{db: db}
}

func (s *store) Create(ctx context.Context, userID, title, message, notifType string) (*Notification, error) {
	n := &Notification{}
	err := s.db.QueryRow(ctx, `
		INSERT INTO notifications (user_id, title, message, type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, message, type, is_read, created_at
	`, userID, title, message, notifType).Scan(
		&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt,
	)
	return n, err
}

func (s *store) GetByUser(ctx context.Context, userID string) ([]*Notification, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, user_id, title, message, type, is_read, created_at
		FROM notifications WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		n := &Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

func (s *store) MarkRead(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `UPDATE notifications SET is_read = TRUE WHERE id = $1`, id)
	return err
}

func (s *store) MarkAllRead(ctx context.Context, userID string) error {
	_, err := s.db.Exec(ctx, `UPDATE notifications SET is_read = TRUE WHERE user_id = $1`, userID)
	return err
}

func (s *store) Delete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM notifications WHERE id = $1`, id)
	return err
}
