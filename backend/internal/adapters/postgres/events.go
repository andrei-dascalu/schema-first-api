package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"spec-first-backend/internal/app"
	"spec-first-backend/internal/domain"
	"spec-first-backend/internal/store/postgres/db"
)

var _ app.EventRepository = (*Store)(nil)

func (s *Store) CreateEvent(ctx context.Context, name string, date time.Time, location *string) (domain.Event, error) {
	e, err := s.q.CreateEvent(ctx, db.CreateEventParams{
		Name:     name,
		Date:     pgtype.Timestamptz{Time: date, Valid: true},
		Location: location,
	})
	if err != nil {
		return domain.Event{}, err
	}
	return toDomainEvent(e), nil
}

func (s *Store) ListEvents(ctx context.Context) ([]domain.Event, error) {
	rows, err := s.q.ListEvents(ctx)
	if err != nil {
		return nil, err
	}
	events := make([]domain.Event, 0, len(rows))
	for _, e := range rows {
		events = append(events, toDomainEvent(e))
	}
	return events, nil
}

func (s *Store) GetEvent(ctx context.Context, id uuid.UUID) (domain.Event, error) {
	e, err := s.q.GetEvent(ctx, toPgUUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Event{}, domain.ErrNotFound
		}
		return domain.Event{}, err
	}
	return toDomainEvent(e), nil
}

func (s *Store) UpdateEvent(ctx context.Context, id uuid.UUID, name string, date time.Time, location *string) (domain.Event, error) {
	e, err := s.q.UpdateEvent(ctx, db.UpdateEventParams{
		ID:       toPgUUID(id),
		Name:     name,
		Date:     pgtype.Timestamptz{Time: date, Valid: true},
		Location: location,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Event{}, domain.ErrNotFound
		}
		return domain.Event{}, err
	}
	return toDomainEvent(e), nil
}

func (s *Store) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	n, err := s.q.DeleteEvent(ctx, toPgUUID(id))
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
