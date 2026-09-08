package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spec-first-backend/internal/app"
	"spec-first-backend/internal/domain"
)

type fakeUserRepository struct {
	usersByEmail map[string]domain.User
	hashesByID   map[uuid.UUID]string
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		usersByEmail: map[string]domain.User{},
		hashesByID:   map[uuid.UUID]string{},
	}
}

func (f *fakeUserRepository) CreateUser(_ context.Context, email, passwordHash, name string) (domain.User, error) {
	if _, ok := f.usersByEmail[email]; ok {
		return domain.User{}, domain.ErrDuplicateEmail
	}
	user := domain.User{ID: uuid.New(), Email: email, Name: name, CreatedAt: time.Now()}
	f.usersByEmail[email] = user
	f.hashesByID[user.ID] = passwordHash
	return user, nil
}

func (f *fakeUserRepository) UserByEmail(_ context.Context, email string) (domain.User, string, error) {
	user, ok := f.usersByEmail[email]
	if !ok {
		return domain.User{}, "", domain.ErrNotFound
	}
	return user, f.hashesByID[user.ID], nil
}

type fakeTokenIssuer struct{}

func (fakeTokenIssuer) Issue(userID uuid.UUID) (string, time.Time, error) {
	return "token-for-" + userID.String(), time.Now().Add(time.Hour), nil
}

func (fakeTokenIssuer) Parse(_ string) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}

type fakeHasher struct{}

func (fakeHasher) Hash(password string) (string, error) {
	return "hashed:" + password, nil
}

func (fakeHasher) Compare(hash, password string) bool {
	return hash == "hashed:"+password
}

func TestAuthServiceRegisterThenLogin(t *testing.T) {
	users := newFakeUserRepository()
	svc := app.NewAuthService(users, fakeTokenIssuer{}, fakeHasher{})
	ctx := context.Background()

	user, err := svc.Register(ctx, "alice@example.com", "hunter22", "Alice")
	require.NoError(t, err)
	require.Equal(t, "alice@example.com", user.Email)

	token, expiresAt, err := svc.Login(ctx, "alice@example.com", "hunter22")
	require.NoError(t, err)
	assert.Equal(t, "token-for-"+user.ID.String(), token)
	assert.WithinDuration(t, time.Now().Add(time.Hour), expiresAt, time.Second)
}

func TestAuthServiceLoginRejectsWrongPassword(t *testing.T) {
	users := newFakeUserRepository()
	svc := app.NewAuthService(users, fakeTokenIssuer{}, fakeHasher{})
	ctx := context.Background()

	_, err := svc.Register(ctx, "alice@example.com", "hunter22", "Alice")
	require.NoError(t, err)

	_, _, err = svc.Login(ctx, "alice@example.com", "wrong-password")
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestAuthServiceLoginRejectsUnknownEmail(t *testing.T) {
	svc := app.NewAuthService(newFakeUserRepository(), fakeTokenIssuer{}, fakeHasher{})

	_, _, err := svc.Login(context.Background(), "nobody@example.com", "hunter22")
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}
