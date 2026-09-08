package jwttoken_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spec-first-backend/internal/adapters/jwttoken"
)

func TestIssueAndParseRoundTrip(t *testing.T) {
	issuer := jwttoken.NewIssuer("test-secret")
	userID := uuid.New()

	token, expiresAt, err := issuer.Issue(userID)
	require.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(jwttoken.TokenTTL), expiresAt, time.Second)

	parsed, err := issuer.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, userID, parsed)
}

func TestParseRejectsTokenFromDifferentSecret(t *testing.T) {
	token, _, err := jwttoken.NewIssuer("secret-a").Issue(uuid.New())
	require.NoError(t, err)

	_, err = jwttoken.NewIssuer("secret-b").Parse(token)
	assert.ErrorIs(t, err, jwttoken.ErrInvalidToken)
}

func TestParseRejectsGarbage(t *testing.T) {
	_, err := jwttoken.NewIssuer("test-secret").Parse("not-a-token")
	assert.ErrorIs(t, err, jwttoken.ErrInvalidToken)
}
