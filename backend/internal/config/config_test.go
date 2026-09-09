package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadCORSOrigins(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("CORS_ORIGINS", " https://one.example.com, https://two.example.com ")

	cfg, err := Load()

	require.NoError(t, err)
	require.Equal(t, []string{"https://one.example.com", "https://two.example.com"}, cfg.CORSOrigins)
}

func TestCORSOriginSingularIsIgnored(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("CORS_ORIGINS", "")
	t.Setenv("CORS_ORIGIN", "https://singular.example.com")

	cfg, err := Load()

	require.NoError(t, err)
	require.Empty(t, cfg.CORSOrigins)
}
