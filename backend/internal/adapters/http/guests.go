package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"spec-first-backend/internal/api"
	"spec-first-backend/internal/domain"
)

func (h *Handlers) ListGuests(w http.ResponseWriter, r *http.Request, eventId api.EventId) {
	guests, err := h.guests.ListGuestsByEvent(r.Context(), eventId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list guests")
		return
	}
	writeJSON(w, http.StatusOK, toAPIGuests(guests))
}

func (h *Handlers) CreateGuest(w http.ResponseWriter, r *http.Request, eventId api.EventId) {
	var body api.NewGuest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	guest, err := h.guests.CreateGuest(r.Context(), eventId, body.Name, toDomainEmail(body.Email), toDomainRsvpPtr(body.RsvpStatus))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create guest")
		return
	}
	writeJSON(w, http.StatusCreated, toAPIGuest(guest))
}

func (h *Handlers) GetGuest(w http.ResponseWriter, r *http.Request, eventId api.EventId, guestId api.GuestId) {
	guest, err := h.guests.GetGuest(r.Context(), eventId, guestId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "guest not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get guest")
		return
	}
	writeJSON(w, http.StatusOK, toAPIGuest(guest))
}

func (h *Handlers) UpdateGuest(w http.ResponseWriter, r *http.Request, eventId api.EventId, guestId api.GuestId) {
	var body api.NewGuest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	guest, err := h.guests.UpdateGuest(r.Context(), eventId, guestId, body.Name, toDomainEmail(body.Email), toDomainRsvpPtr(body.RsvpStatus))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "guest not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update guest")
		return
	}
	writeJSON(w, http.StatusOK, toAPIGuest(guest))
}

func (h *Handlers) DeleteGuest(w http.ResponseWriter, r *http.Request, eventId api.EventId, guestId api.GuestId) {
	if err := h.guests.DeleteGuest(r.Context(), eventId, guestId); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "guest not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete guest")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
