package main

// import (
// 	"context"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"strings"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func main() {
// 	bytes, err := ioutil.ReadFile(".env")
// 	if err != nil {
// 		log.Fatalf("Failed to read .env file: %v", err)
// 	}

// 	var dbURL string
// 	lines := strings.Split(string(bytes), "\n")
// 	for _, line := range lines {
// 		line = strings.TrimSpace(line)
// 		if strings.HasPrefix(line, "DATABASE_URL=") {
// 			dbURL = strings.TrimPrefix(line, "DATABASE_URL=")
// 			break
// 		}
// 	}

// 	ctx := context.Background()
// 	pool, err := pgxpool.New(ctx, dbURL)
// 	if err != nil {
// 		log.Fatalf("Unable to create pool: %v", err)
// 	}
// 	defer pool.Close()

// 	var userCount int
// 	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
// 	if err != nil {
// 		log.Fatalf("Query failed: %v", err)
// 	}
// 	fmt.Printf("Total users in database: %d\n", userCount)

// 	rows, err := pool.Query(ctx, "SELECT user_id, email, current_phase_id FROM users")
// 	if err != nil {
// 		log.Fatalf("Query failed: %v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var id, email string
// 		var phase *string
// 		if err := rows.Scan(&id, &email, &phase); err != nil {
// 			log.Fatalf("Scan failed: %v", err)
// 		}
// 		pStr := "NULL"
// 		if phase != nil {
// 			pStr = *phase
// 		}
// 		fmt.Printf("User: ID=%s, Email=%s, Phase=%s\n", id, email, pStr)
// 	}
// }
