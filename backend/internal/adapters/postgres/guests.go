package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"spec-first-backend/internal/app"
	"spec-first-backend/internal/domain"
	"spec-first-backend/internal/store/postgres/db"
)

var _ app.GuestRepository = (*Store)(nil)

func (s *Store) CreateGuest(ctx context.Context, eventID uuid.UUID, name string, email *string, rsvp domain.RsvpStatus) (domain.Guest, error) {
	g, err := s.q.CreateGuest(ctx, db.CreateGuestParams{
		EventID:    toPgUUID(eventID),
		Name:       name,
		Email:      email,
		RsvpStatus: db.RsvpStatus(rsvp),
	})
	if err != nil {
		if isForeignKeyViolation(err) {
			return domain.Guest{}, domain.ErrNotFound
		}
		return domain.Guest{}, err
	}
	return toDomainGuest(g), nil
}

func (s *Store) ListGuestsByEvent(ctx context.Context, eventID uuid.UUID) ([]domain.Guest, error) {
	rows, err := s.q.ListGuestsByEvent(ctx, toPgUUID(eventID))
	if err != nil {
		return nil, err
	}
	guests := make([]domain.Guest, 0, len(rows))
	for _, g := range rows {
		guests = append(guests, toDomainGuest(g))
	}
	return guests, nil
}

func (s *Store) GetGuest(ctx context.Context, eventID, guestID uuid.UUID) (domain.Guest, error) {
	g, err := s.q.GetGuest(ctx, db.GetGuestParams{ID: toPgUUID(guestID), EventID: toPgUUID(eventID)})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Guest{}, domain.ErrNotFound
		}
		return domain.Guest{}, err
	}
	return toDomainGuest(g), nil
}

func (s *Store) UpdateGuest(ctx context.Context, eventID, guestID uuid.UUID, name string, email *string, rsvp domain.RsvpStatus) (domain.Guest, error) {
	g, err := s.q.UpdateGuest(ctx, db.UpdateGuestParams{
		ID:         toPgUUID(guestID),
		EventID:    toPgUUID(eventID),
		Name:       name,
		Email:      email,
		RsvpStatus: db.RsvpStatus(rsvp),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Guest{}, domain.ErrNotFound
		}
		return domain.Guest{}, err
	}
	return toDomainGuest(g), nil
}

func (s *Store) DeleteGuest(ctx context.Context, eventID, guestID uuid.UUID) error {
	n, err := s.q.DeleteGuest(ctx, db.DeleteGuestParams{ID: toPgUUID(guestID), EventID: toPgUUID(eventID)})
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
