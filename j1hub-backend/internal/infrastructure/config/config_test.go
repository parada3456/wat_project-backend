package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/parada3456/wat_project-backend/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig_MustLoad(t *testing.T) {
	err := os.WriteFile(".env", []byte("DATABASE_URL=postgres://localhost/test\nJWT_SECRET=secret\nJWT_EXPIRY_HOURS=24\n"), 0644)
	if err != nil {
		t.Fatalf("failed to create temporary .env: %v", err)
	}
	defer os.Remove(".env")

	cfg := config.MustLoad()
	assert.NotNil(t, cfg)
	assert.Equal(t, "postgres://localhost/test", cfg.DatabaseURL)
	assert.Equal(t, 24*time.Hour, cfg.JWTExpiry())
}
