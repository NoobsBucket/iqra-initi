package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/NoobsBucket/iqra-initi/internal/mailer"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already in use")
	ErrInvalidOTP         = errors.New("invalid or expired OTP")
	ErrAlreadyVerified    = errors.New("email already verified")
	ErrNotVerified        = errors.New("email not verified")
)

type service interface {
	RegisterUser(ctx context.Context, name, email, password string) error
	LoginUser(ctx context.Context, email, password string) (*User, string, error)
	VerifyOTP(ctx context.Context, email, otp string) (*User, string, error)
	ResendOTP(ctx context.Context, email string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, otp, newPassword string) error
}

type Service struct {
	store     Store
	jwtSecret string
	jwtExp    time.Duration
	mailer    *mailer.Mailer
}

func NewService(store Store, jwtSecret string, jwtExp time.Duration, mailer *mailer.Mailer) *Service {
	return &Service{
		store:     store,
		jwtSecret: jwtSecret,
		jwtExp:    jwtExp,
		mailer:    mailer,
	}
}

func (s *Service) RegisterUser(ctx context.Context, name, email, password string) error {
	name = strings.TrimSpace(name)
	email = normalizeEmail(email)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := s.store.CreateUser(ctx, name, email, string(hash))
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			return ErrEmailTaken
		}
		return fmt.Errorf("create user: %w", err)
	}

	otp, err := mailer.GenerateOTP()
	if err != nil {
		return err
	}

	if err := s.store.SaveOTP(ctx, user.ID, otp, time.Now().Add(10*time.Minute)); err != nil {
		return err
	}

	return s.mailer.SendOTP(email, name, otp)
}

func (s *Service) VerifyOTP(ctx context.Context, email, otp string) (*User, string, error) {
	user, err := s.store.VerifyOTP(ctx, normalizeEmail(email), otp)
	if err != nil {
		return nil, "", ErrInvalidOTP
	}

	if err := s.store.MarkVerified(ctx, user.ID); err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (*User, string, error) {
	email = normalizeEmail(email)

	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, "", ErrNotVerified
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// generate token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *Service) ResendOTP(ctx context.Context, email string) error {
	email = normalizeEmail(email)
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil
	}

	if user.IsVerified {
		return ErrAlreadyVerified
	}

	otp, err := mailer.GenerateOTP()
	if err != nil {
		return err
	}

	if err := s.store.SaveOTP(ctx, user.ID, otp, time.Now().Add(10*time.Minute)); err != nil {
		return err
	}

	return s.mailer.SendResendOTP(email, user.Name, otp)
}

func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	email = normalizeEmail(email)
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil
	}

	otp, err := mailer.GenerateOTP()
	if err != nil {
		return err
	}

	if err := s.store.SaveResetOTP(ctx, user.ID, otp, time.Now().Add(10*time.Minute)); err != nil {
		return err
	}

	return s.mailer.SendResetOTP(email, user.Name, otp)
}

func (s *Service) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	user, err := s.store.VerifyResetOTP(ctx, normalizeEmail(email), otp)
	if err != nil {
		return ErrInvalidOTP
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.store.UpdatePassword(ctx, user.ID, string(hash))
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (s *Service) generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(s.jwtExp).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
