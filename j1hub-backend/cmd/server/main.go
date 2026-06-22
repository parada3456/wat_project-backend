package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/j1hub/backend/internal/adapter/auth"
	"github.com/j1hub/backend/internal/adapter/http/handler"
	"github.com/j1hub/backend/internal/adapter/notification"
	"github.com/j1hub/backend/internal/adapter/postgres"
	"github.com/j1hub/backend/internal/adapter/storage"
	expenseusecase "github.com/j1hub/backend/internal/expense/usecase"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/j1hub/backend/internal/infrastructure/db"
	"github.com/j1hub/backend/internal/infrastructure/scheduler"
	jobusecase "github.com/j1hub/backend/internal/job/usecase"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/j1hub/backend/pkg/timeutil"
)

func main() {
	log.Println("debugprint: entering main")
	cfg := config.MustLoad()

	pool, err := db.NewPool(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(cfg.DatabaseURL, "migrations"); err != nil {
		log.Printf("migration error (might be expected if already applied): %v", err)
	}

	clock := timeutil.RealClock{}

	// Repos
	userRepo := postgres.NewUserRepository(pool)
	profileRepo := postgres.NewProfileRepository(pool)
	creditRepo := postgres.NewCreditScoreRepository(pool)
	phaseRepo := postgres.NewJourneyPhaseRepository(pool)
	historyRepo := postgres.NewUserPhaseHistoryRepository(pool)
	missionRepo := postgres.NewMissionRepository(pool)
	umRepo := postgres.NewUserMissionRepository(pool)
	ledgerRepo := postgres.NewPointLedgerRepository(pool)
	splitRepo := postgres.NewExpenseSplitRepository(pool)
	adminRepo := postgres.NewAdminRepository(pool)
	friendRepo := postgres.NewFriendshipRepository(pool)

	// Adapters
	hasher := auth.NewArgon2Hasher()
	issuer := auth.NewJWTIssuer(cfg)
	storage := storage.NewSupabaseStorage(cfg)
	notifier := notification.NewFCMNotifier(cfg, userRepo)

	// Usecases
	rewardEngine := usecase.NewRewardEngine(cfg, userRepo, umRepo)
	registerUC := authauthusecase.NewRegisterUserUseCase(pool, userRepo, profileRepo, creditRepo, phaseRepo, historyRepo, missionRepo, umRepo, hasher, issuer, clock)
	loginUC := usecase.NewLoginUseCase(userRepo, hasher, issuer)
	userUC := usecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	taskRepo := postgres.NewTaskRepository(pool)
	utRepo := postgres.NewUserTaskRepository(pool)
	badgeRepo := postgres.NewBadgeRepository(pool)
	ubRepo := postgres.NewUserBadgeRepository(pool)

	leaderRepo := postgres.NewLeaderboardRepository(pool)

	missionUC := usecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)
	completeUC := missionmissionusecase.NewCompleteMissionUseCase(umRepo, missionRepo, taskRepo, utRepo, userRepo, ledgerRepo, badgeRepo, ubRepo, storage, notifier, rewardEngine, clock)
	journeyUC := usecase.NewJourneyUseCase(phaseRepo, historyRepo, badgeRepo, ubRepo, creditRepo, ledgerRepo)
	advanceUC := usecase.NewAdvancePhaseUseCase(userRepo, umRepo, phaseRepo, historyRepo, missionRepo, notifier, clock)
	leaderboardUC := gamificationgamificationusecase.NewLeaderboardUseCase(leaderRepo, profileRepo, ubRepo)

	radarRepo := postgres.NewRadarRepository(pool)

	friendshipUC := friendfriendusecase.NewManageFriendshipUseCase(friendRepo, userRepo, notifier, clock)
	radarUC := gamificationgamificationusecase.NewRadarUseCase(cfg, profileRepo, radarRepo, friendRepo)

	txnRepo := postgres.NewExpenseRepository(pool)
	expenseUC := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, storage, notifier, clock)

	notifRepo := postgres.NewNotificationRepository(pool)
	notifUC := usecase.NewNotificationUseCase(notifRepo)
	adminUC := usecase.NewAdminUseCase(pool, adminRepo, userRepo, umRepo, missionRepo, ledgerRepo, notifier, rewardEngine, clock)

	jobRepo := postgres.NewJobRepository(pool)
	housingRepo := postgres.NewJobHousingRepository(pool)
	ratingRepo := postgres.NewJobOverallRatingRepository(pool)
	reviewRepo := postgres.NewJobReviewRepository(pool)
	cartRepo := postgres.NewUserCartRepository(pool)

	jobUC := jobusecase.NewManageJobUseCase(jobRepo, housingRepo, ratingRepo, reviewRepo, cartRepo, clock)
	scrapeJobsUC := usecase.NewScrapeJobsUseCase(jobRepo, housingRepo)

	// Handlers
	authH := handler.NewAuthHandler(registerUC, loginUC)
	userH := handler.NewUserHandler(userUC)
	missionH := handler.NewMissionHandler(missionUC, completeUC)
	journeyH := handler.NewJourneyHandler(journeyUC, advanceUC, leaderboardUC)
	friendH := handler.NewFriendHandler(friendshipUC, radarUC)
	expenseH := handler.NewExpenseHandler(expenseUC)
	notifH := handler.NewNotificationHandler(notifUC)
	jobH := handler.NewJobHandler(jobUC)
	adminH := handler.NewAdminHandler(adminUC)

	r := handler.NewRouter(issuer, authH, userH, missionH, journeyH, friendH, expenseH, notifH, jobH, adminH)

	// Jobs
	overdueExpenseJob := usecase.NewOverdueExpenseJob(splitRepo, creditRepo, ledgerRepo, notifier)
	overdueMissionJob := usecase.NewOverdueMissionJob(umRepo, missionRepo, userRepo, notifier)
	cron := scheduler.NewScheduler(cfg, overdueExpenseJob, overdueMissionJob, scrapeJobsUC)
	cron.Start()
	defer cron.Stop()

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
