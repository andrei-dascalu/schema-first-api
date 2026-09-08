//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"spec-first-backend/internal/adapters/postgres"
	"spec-first-backend/internal/domain"
)

func newTestStore(t *testing.T) *postgres.Store {
	t.Helper()
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx, "docker.io/postgres:16-alpine",
		tcpostgres.WithDatabase("spec_first_test"),
		tcpostgres.WithUsername("spec_first"),
		tcpostgres.WithPassword("spec_first"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(context.Background()) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	require.NoError(t, postgres.Migrate(dsn))

	pool, err := postgres.NewPool(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	return postgres.NewStore(pool)
}

func TestEventAndGuestLifecycle(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	location := "HQ"
	event, err := store.CreateEvent(ctx, "Launch Party", time.Now().Add(24*time.Hour).UTC(), &location)
	require.NoError(t, err)
	require.Equal(t, "Launch Party", event.Name)

	guest, err := store.CreateGuest(ctx, event.ID, "Bob", nil, domain.RsvpPending)
	require.NoError(t, err)
	require.Equal(t, domain.RsvpPending, guest.RsvpStatus)

	guests, err := store.ListGuestsByEvent(ctx, event.ID)
	require.NoError(t, err)
	require.Len(t, guests, 1)

	updated, err := store.UpdateGuest(ctx, event.ID, guest.ID, "Bob", nil, domain.RsvpAccepted)
	require.NoError(t, err)
	require.Equal(t, domain.RsvpAccepted, updated.RsvpStatus)

	require.NoError(t, store.DeleteGuest(ctx, event.ID, guest.ID))
	require.NoError(t, store.DeleteEvent(ctx, event.ID))

	_, err = store.GetEvent(ctx, event.ID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestCreateUserRejectsDuplicateEmail(t *testing.T) {
	store := newTestStore(t)
	ctx := context.Background()

	_, err := store.CreateUser(ctx, "alice@example.com", "hash", "Alice")
	require.NoError(t, err)

	_, err = store.CreateUser(ctx, "alice@example.com", "hash", "Alice Again")
	require.ErrorIs(t, err, domain.ErrDuplicateEmail)
}
