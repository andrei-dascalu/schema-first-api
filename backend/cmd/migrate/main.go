// Command migrate applies or rolls back database migrations against
// DATABASE_URL. It's a separate binary from cmd/api so that migrating the
// schema is always an explicit, standalone step rather than a side effect of
// starting the API server.
package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"spec-first-backend/internal/adapters/postgres"
)

func main() {
	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal().Msg("DATABASE_URL is required")
	}

	var err error
	switch direction {
	case "up":
		err = postgres.Migrate(databaseURL)
	case "down":
		err = postgres.MigrateDown(databaseURL)
	default:
		log.Fatal().Str("direction", direction).Msg(`unknown migration direction (want "up" or "down")`)
	}
	if err != nil {
		log.Fatal().Err(err).Str("direction", direction).Msg("migrate failed")
	}
	log.Info().Str("direction", direction).Msg("migrate done")
}
