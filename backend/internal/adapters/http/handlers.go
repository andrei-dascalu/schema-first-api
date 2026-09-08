// Package handlers implements api.ServerInterface, the oapi-codegen
// generated interface for every operation in docs/api_schema.yaml, and owns
// the bearer-auth and request-validation middleware. It translates
// api.* <-> domain.* and calls internal/app services - no SQL here.
package handlers

import (
	"encoding/json"
	"net/http"

	"spec-first-backend/internal/app"
)

type Handlers struct {
	auth   *app.AuthService
	events *app.EventService
	guests *app.GuestService
}

func New(auth *app.AuthService, events *app.EventService, guests *app.GuestService) *Handlers {
	return &Handlers{auth: auth, events: events, guests: guests}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"message": message})
}
