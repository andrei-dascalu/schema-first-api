package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"spec-first-backend/internal/api"
	"spec-first-backend/internal/domain"
)

func (h *Handlers) ListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.events.ListEvents(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list events")
		return
	}
	writeJSON(w, http.StatusOK, toAPIEvents(events))
}

func (h *Handlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var body api.NewEvent
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.events.CreateEvent(r.Context(), body.Name, body.Date, body.Location)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create event")
		return
	}
	writeJSON(w, http.StatusCreated, toAPIEvent(event))
}

func (h *Handlers) GetEvent(w http.ResponseWriter, r *http.Request, eventId api.EventId) {
	event, err := h.events.GetEvent(r.Context(), eventId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get event")
		return
	}
	writeJSON(w, http.StatusOK, toAPIEvent(event))
}

func (h *Handlers) UpdateEvent(w http.ResponseWriter, r *http.Request, eventId api.EventId) {
	var body api.NewEvent
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.events.UpdateEvent(r.Context(), eventId, body.Name, body.Date, body.Location)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update event")
		return
	}
	writeJSON(w, http.StatusOK, toAPIEvent(event))
}

func (h *Handlers) DeleteEvent(w http.ResponseWriter, r *http.Request, eventId api.EventId) {
	if err := h.events.DeleteEvent(r.Context(), eventId); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "event not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete event")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
