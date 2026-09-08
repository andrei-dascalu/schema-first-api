package postgres

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"spec-first-backend/internal/domain"
	"spec-first-backend/internal/store/postgres/db"
)

func toPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func fromPgUUID(id pgtype.UUID) uuid.UUID {
	return uuid.UUID(id.Bytes)
}

func toDomainEvent(e db.Event) domain.Event {
	return domain.Event{
		ID:        fromPgUUID(e.ID),
		Name:      e.Name,
		Date:      e.Date.Time,
		Location:  e.Location,
		CreatedAt: e.CreatedAt.Time,
	}
}

func toDomainGuest(g db.Guest) domain.Guest {
	return domain.Guest{
		ID:         fromPgUUID(g.ID),
		EventID:    fromPgUUID(g.EventID),
		Name:       g.Name,
		Email:      g.Email,
		RsvpStatus: domain.RsvpStatus(g.RsvpStatus),
		CreatedAt:  g.CreatedAt.Time,
	}
}

func toDomainUser(u db.User) domain.User {
	return domain.User{
		ID:        fromPgUUID(u.ID),
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt.Time,
	}
}
