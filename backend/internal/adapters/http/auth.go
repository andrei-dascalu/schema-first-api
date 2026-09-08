package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"spec-first-backend/internal/api"
	"spec-first-backend/internal/domain"
)

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var body api.NewUser
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.auth.Register(r.Context(), string(body.Email), body.Password, body.Name)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			writeError(w, http.StatusConflict, "email already registered")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	writeJSON(w, http.StatusCreated, toAPIUser(user))
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var body api.Credentials
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, expiresAt, err := h.auth.Login(r.Context(), string(body.Email), body.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	writeJSON(w, http.StatusOK, api.LoginResponse{Token: token, ExpiresAt: expiresAt})
}
