package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	adminhandler "github.com/j1hub/backend/internal/admin/adapter/http"
	adminpostgres "github.com/j1hub/backend/internal/admin/adapter/postgres"
	adminusecase "github.com/j1hub/backend/internal/admin/usecase"
	authhandler "github.com/j1hub/backend/internal/auth/adapter/http"
	authusecase "github.com/j1hub/backend/internal/auth/usecase"
	expensehandler "github.com/j1hub/backend/internal/expense/adapter/http"
	expensepostgres "github.com/j1hub/backend/internal/expense/adapter/postgres"
	expenseusecase "github.com/j1hub/backend/internal/expense/usecase"
	friendhandler "github.com/j1hub/backend/internal/friend/adapter/http"
	friendpostgres "github.com/j1hub/backend/internal/friend/adapter/postgres"
	friendusecase "github.com/j1hub/backend/internal/friend/usecase"
	journeyhandler "github.com/j1hub/backend/internal/gamification/adapter/http"
	gamificationpostgres "github.com/j1hub/backend/internal/gamification/adapter/postgres"
	gamificationusecase "github.com/j1hub/backend/internal/gamification/usecase"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/j1hub/backend/internal/infrastructure/db"
	"github.com/j1hub/backend/internal/infrastructure/notification"
	scraper "github.com/j1hub/backend/internal/infrastructure/outbound/scraper"
	"github.com/j1hub/backend/internal/infrastructure/scheduler"
	"github.com/j1hub/backend/internal/infrastructure/security"
	"github.com/j1hub/backend/internal/infrastructure/storage"
	jobhandler "github.com/j1hub/backend/internal/job/adapter/http"
	jobpostgres "github.com/j1hub/backend/internal/job/adapter/postgres"
	jobusecase "github.com/j1hub/backend/internal/job/usecase"
	"github.com/j1hub/backend/internal/media"
	missionhandler "github.com/j1hub/backend/internal/mission/adapter/http"
	missionpostgres "github.com/j1hub/backend/internal/mission/adapter/postgres"
	missionusecase "github.com/j1hub/backend/internal/mission/usecase"
	notifhandler "github.com/j1hub/backend/internal/notification/adapter/http"
	notificationpostgres "github.com/j1hub/backend/internal/notification/adapter/postgres"
	notificationusecase "github.com/j1hub/backend/internal/notification/usecase"
	transporthttp "github.com/j1hub/backend/internal/transport/http"
	userhandler "github.com/j1hub/backend/internal/user/adapter/http"
	userpostgres "github.com/j1hub/backend/internal/user/adapter/postgres"
	userusecase "github.com/j1hub/backend/internal/user/usecase"
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

	if err := db.SeedMockData(pool); err != nil {
		log.Printf("failed to seed mock data: %v", err)
	}

	clock := timeutil.RealClock{}

	// Repos
	userRepo := userpostgres.NewUserRepository(pool)
	profileRepo := userpostgres.NewProfileRepository(pool)
	creditRepo := gamificationpostgres.NewCreditScoreRepository(pool)
	phaseRepo := missionpostgres.NewJourneyPhaseRepository(pool)
	historyRepo := missionpostgres.NewUserPhaseHistoryRepository(pool)
	missionRepo := missionpostgres.NewMissionRepository(pool)
	umRepo := missionpostgres.NewUserMissionRepository(pool)
	ledgerRepo := gamificationpostgres.NewPointLedgerRepository(pool)
	splitRepo := expensepostgres.NewExpenseSplitRepository(pool)
	adminRepo := adminpostgres.NewAdminRepository(pool)
	friendRepo := friendpostgres.NewFriendshipRepository(pool)

	// Adapters
	hasher := security.NewArgon2Hasher()
	issuer := security.NewJWTIssuer(cfg)
	storage := storage.NewSupabaseStorage(cfg)
	notifier := notification.NewFCMNotifier(cfg, userRepo)

	// Usecases
	rewardEngine := gamificationusecase.NewRewardEngine(cfg, userRepo, umRepo)
	registerUC := authusecase.NewRegisterUserUseCase(pool, userRepo, profileRepo, creditRepo, phaseRepo, historyRepo, missionRepo, umRepo, hasher, issuer, clock)
	loginUC := authusecase.NewLoginUseCase(userRepo, hasher, issuer)
	userUC := userusecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	taskRepo := missionpostgres.NewTaskRepository(pool)
	utRepo := missionpostgres.NewUserTaskRepository(pool)
	badgeRepo := gamificationpostgres.NewBadgeRepository(pool)
	ubRepo := gamificationpostgres.NewUserBadgeRepository(pool)

	leaderRepo := gamificationpostgres.NewLeaderboardRepository(pool)

	missionUC := missionusecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)
	completeUC := missionusecase.NewCompleteMissionUseCase(umRepo, missionRepo, taskRepo, utRepo, userRepo, ledgerRepo, badgeRepo, ubRepo, storage, notifier, rewardEngine, clock)
	journeyUC := gamificationusecase.NewJourneyUseCase(phaseRepo, historyRepo, badgeRepo, ubRepo, creditRepo, ledgerRepo)
	advanceUC := gamificationusecase.NewAdvancePhaseUseCase(userRepo, umRepo, phaseRepo, historyRepo, missionRepo, notifier, clock)
	leaderboardUC := gamificationusecase.NewLeaderboardUseCase(leaderRepo, profileRepo, ubRepo)

	radarRepo := gamificationpostgres.NewRadarRepository(pool)

	friendshipUC := friendusecase.NewManageFriendshipUseCase(friendRepo, userRepo, notifier, clock)
	radarUC := gamificationusecase.NewRadarUseCase(cfg, profileRepo, radarRepo, friendRepo)

	txnRepo := expensepostgres.NewExpenseRepository(pool)
	expenseUC := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, storage, notifier, clock)

	notifRepo := notificationpostgres.NewNotificationRepository(pool)
	notifUC := notificationusecase.NewNotificationUseCase(notifRepo)
	adminUC := adminusecase.NewAdminUseCase(pool, adminRepo, userRepo, umRepo, missionRepo, ledgerRepo, notifier, rewardEngine, clock)

	jobRepo := jobpostgres.NewJobRepository(pool)
	housingRepo := jobpostgres.NewJobHousingRepository(pool)
	ratingRepo := jobpostgres.NewJobOverallRatingRepository(pool)
	reviewRepo := jobpostgres.NewJobReviewRepository(pool)
	cartRepo := jobpostgres.NewUserCartRepository(pool)

	jobUC := jobusecase.NewManageJobUseCase(jobRepo, housingRepo, ratingRepo, reviewRepo, cartRepo, clock)
	scrapeJobsUC := scraper.NewScrapeJobsUseCase(jobRepo, housingRepo)

	// Handlers
	authH := authhandler.NewAuthHandler(registerUC, loginUC)
	userH := userhandler.NewUserHandler(userUC)
	missionH := missionhandler.NewMissionHandler(missionUC, completeUC)
	journeyH := journeyhandler.NewJourneyHandler(journeyUC, advanceUC, leaderboardUC)
	friendH := friendhandler.NewFriendHandler(friendshipUC, radarUC)
	expenseH := expensehandler.NewExpenseHandler(expenseUC)
	notifH := notifhandler.NewNotificationHandler(notifUC)
	jobH := jobhandler.NewJobHandler(jobUC)
	adminH := adminhandler.NewAdminHandler(adminUC)
	mediaH := media.NewMediaHandler(storage)

	r := transporthttp.NewRouter(issuer, authH, userH, missionH, journeyH, friendH, expenseH, notifH, jobH, adminH, mediaH)

	// Jobs
	overdueExpenseJob := scheduler.NewOverdueExpenseJob(splitRepo, creditRepo, ledgerRepo, notifier)
	overdueMissionJob := scheduler.NewOverdueMissionJob(umRepo, missionRepo, userRepo, notifier)
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
