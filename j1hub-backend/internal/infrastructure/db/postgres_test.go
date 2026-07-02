package db_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/parada3456/wat_project-backend/internal/infrastructure/config"
	"github.com/parada3456/wat_project-backend/internal/infrastructure/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadEnv() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			content, err := os.ReadFile(envPath)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "#") {
						continue
					}
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						val := strings.Trim(parts[1], `"'`)
						os.Setenv(parts[0], val)
					}
				}
			}
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

func findMigrationsPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		migPath := filepath.Join(dir, "migrations")
		if _, err := os.Stat(migPath); err == nil {
			return migPath, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("migrations directory not found")
}

func TestPostgres_NewPool(t *testing.T) {
	loadEnv()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	// Test with invalid URL
	cfgInvalid := &config.Config{DatabaseURL: "invalid://url"}
	pool, err := db.NewPool(cfgInvalid)
	assert.Error(t, err)
	assert.Nil(t, pool)

	// Test with valid URL
	cfgValid := &config.Config{DatabaseURL: dbURL}
	pool, err = db.NewPool(cfgValid)
	assert.NoError(t, err)
	assert.NotNil(t, pool)
	if pool != nil {
		pool.Close()
	}
}

func TestPostgres_RunMigrations(t *testing.T) {
	loadEnv()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	migPath, err := findMigrationsPath()
	require.NoError(t, err)

	// Test invalid DB URL
	err = db.RunMigrations("postgres://invalid:invalid@localhost:5432/invalid", migPath)
	assert.Error(t, err)

	// Test valid migrations path and DB URL
	err = db.RunMigrations(dbURL, migPath)
	assert.NoError(t, err)
}
