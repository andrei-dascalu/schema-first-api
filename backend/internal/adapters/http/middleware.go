package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"spec-first-backend/internal/app"
)

type contextKey string

const userIDContextKey contextKey = "userID"

// publicPaths lists the operations marked `security: []` in
// docs/api_schema.yaml - oapi-codegen's generated chi router applies
// Middlewares uniformly to every route, so the exemption is handled here.
var publicPaths = map[string]bool{
	"/auth/register": true,
	"/auth/login":    true,
}

// Middleware validates the Authorization: Bearer <token> header and stores
// the authenticated user's ID in the request context. Matches the
// bearerAuth security scheme in docs/api_schema.yaml.
func Middleware(issuer app.TokenIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if publicPaths[r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			header := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" {
				writeUnauthorized(w)
				return
			}

			userID, err := issuer.Parse(token)
			if err != nil {
				writeUnauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	return userID, ok
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "missing or invalid credentials"})
}
