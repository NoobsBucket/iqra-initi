package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	IsVerified   bool      `json:"is_verified"`
	AvatarURL    string    `json:"avatar_url"`
	CreatedAt    time.Time `json:"created_at"`
}

func (u *User) SetAvatar() {
	u.AvatarURL = fmt.Sprintf(
		"http://localhost:4000/v1/users/%s/avatar",
		u.ID,
	)
}

// Store interface
type Store interface {
	CreateUser(ctx context.Context, name, email, passwordHash string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	SaveOTP(ctx context.Context, userID, otp string, expiresAt time.Time) error
	VerifyOTP(ctx context.Context, email, otp string) (*User, error)
	MarkVerified(ctx context.Context, userID string) error
	SaveResetOTP(ctx context.Context, userID, otp string, expiresAt time.Time) error
	VerifyResetOTP(ctx context.Context, email, otp string) (*User, error)
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
}

type store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &store{db: db}
}

func (s *store) CreateUser(ctx context.Context, name, email, passwordHash string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(ctx, `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, role, is_verified, created_at
	`, name, email, passwordHash).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.IsVerified,
		&user.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrEmailTaken
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}
	user.SetAvatar()
	return user, nil
}

func (s *store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(ctx, `
		SELECT id, name, email, password_hash, role, is_verified, created_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsVerified,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	user.SetAvatar()
	return user, nil
}

func (s *store) SaveOTP(ctx context.Context, userID, otp string, expiresAt time.Time) error {
	_, err := s.db.Exec(ctx, `
		UPDATE users 
		SET otp_code = $1, otp_expires_at = $2
		WHERE id = $3
	`, otp, expiresAt, userID)
	return err
}

func (s *store) VerifyOTP(ctx context.Context, email, otp string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(ctx, `
		SELECT id, name, email, role, created_at
		FROM users
		WHERE email = $1
		  AND otp_code = $2
		  AND otp_expires_at > NOW()
		  AND is_verified = FALSE
	`, email, otp).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	user.SetAvatar()
	return user, nil
}

func (s *store) MarkVerified(ctx context.Context, userID string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE users 
		SET is_verified = TRUE, otp_code = NULL, otp_expires_at = NULL
		WHERE id = $1
	`, userID)
	return err
}

func (s *store) SaveResetOTP(ctx context.Context, userID, otp string, expiresAt time.Time) error {
	_, err := s.db.Exec(ctx, `
		UPDATE users
		SET reset_otp = $1, reset_otp_expires_at = $2
		WHERE id = $3
	`, otp, expiresAt, userID)
	return err
}

func (s *store) VerifyResetOTP(ctx context.Context, email, otp string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(ctx, `
		SELECT id, name, email, role, created_at
		FROM users
		WHERE email = $1
		  AND reset_otp = $2
		  AND reset_otp_expires_at > NOW()
	`, email, otp).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	user.SetAvatar()
	return user, nil
}

func (s *store) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE users
		SET password_hash = $1, reset_otp = NULL, reset_otp_expires_at = NULL
		WHERE id = $2
	`, passwordHash, userID)
	return err
}

func (s *store) GetUserByID(ctx context.Context, id string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(ctx, `
		SELECT id, name, email, role, is_verified, created_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.IsVerified,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	user.SetAvatar()
	return user, nil
}
