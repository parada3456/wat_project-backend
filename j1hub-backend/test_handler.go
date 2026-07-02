package main

// import (
// 	"context"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"strings"

// 	"github.com/jackc/pgx/v5/pgxpool"
// 	missionpostgres "github.com/parada3456/wat_project-backend/internal/mission/adapter/postgres"
// 	missionusecase "github.com/parada3456/wat_project-backend/internal/mission/usecase"
// 	userpostgres "github.com/parada3456/wat_project-backend/internal/user/adapter/postgres"
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

// 	userRepo := userpostgres.NewUserRepository(pool)
// 	missionRepo := missionpostgres.NewMissionRepository(pool)
// 	umRepo := missionpostgres.NewUserMissionRepository(pool)
// 	taskRepo := missionpostgres.NewTaskRepository(pool)
// 	utRepo := missionpostgres.NewUserTaskRepository(pool)

// 	uc := missionusecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)

// 	fmt.Println("Testing ListStaticMissions for usr_alice...")
// 	missions, err := uc.ListStaticMissions(ctx, "usr_alice", nil)
// 	if err != nil {
// 		fmt.Printf("FAILED WITH ERROR: %v\n", err)
// 	} else {
// 		fmt.Printf("SUCCESS! Found %d missions for usr_alice\n", len(missions))
// 		for _, m := range missions {
// 			fmt.Printf(" - Mission: ID=%s, Title=%s, Phase=%s\n", m.MissionID, m.Title, m.PhaseID)
// 		}
// 	}

// 	fmt.Println("Testing ListAvailableMissions for usr_alice...")
// 	ums, err := uc.ListAvailableMissions(ctx, "usr_alice", nil)
// 	if err != nil {
// 		fmt.Printf("FAILED WITH ERROR: %v\n", err)
// 	} else {
// 		fmt.Printf("SUCCESS! Found %d user missions for usr_alice\n", len(ums))
// 		for _, um := range ums {
// 			fmt.Printf(" - UserMission: ID=%s, MissionID=%s, Status=%s\n", um.UserMissionID, um.MissionID, um.Status)
// 		}
// 	}
// }
