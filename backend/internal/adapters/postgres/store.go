package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"spec-first-backend/internal/store/postgres/db"
)

type Store struct {
	q *db.Queries
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{q: db.New(pool)}
}
