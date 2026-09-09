package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestCORSMiddlewareAllowsConfiguredOrigin(t *testing.T) {
	router := chi.NewRouter()
	router.Use(corsMiddleware([]string{"https://frontend.example.com"}))
	router.Get("/resource", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodOptions, "/resource", nil)
	req.Header.Set("Origin", "https://frontend.example.com")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, req)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "https://frontend.example.com", response.Header().Get("Access-Control-Allow-Origin"))
}
