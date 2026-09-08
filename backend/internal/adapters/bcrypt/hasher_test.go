package bcrypt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spec-first-backend/internal/adapters/bcrypt"
)

func TestHashAndCompare(t *testing.T) {
	hasher := bcrypt.New()

	hash, err := hasher.Hash("hunter22")
	require.NoError(t, err)

	assert.True(t, hasher.Compare(hash, "hunter22"))
	assert.False(t, hasher.Compare(hash, "wrong-password"))
}
