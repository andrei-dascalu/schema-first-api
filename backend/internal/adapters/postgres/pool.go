// Package postgres implements the app repository ports against Postgres.
// Table/query definitions live in migrations/ (applied with goose, and
// embedded here so Migrate can run them) and queries/ (compiled by sqlc
// into internal/store/postgres/db) - see that package for the generated
// code and the Makefile for the `migrate-*` and `generate` targets.
package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib" // database/sql driver used by goose
)

//go:embed migrations/*.sql
var migrations embed.FS

// Migrate applies all pending goose migrations found in migrations/. Callers
// invoke this explicitly (see cmd/migrate) - the API server does not run
// migrations on startup.
func Migrate(databaseURL string) error {
	return withGooseDB(databaseURL, func(db *sql.DB) error {
		if err := goose.Up(db, "migrations"); err != nil {
			return fmt.Errorf("run migrations: %w", err)
		}
		return nil
	})
}

// MigrateDown rolls back the most recently applied migration.
func MigrateDown(databaseURL string) error {
	return withGooseDB(databaseURL, func(db *sql.DB) error {
		if err := goose.Down(db, "migrations"); err != nil {
			return fmt.Errorf("roll back migration: %w", err)
		}
		return nil
	})
}

func withGooseDB(databaseURL string, fn func(*sql.DB) error) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open migration connection: %w", err)
	}
	defer func() { _ = db.Close() }()

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	return fn(db)
}

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}
