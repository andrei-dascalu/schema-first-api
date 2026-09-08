package handlers

import (
	openapi_types "github.com/oapi-codegen/runtime/types"

	"spec-first-backend/internal/api"
	"spec-first-backend/internal/domain"
)

func toAPIUser(u domain.User) api.User {
	return api.User{
		Id:        u.ID,
		Email:     openapi_types.Email(u.Email),
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}

func toAPIEvent(e domain.Event) api.Event {
	return api.Event{
		Id:        e.ID,
		Name:      e.Name,
		Date:      e.Date,
		Location:  e.Location,
		CreatedAt: e.CreatedAt,
	}
}

func toAPIEvents(events []domain.Event) []api.Event {
	out := make([]api.Event, 0, len(events))
	for _, e := range events {
		out = append(out, toAPIEvent(e))
	}
	return out
}

func toAPIGuest(g domain.Guest) api.Guest {
	return api.Guest{
		Id:         g.ID,
		EventId:    g.EventID,
		Name:       g.Name,
		Email:      toAPIEmail(g.Email),
		RsvpStatus: api.RsvpStatus(g.RsvpStatus),
		CreatedAt:  g.CreatedAt,
	}
}

func toAPIGuests(guests []domain.Guest) []api.Guest {
	out := make([]api.Guest, 0, len(guests))
	for _, g := range guests {
		out = append(out, toAPIGuest(g))
	}
	return out
}

func toDomainRsvpPtr(status *api.RsvpStatus) *domain.RsvpStatus {
	if status == nil {
		return nil
	}
	v := domain.RsvpStatus(*status)
	return &v
}

func toAPIEmail(email *string) *openapi_types.Email {
	if email == nil {
		return nil
	}
	v := openapi_types.Email(*email)
	return &v
}

func toDomainEmail(email *openapi_types.Email) *string {
	if email == nil {
		return nil
	}
	v := string(*email)
	return &v
}
