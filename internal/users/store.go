package users

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
	AvatarURL  string    `json:"avatar_url"`
	CreatedAt  time.Time `json:"created_at"`
}

type Store interface {
	GetAll(ctx context.Context, search string) ([]*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	UpdateRole(ctx context.Context, id, role string) (*User, error)
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{db: db}
}

func (s *store) GetAll(ctx context.Context, search string) ([]*User, error) {
	query := `
		SELECT id, name, email, role, is_verified, created_at
		FROM users
		WHERE ($1 = '' OR name ILIKE '%' || $1 || '%' OR email ILIKE '%' || $1 || '%')
		ORDER BY created_at DESC
	`
	rows, err := s.db.Query(ctx, query, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.Role, &u.IsVerified, &u.CreatedAt,
		); err != nil {
			return nil, err
		}
		u.AvatarURL = "http://localhost:4000/v1/users/" + u.ID + "/avatar"
		users = append(users, u)
	}
	return users, rows.Err()
}

func (s *store) GetByID(ctx context.Context, id string) (*User, error) {
	u := &User{}
	err := s.db.QueryRow(ctx, `
		SELECT id, name, email, role, is_verified, created_at
		FROM users WHERE id = $1
	`, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.Role, &u.IsVerified, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.AvatarURL = "http://localhost:4000/v1/users/" + u.ID + "/avatar"
	return u, nil
}

func (s *store) UpdateRole(ctx context.Context, id, role string) (*User, error) {
	u := &User{}
	err := s.db.QueryRow(ctx, `
		UPDATE users SET role = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, name, email, role, is_verified, created_at
	`, role, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.Role, &u.IsVerified, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.AvatarURL = "http://localhost:4000/v1/users/" + u.ID + "/avatar"
	return u, nil
}
