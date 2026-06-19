package db

import (
	"context"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(cfg *config.Config) (*pgxpool.Pool, error) {
	log.Println("debugprint: entering NewPool")
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Printf("Connected to database. Pool stats: %+v", pool.Stat())
	return pool, nil
}

func RunMigrations(databaseURL string, migrationsPath string) error {
	log.Println("debugprint: entering RunMigrations")
	m, err := migrate.New("file://"+migrationsPath, databaseURL)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run up migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
