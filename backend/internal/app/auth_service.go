package app

import (
	"context"
	"time"

	"spec-first-backend/internal/domain"
)

type AuthService struct {
	users  UserRepository
	tokens TokenIssuer
	hasher PasswordHasher
}

func NewAuthService(users UserRepository, tokens TokenIssuer, hasher PasswordHasher) *AuthService {
	return &AuthService{users: users, tokens: tokens, hasher: hasher}
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (domain.User, error) {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return domain.User{}, err
	}
	return s.users.CreateUser(ctx, email, hash, name)
}

// Login collapses "no such user" and "wrong password" into a single
// ErrInvalidCredentials, so the HTTP layer can't be used to enumerate
// registered emails.
func (s *AuthService) Login(ctx context.Context, email, password string) (token string, expiresAt time.Time, err error) {
	user, hash, err := s.users.UserByEmail(ctx, email)
	if err != nil {
		return "", time.Time{}, domain.ErrInvalidCredentials
	}

	if !s.hasher.Compare(hash, password) {
		return "", time.Time{}, domain.ErrInvalidCredentials
	}

	return s.tokens.Issue(user.ID)
}
