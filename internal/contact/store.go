package contact

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Message struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type Store interface {
	Create(ctx context.Context, m *Message) (*Message, error)
	GetAll(ctx context.Context) ([]*Message, error)
	MarkRead(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{db: db}
}

func (s *store) Create(ctx context.Context, m *Message) (*Message, error) {
	err := s.db.QueryRow(ctx, `
		INSERT INTO contact_messages (name, email, phone, subject, message)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, email, phone, subject, message, is_read, created_at
	`, m.Name, m.Email, m.Phone, m.Subject, m.Message).Scan(
		&m.ID, &m.Name, &m.Email, &m.Phone, &m.Subject, &m.Message, &m.IsRead, &m.CreatedAt,
	)
	return m, err
}

func (s *store) GetAll(ctx context.Context) ([]*Message, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, name, email, phone, subject, message, is_read, created_at
		FROM contact_messages ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		m := &Message{}
		if err := rows.Scan(
			&m.ID, &m.Name, &m.Email, &m.Phone, &m.Subject, &m.Message, &m.IsRead, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (s *store) MarkRead(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `UPDATE contact_messages SET is_read = TRUE WHERE id = $1`, id)
	return err
}

func (s *store) Delete(ctx context.Context, id string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM contact_messages WHERE id = $1`, id)
	return err
}
