package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/j1hub/backend/internal/infrastructure/db"
)

func main() {
	bytes, err := ioutil.ReadFile(".env")
	if err != nil {
		log.Fatalf("Failed to read .env file: %v", err)
	}

	var dbURL string
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "DATABASE_URL=") {
			dbURL = strings.TrimPrefix(line, "DATABASE_URL=")
			break
		}
	}

	if dbURL == "" {
		log.Fatal("DATABASE_URL not found in .env")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to create pool: %v", err)
	}
	defer pool.Close()

	fmt.Println("Running SeedMockData...")
	err = db.SeedMockData(pool)
	if err != nil {
		fmt.Printf("SEEDING FAILED WITH ERROR:\n%v\n", err)
		os.Exit(1)
	}
	fmt.Println("Seeding completed successfully!")
}
