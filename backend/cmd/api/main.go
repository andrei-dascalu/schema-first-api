package main

import (
	"context"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"spec-first-backend/internal/adapters/bcrypt"
	handlers "spec-first-backend/internal/adapters/http"
	"spec-first-backend/internal/adapters/jwttoken"
	"spec-first-backend/internal/adapters/postgres"
	"spec-first-backend/internal/api"
	"spec-first-backend/internal/app"
	"spec-first-backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("load config")
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("connect to database")
	}
	defer pool.Close()

	store := postgres.NewStore(pool)
	hasher := bcrypt.New()
	issuer := jwttoken.NewIssuer(cfg.JWTSecret)

	authService := app.NewAuthService(store, issuer, hasher)
	eventService := app.NewEventService(store)
	guestService := app.NewGuestService(store)

	h := handlers.New(authService, eventService, guestService)

	spec, err := api.GetSpec()
	if err != nil {
		log.Fatal().Err(err).Msg("load embedded spec")
	}
	spec.Servers = nil // servers in the spec would restrict validation to those hosts

	// Request bodies/params are validated against the embedded spec first;
	// actual bearer-token identity checks happen in handlers.Middleware, since
	// AuthenticationFunc here only needs to satisfy openapi3filter - see
	// docs/api_schema.yaml for the bearerAuth security scheme.
	validator := nethttpmiddleware.OapiRequestValidatorWithOptions(spec, &nethttpmiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
	})

	router := chi.NewRouter()

	// /metrics is intentionally outside the group below: it's not part of the
	// OpenAPI contract (so the spec validator would reject it) and is scraped
	// unauthenticated, so it must also skip the bearer-auth middleware.
	router.Handle("/metrics", promhttp.Handler())

	router.Group(func(r chi.Router) {
		// CORS is opt-in via CORS_ORIGINS (comma-separated) - e.g. for a
		// frontend dev server on a different origin/port. No policy is
		// applied at all when it's unset, rather than defaulting to some
		// baked-in origin list. This must run before the spec validator: it
		// short-circuits OPTIONS preflights itself, which otherwise wouldn't
		// match any documented operation and would fail validation.
		if len(cfg.CORSOrigins) > 0 {
			r.Use(corsMiddleware(cfg.CORSOrigins))
		}
		r.Use(handlers.MetricsMiddleware())
		r.Use(validator)
		r.Use(handlers.Middleware(issuer))

		api.HandlerFromMux(h, r)
	})

	log.Info().Strs("cors_origins", cfg.CORSOrigins).Msg("CORS configured")
	log.Info().Str("port", cfg.Port).Msg("listening")
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}

func corsMiddleware(origins []string) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
