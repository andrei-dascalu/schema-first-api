package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"spec-first-backend/internal/app"
	"spec-first-backend/internal/domain"
	"spec-first-backend/internal/store/postgres/db"
)

var _ app.UserRepository = (*Store)(nil)

func (s *Store) CreateUser(ctx context.Context, email, passwordHash, name string) (domain.User, error) {
	u, err := s.q.CreateUser(ctx, db.CreateUserParams{Email: email, PasswordHash: passwordHash, Name: name})
	if err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, domain.ErrDuplicateEmail
		}
		return domain.User{}, err
	}
	return toDomainUser(u), nil
}

// UserByEmail returns the stored user and password hash for login verification.
func (s *Store) UserByEmail(ctx context.Context, email string) (domain.User, string, error) {
	u, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, "", domain.ErrNotFound
		}
		return domain.User{}, "", err
	}
	return toDomainUser(u), u.PasswordHash, nil
}
