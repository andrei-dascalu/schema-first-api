package app

import (
	"context"

	"github.com/google/uuid"

	"spec-first-backend/internal/domain"
)

type GuestService struct {
	guests GuestRepository
}

func NewGuestService(guests GuestRepository) *GuestService {
	return &GuestService{guests: guests}
}

func (s *GuestService) ListGuestsByEvent(ctx context.Context, eventID uuid.UUID) ([]domain.Guest, error) {
	return s.guests.ListGuestsByEvent(ctx, eventID)
}

// CreateGuest defaults rsvp to RsvpPending when the caller omits it.
func (s *GuestService) CreateGuest(ctx context.Context, eventID uuid.UUID, name string, email *string, rsvp *domain.RsvpStatus) (domain.Guest, error) {
	return s.guests.CreateGuest(ctx, eventID, name, email, rsvpOrDefault(rsvp))
}

func (s *GuestService) GetGuest(ctx context.Context, eventID, guestID uuid.UUID) (domain.Guest, error) {
	return s.guests.GetGuest(ctx, eventID, guestID)
}

// UpdateGuest defaults rsvp to RsvpPending when the caller omits it.
func (s *GuestService) UpdateGuest(ctx context.Context, eventID, guestID uuid.UUID, name string, email *string, rsvp *domain.RsvpStatus) (domain.Guest, error) {
	return s.guests.UpdateGuest(ctx, eventID, guestID, name, email, rsvpOrDefault(rsvp))
}

func (s *GuestService) DeleteGuest(ctx context.Context, eventID, guestID uuid.UUID) error {
	return s.guests.DeleteGuest(ctx, eventID, guestID)
}

func rsvpOrDefault(status *domain.RsvpStatus) domain.RsvpStatus {
	if status == nil {
		return domain.RsvpPending
	}
	return *status
}
