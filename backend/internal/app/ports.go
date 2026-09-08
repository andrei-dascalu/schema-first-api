// Package app orchestrates domain logic behind ports implemented by
// internal/adapters/*. It imports only internal/domain and its own ports -
// never a concrete adapter package.
package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"spec-first-backend/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash, name string) (domain.User, error)
	UserByEmail(ctx context.Context, email string) (domain.User, string, error)
}

type EventRepository interface {
	CreateEvent(ctx context.Context, name string, date time.Time, location *string) (domain.Event, error)
	ListEvents(ctx context.Context) ([]domain.Event, error)
	GetEvent(ctx context.Context, id uuid.UUID) (domain.Event, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, name string, date time.Time, location *string) (domain.Event, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
}

type GuestRepository interface {
	CreateGuest(ctx context.Context, eventID uuid.UUID, name string, email *string, rsvp domain.RsvpStatus) (domain.Guest, error)
	ListGuestsByEvent(ctx context.Context, eventID uuid.UUID) ([]domain.Guest, error)
	GetGuest(ctx context.Context, eventID, guestID uuid.UUID) (domain.Guest, error)
	UpdateGuest(ctx context.Context, eventID, guestID uuid.UUID, name string, email *string, rsvp domain.RsvpStatus) (domain.Guest, error)
	DeleteGuest(ctx context.Context, eventID, guestID uuid.UUID) error
}

type TokenIssuer interface {
	Issue(userID uuid.UUID) (token string, expiresAt time.Time, err error)
	Parse(token string) (uuid.UUID, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}
