package app

import (
	"context"
	"time"

	"github.com/google/uuid"

	"spec-first-backend/internal/domain"
)

type EventService struct {
	events EventRepository
}

func NewEventService(events EventRepository) *EventService {
	return &EventService{events: events}
}

func (s *EventService) ListEvents(ctx context.Context) ([]domain.Event, error) {
	return s.events.ListEvents(ctx)
}

func (s *EventService) CreateEvent(ctx context.Context, name string, date time.Time, location *string) (domain.Event, error) {
	return s.events.CreateEvent(ctx, name, date, location)
}

func (s *EventService) GetEvent(ctx context.Context, id uuid.UUID) (domain.Event, error) {
	return s.events.GetEvent(ctx, id)
}

func (s *EventService) UpdateEvent(ctx context.Context, id uuid.UUID, name string, date time.Time, location *string) (domain.Event, error) {
	return s.events.UpdateEvent(ctx, id, name, date, location)
}

func (s *EventService) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return s.events.DeleteEvent(ctx, id)
}
