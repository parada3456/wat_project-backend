package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/j1hub/backend/internal/adapter/postgres"
	"github.com/j1hub/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
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

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	loadEnv()
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	// Clean/wipe database schema to resolve any dirty migration state
	sqlDB, err := sql.Open("postgres", dbURL)
	require.NoError(t, err)
	_, err = sqlDB.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	require.NoError(t, err)
	sqlDB.Close()

	migPath, err := findMigrationsPath()
	require.NoError(t, err)

	// golang-migrate requires database driver configuration
	m, err := migrate.New("file://"+migPath, dbURL)
	require.NoError(t, err)
	
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}
	
	dbErr1, dbErr2 := m.Close()
	require.NoError(t, dbErr1)
	require.NoError(t, dbErr2)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)

	cleanup := func() {
		// Truncate tables to clean up
		tables := []string{
			"users", "profiles", "friendships", "journey_phases",
			"user_phase_history", "missions", "user_missions",
			"tasks", "user_tasks", "point_ledger", "badges",
			"user_badges", "credit_scores", "expense_splits", "expense_transactions",
			"user_carts", "job_reviews", "job_overall_ratings", "job_housings", "job_postings",
		}
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", strings.Join(tables, ", "))
		_, err := pool.Exec(ctx, query)
		if err != nil {
			t.Logf("Failed to truncate tables: %v", err)
		}
		pool.Close()
	}

	return pool, cleanup
}

func TestUserRepository(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	repo := postgres.NewUserRepository(pool)
	ctx := context.Background()

	// Insert dummy phases first to avoid FK constraint violation
	_, err := pool.Exec(ctx, `INSERT INTO journey_phases (phase_id, phase_number, title) VALUES ('phase_1', 1, 'Phase 1'), ('phase_2', 2, 'Phase 2')`)
	require.NoError(t, err)

	u := &domain.User{
		UserID:              "usr_1",
		Email:               "user1@test.com",
		PasswordHash:        "hash1",
		FirstName:           "First",
		LastName:            "Last",
		CurrentPhaseID:      "phase_1",
		TotalLifetimePoints: 10,
		CurrentPhasePoints:  5,
		MissionStreak:       1,
		ArrivalDate:         time.Now().Truncate(time.Second),
		JobStartDate:        time.Now().Truncate(time.Second),
		CreatedAt:           time.Now().Truncate(time.Second),
		UpdatedAt:           time.Now().Truncate(time.Second),
	}

	// Create
	err = repo.Create(ctx, u)
	assert.NoError(t, err)

	// FindByID
	found, err := repo.FindByID(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, u.Email, found.Email)
	assert.Equal(t, u.FirstName, found.FirstName)

	// FindByEmail
	foundEmail, err := repo.FindByEmail(ctx, "user1@test.com")
	assert.NoError(t, err)
	assert.Equal(t, u.UserID, foundEmail.UserID)

	// Update
	u.FirstName = "Updated"
	err = repo.Update(ctx, u)
	assert.NoError(t, err)
	foundUpdated, err := repo.FindByID(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, "Updated", foundUpdated.FirstName)

	// IncrementPoints
	err = repo.IncrementPoints(ctx, "usr_1", 20, 15)
	assert.NoError(t, err)
	foundPoints, err := repo.FindByID(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, 30, foundPoints.TotalLifetimePoints)
	assert.Equal(t, 20, foundPoints.CurrentPhasePoints)

	// ResetStreak
	err = repo.ResetStreak(ctx, "usr_1")
	assert.NoError(t, err)
	foundStreak, err := repo.FindByID(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, 0, foundStreak.MissionStreak)

	// SetPhase
	err = repo.SetPhase(ctx, "usr_1", "phase_2")
	assert.NoError(t, err)
	foundPhase, err := repo.FindByID(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, "phase_2", foundPhase.CurrentPhaseID)

	// Delete
	err = repo.Delete(ctx, "usr_1")
	assert.NoError(t, err)
	_, err = repo.FindByID(ctx, "usr_1")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestProfileRepository(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(pool)
	repo := postgres.NewProfileRepository(pool)
	ctx := context.Background()

	// Insert dummy phases first to avoid FK constraint violation
	_, err := pool.Exec(ctx, `INSERT INTO journey_phases (phase_id, phase_number, title) VALUES ('phase_1', 1, 'Phase 1'), ('phase_2', 2, 'Phase 2')`)
	require.NoError(t, err)

	// Need a user first
	u := &domain.User{
		UserID:         "usr_1",
		Email:          "user1@test.com",
		PasswordHash:   "hash",
		FirstName:      "First",
		LastName:       "Last",
		CurrentPhaseID: "phase_1",
	}
	err = userRepo.Create(ctx, u)
	require.NoError(t, err)

	p := &domain.Profile{
		ProfileID:         "prof_1",
		UserID:            "usr_1",
		PhoneNumber:       "123456",
		Bio:               "My bio",
		AvatarURL:         "avatar.jpg",
		RadarVisibility:   domain.VisibilityShowFriends,
		Lat:               1.23,
		Lng:               4.56,
		LocationUpdatedAt: time.Now().Truncate(time.Second),
		UpdatedAt:         time.Now().Truncate(time.Second),
	}

	// Create
	err = repo.Create(ctx, p)
	assert.NoError(t, err)

	// FindByUserID
	found, err := repo.FindByUserID(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, p.Bio, found.Bio)

	// Update
	p.Bio = "New Bio"
	err = repo.Update(ctx, p)
	assert.NoError(t, err)
	foundUpdated, err := repo.FindByUserID(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, "New Bio", foundUpdated.Bio)

	// UpdateLocation
	err = repo.UpdateLocation(ctx, "usr_1", 2.34, 5.67)
	assert.NoError(t, err)
	foundLoc, err := repo.FindByUserID(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, 2.34, foundLoc.Lat)
	assert.Equal(t, 5.67, foundLoc.Lng)

	// UpdateVisibility
	err = repo.UpdateVisibility(ctx, "usr_1", domain.VisibilityHidden)
	assert.NoError(t, err)
	foundVis, err := repo.FindByUserID(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, domain.VisibilityHidden, foundVis.RadarVisibility)
}

func TestFriendshipRepository(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(pool)
	repo := postgres.NewFriendshipRepository(pool)
	ctx := context.Background()

	_, err := pool.Exec(ctx, `INSERT INTO journey_phases (phase_id, phase_number, title) VALUES ('phase_1', 1, 'Phase 1')`)
	require.NoError(t, err)

	u1 := &domain.User{UserID: "usr_1", Email: "u1@t.com", CurrentPhaseID: "phase_1"}
	u2 := &domain.User{UserID: "usr_2", Email: "u2@t.com", CurrentPhaseID: "phase_1"}
	require.NoError(t, userRepo.Create(ctx, u1))
	require.NoError(t, userRepo.Create(ctx, u2))

	f := &domain.Friendship{
		FriendshipID: "fr_1",
		UserID1:      "usr_1",
		UserID2:      "usr_2",
		Status:       domain.FriendshipPending,
		CreatedAt:    time.Now().Truncate(time.Second),
		UpdatedAt:    time.Now().Truncate(time.Second),
	}

	err = repo.Insert(ctx, f)
	assert.NoError(t, err)

	found, err := repo.FindByID(ctx, "fr_1")
	assert.NoError(t, err)
	assert.Equal(t, f.Status, found.Status)

	foundPair, err := repo.FindByCanonicalPair(ctx, "usr_1", "usr_2")
	assert.NoError(t, err)
	assert.Equal(t, "fr_1", foundPair.FriendshipID)

	err = repo.UpdateStatus(ctx, "fr_1", domain.FriendshipAccepted)
	assert.NoError(t, err)

	friends, err := repo.FindFriendsOf(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Len(t, friends, 1)
}

func TestNotificationRepository(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()
	repo := postgres.NewNotificationRepository(pool)
	ctx := context.Background()

	err := repo.Insert(ctx, &domain.Notification{})
	assert.NoError(t, err)

	list, err := repo.FindByUser(ctx, "user")
	assert.Nil(t, list)
	assert.NoError(t, err)

	err = repo.MarkAsRead(ctx, "123")
	assert.NoError(t, err)

	err = repo.MarkAllAsRead(ctx, "user")
	assert.NoError(t, err)

	err = repo.Delete(ctx, "123")
	assert.NoError(t, err)
}

func TestRadarRepository(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(pool)
	profileRepo := postgres.NewProfileRepository(pool)
	radarRepo := postgres.NewRadarRepository(pool)
	ctx := context.Background()

	_, err := pool.Exec(ctx, `INSERT INTO journey_phases (phase_id, phase_number, title) VALUES ('phase_1', 1, 'Phase 1')`)
	require.NoError(t, err)

	u1 := &domain.User{UserID: "usr_1", Email: "u1@t.com", CurrentPhaseID: "phase_1"}
	err = userRepo.Create(ctx, u1)
	require.NoError(t, err)

	p1 := &domain.Profile{
		ProfileID:       "prof_1",
		UserID:          "usr_1",
		RadarVisibility: domain.VisibilityShowFriends,
		Lat:             1.23,
		Lng:             4.56,
	}
	err = profileRepo.Create(ctx, p1)
	require.NoError(t, err)

	err = profileRepo.UpdateLocation(ctx, "usr_1", 1.23, 4.56)
	require.NoError(t, err)

	profiles, err := radarRepo.FindNearby(ctx, 1.23, 4.56, 5000, 30)
	assert.NoError(t, err)
	assert.Len(t, profiles, 1)
	assert.Equal(t, "usr_1", profiles[0].UserID)
}

func TestLeaderboardRepository(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(pool)
	leaderboardRepo := postgres.NewLeaderboardRepository(pool)
	ctx := context.Background()

	_, err := pool.Exec(ctx, `INSERT INTO journey_phases (phase_id, phase_number, title) VALUES ('phase_1', 1, 'Phase 1')`)
	require.NoError(t, err)

	u1 := &domain.User{UserID: "usr_1", Email: "u1@t.com", CurrentPhaseID: "phase_1", CurrentPhasePoints: 100}
	u2 := &domain.User{UserID: "usr_2", Email: "u2@t.com", CurrentPhaseID: "phase_1", CurrentPhasePoints: 200}
	require.NoError(t, userRepo.Create(ctx, u1))
	require.NoError(t, userRepo.Create(ctx, u2))

	users, err := leaderboardRepo.FindByScope(ctx, "global", "")
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "usr_2", users[0].UserID)
	assert.Equal(t, "usr_1", users[1].UserID)
}

func TestExpenseRepository(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(pool)
	expenseRepo := postgres.NewExpenseRepository(pool)
	splitRepo := postgres.NewExpenseSplitRepository(pool)
	ctx := context.Background()

	_, err := pool.Exec(ctx, `INSERT INTO journey_phases (phase_id, phase_number, title) VALUES ('phase_1', 1, 'Phase 1')`)
	require.NoError(t, err)

	u1 := &domain.User{UserID: "usr_1", Email: "u1@t.com", CurrentPhaseID: "phase_1"}
	u2 := &domain.User{UserID: "usr_2", Email: "u2@t.com", CurrentPhaseID: "phase_1"}
	require.NoError(t, userRepo.Create(ctx, u1))
	require.NoError(t, userRepo.Create(ctx, u2))

	tx := &domain.ExpenseTransaction{
		TransactionID:   "tx_1",
		PaidByUserID:    "usr_1",
		Title:           "Dinner",
		TotalAmount:     100.0,
		Currency:        "USD",
		Memo:            "Shared dinner",
		TransactionDate: time.Now().Truncate(time.Second),
		DueDate:         time.Now().Add(-1 * time.Hour).Truncate(time.Second),
		CreatedAt:       time.Now().Truncate(time.Second),
		UpdatedAt:       time.Now().Truncate(time.Second),
	}

	err = expenseRepo.Insert(ctx, tx)
	assert.NoError(t, err)

	foundTx, err := expenseRepo.FindByID(ctx, "tx_1")
	assert.NoError(t, err)
	assert.Equal(t, tx.Title, foundTx.Title)

	err = expenseRepo.MarkSettled(ctx, "tx_1")
	assert.NoError(t, err)

	splits := []domain.ExpenseSplit{
		{
			SplitID:        "sp_1",
			TransactionID:  "tx_1",
			UserID:         "usr_2",
			OweAmount:      50.0,
			PaymentStatus:  domain.PaymentPending,
			ApprovalStatus: domain.ApprovalPending,
			UpdatedAt:      time.Now().Truncate(time.Second),
		},
	}

	err = splitRepo.BulkInsert(ctx, splits)
	assert.NoError(t, err)

	foundSplit, err := splitRepo.FindByID(ctx, "sp_1")
	assert.NoError(t, err)
	assert.Equal(t, 50.0, foundSplit.OweAmount)

	err = splitRepo.UpdatePaymentStatus(ctx, "sp_1", domain.PaymentSubmitted, "http://slip.png")
	assert.NoError(t, err)

	now := time.Now().Truncate(time.Second)
	err = splitRepo.UpdateApproval(ctx, "sp_1", domain.ApprovalApproved, &now)
	assert.NoError(t, err)

	count, err := splitRepo.CountUnsettled(ctx, "tx_1")
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	tx2 := &domain.ExpenseTransaction{
		TransactionID:   "tx_2",
		PaidByUserID:    "usr_1",
		Title:           "Lunch",
		TotalAmount:     40.0,
		Currency:        "USD",
		DueDate:         time.Now().Add(-24 * time.Hour).Truncate(time.Second),
		TransactionDate: time.Now().Truncate(time.Second),
		CreatedAt:       time.Now().Truncate(time.Second),
		UpdatedAt:       time.Now().Truncate(time.Second),
	}
	require.NoError(t, expenseRepo.Insert(ctx, tx2))

	splits2 := []domain.ExpenseSplit{
		{
			SplitID:        "sp_2",
			TransactionID:  "tx_2",
			UserID:         "usr_2",
			OweAmount:      20.0,
			PaymentStatus:  domain.PaymentPending,
			ApprovalStatus: domain.ApprovalPending,
			UpdatedAt:      time.Now().Truncate(time.Second),
		},
	}
	require.NoError(t, splitRepo.BulkInsert(ctx, splits2))

	overdues, err := splitRepo.FindOverdue(ctx)
	assert.NoError(t, err)
	assert.Len(t, overdues, 1)
	assert.Equal(t, "sp_2", overdues[0].SplitID)
}

func TestJobRepository(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	jobRepo := postgres.NewJobRepository(pool)
	housingRepo := postgres.NewJobHousingRepository(pool)
	ratingRepo := postgres.NewJobOverallRatingRepository(pool)
	reviewRepo := postgres.NewJobReviewRepository(pool)
	cartRepo := postgres.NewUserCartRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	ctx := context.Background()

	_, err := pool.Exec(ctx, `INSERT INTO journey_phases (phase_id, phase_number, title) VALUES ('phase_1', 1, 'Phase 1')`)
	require.NoError(t, err)
	u1 := &domain.User{UserID: "usr_1", Email: "u1@t.com", CurrentPhaseID: "phase_1"}
	require.NoError(t, userRepo.Create(ctx, u1))

	_, err = pool.Exec(ctx, `
		INSERT INTO job_postings (
			job_id, agency_name, employer_title, position, position_type, 
			location_city, location_state, group_location, us_sponsor, 
			salary_range_min, salary_range_max, available_slots, description, 
			source_url, scrape_at, posted_at, updated_at
		) VALUES (
			'job_1', 'Agency A', 'Title T', 'Position P', 'Full-Time', 
			'City C', 'State S', 'Group G', true, 
			15.0, 20.0, 5, 'Desc', 
			'http://source', NOW(), NOW(), NOW()
		)`)
	require.NoError(t, err)

	jobs, err := jobRepo.FindWithFilters(ctx, map[string]interface{}{"position_type": "Full-Time"})
	assert.NoError(t, err)
	assert.Len(t, jobs, 1)
	assert.Equal(t, "job_1", jobs[0].JobID)

	job, err := jobRepo.FindByID(ctx, "job_1")
	assert.NoError(t, err)
	assert.Equal(t, "Agency A", job.AgencyName)

	_, err = pool.Exec(ctx, `INSERT INTO job_housings (housing_id, job_id, description, weekly_rate, deposit, transportation, range_min_start_date, range_max_start_date, created_at, updated_at) VALUES ('h_1', 'job_1', 'Desc', 150.0, 200.0, 'Bus', NOW(), NOW(), NOW(), NOW())`)
	require.NoError(t, err)

	housings, err := housingRepo.FindByJobID(ctx, "job_1")
	assert.NoError(t, err)
	assert.Len(t, housings, 1)
	assert.Equal(t, "h_1", housings[0].HousingID)

	rv := &domain.JobReview{
		ReviewID:                  "rev_1",
		JobID:                     "job_1",
		UserID:                    "usr_1",
		RatingStars:               4.5,
		ReviewText:                "Great",
		TipsForNextGeneration:     "None",
		ScoreAgency:               4.0,
		ScoreJob:                  4.0,
		ScoreCoworkers:            4.0,
		ScoreTown:                 4.0,
		ScoreHours:                4.0,
		ScoreHousing:              4.0,
		ScoreSecondJobFeasibility: 4.0,
		ScoreOvertimeAvailability: 4.0,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}
	err = reviewRepo.Insert(ctx, rv)
	assert.NoError(t, err)

	reviews, err := reviewRepo.FindByJobID(ctx, "job_1")
	assert.NoError(t, err)
	assert.Len(t, reviews, 1)

	err = ratingRepo.Recalculate(ctx, "job_1")
	assert.NoError(t, err)

	rating, err := ratingRepo.FindByJobID(ctx, "job_1")
	assert.NoError(t, err)
	assert.Equal(t, 4.5, rating.OverallRate)

	cartItem := &domain.UserCart{
		CartID:    "cart_1",
		UserID:    "usr_1",
		JobID:     "job_1",
		Status:    domain.CartSaved,
		AddedAt:   time.Now(),
		UpdatedAt: time.Now(),
	}
	err = cartRepo.Insert(ctx, cartItem)
	assert.NoError(t, err)

	foundCart, err := cartRepo.FindByUserAndJob(ctx, "usr_1", "job_1")
	assert.NoError(t, err)
	assert.Equal(t, "cart_1", foundCart.CartID)

	foundCartByID, err := cartRepo.FindByID(ctx, "cart_1")
	assert.NoError(t, err)
	assert.Equal(t, "cart_1", foundCartByID.CartID)

	err = cartRepo.UpdateStatus(ctx, "cart_1", domain.CartApplied)
	assert.NoError(t, err)
}

func TestMissionRepository(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := postgres.NewUserRepository(pool)
	missionRepo := postgres.NewMissionRepository(pool)
	umRepo := postgres.NewUserMissionRepository(pool)
	taskRepo := postgres.NewTaskRepository(pool)
	utRepo := postgres.NewUserTaskRepository(pool)
	phaseRepo := postgres.NewJourneyPhaseRepository(pool)
	historyRepo := postgres.NewUserPhaseHistoryRepository(pool)
	ctx := context.Background()

	_, err := pool.Exec(ctx, `INSERT INTO journey_phases (phase_id, phase_number, title) VALUES ('phase_1', 1, 'Phase 1'), ('phase_2', 2, 'Phase 2')`)
	require.NoError(t, err)

	u1 := &domain.User{UserID: "usr_1", Email: "u1@t.com", CurrentPhaseID: "phase_1"}
	require.NoError(t, userRepo.Create(ctx, u1))

	ph, err := phaseRepo.FindByID(ctx, "phase_1")
	assert.NoError(t, err)
	assert.Equal(t, 1, ph.PhaseNumber)

	phNum, err := phaseRepo.FindByNumber(ctx, 2)
	assert.NoError(t, err)
	assert.Equal(t, "phase_2", phNum.PhaseID)

	hist := &domain.UserPhaseHistory{
		HistoryID:         "hist_1",
		UserID:            "usr_1",
		PhaseID:           "phase_1",
		PhasePointsEarned: 10,
		EnteredAt:         time.Now().Truncate(time.Second),
	}
	err = historyRepo.Insert(ctx, hist)
	assert.NoError(t, err)

	err = historyRepo.CompleteCurrentPhase(ctx, "usr_1", 50, time.Now().Truncate(time.Second))
	assert.NoError(t, err)

	foundHist, err := historyRepo.FindByUserAndPhase(ctx, "usr_1", "phase_1")
	assert.NoError(t, err)
	assert.Equal(t, 50, foundHist.PhasePointsEarned)

	_, err = pool.Exec(ctx, `INSERT INTO missions (mission_id, phase_id, title, description, location, base_points, is_mandatory, verification_type, due_date_type, fixed_due_date, relative_trigger_event, relative_days_offset, created_at, updated_at) VALUES ('m_1', 'phase_1', 'M1', 'Desc', 'Loc', 100, true, 'None', 'Relative', NULL, 'arrival_date', 5, NOW(), NOW())`)
	require.NoError(t, err)

	missions, err := missionRepo.FindByPhase(ctx, "phase_1")
	assert.NoError(t, err)
	assert.Len(t, missions, 1)

	m, err := missionRepo.FindByID(ctx, "m_1")
	assert.NoError(t, err)
	assert.Equal(t, "M1", m.Title)

	ums := []domain.UserMission{
		{
			UserMissionID:     "um_1",
			UserID:            "usr_1",
			MissionID:         "m_1",
			Status:            domain.StatusInProgress,
			CalculatedDueDate: time.Now().Add(-1 * time.Hour).Truncate(time.Second),
			CreatedAt:         time.Now().Truncate(time.Second),
			UpdatedAt:         time.Now().Truncate(time.Second),
		},
	}
	err = umRepo.BulkInsert(ctx, ums)
	assert.NoError(t, err)

	foundUM, err := umRepo.FindByID(ctx, "um_1")
	assert.NoError(t, err)
	assert.Equal(t, domain.StatusInProgress, foundUM.Status)

	foundUMs, err := umRepo.FindByUserAndPhase(ctx, "usr_1", "phase_1")
	assert.NoError(t, err)
	assert.Len(t, foundUMs, 1)

	err = umRepo.UpdateStatus(ctx, "um_1", domain.StatusPendingVerification)
	assert.NoError(t, err)

	err = umRepo.UpdateVerification(ctx, "um_1", time.Now().Truncate(time.Second), "admin")
	assert.NoError(t, err)

	reward := &domain.PointReward{Base: 100, SpeedBonus: 10, StreakBonus: 5, Total: 115}
	err = umRepo.UpdateReward(ctx, "um_1", reward, time.Now().Truncate(time.Second))
	assert.NoError(t, err)

	overdues, err := umRepo.FindOverdue(ctx)
	assert.NoError(t, err)
	assert.Len(t, overdues, 1)

	_, err = pool.Exec(ctx, `INSERT INTO tasks (task_id, mission_id, title, description, created_at, updated_at) VALUES ('t_1', 'm_1', 'T1', 'Desc', NOW(), NOW())`)
	require.NoError(t, err)

	tasks, err := taskRepo.FindByMission(ctx, "m_1")
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)

	ut := &domain.UserTask{
		UserTaskID:    "ut_1",
		UserID:        "usr_1",
		TaskID:        "t_1",
		UserMissionID: "um_1",
		IsCompleted:   true,
		UpdatedAt:     time.Now().Truncate(time.Second),
	}
	err = utRepo.Upsert(ctx, ut)
	assert.NoError(t, err)

	uts, err := utRepo.FindByUserMission(ctx, "um_1")
	assert.NoError(t, err)
	assert.Len(t, uts, 1)
	assert.True(t, uts[0].IsCompleted)
}

func TestGamificationRepository(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()

	ledgerRepo := postgres.NewPointLedgerRepository(pool)
	badgeRepo := postgres.NewBadgeRepository(pool)
	ubRepo := postgres.NewUserBadgeRepository(pool)
	creditRepo := postgres.NewCreditScoreRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	ctx := context.Background()

	u1 := &domain.User{UserID: "usr_1", Email: "u1@t.com"}
	require.NoError(t, userRepo.Create(ctx, u1))

	// Test PointLedgerRepository
	l1 := &domain.PointLedger{
		LedgerID:             "led_1",
		UserID:               "usr_1",
		SourceType:           domain.SourceMissionBase,
		SourceID:             "m_1",
		Delta:                100,
		LifetimeBalanceAfter: 100,
		PhaseBalanceAfter:    100,
		CreatedAt:            time.Now().Truncate(time.Second),
	}
	assert.NoError(t, ledgerRepo.Insert(ctx, l1))

	l2 := &domain.PointLedger{
		LedgerID:             "led_2",
		UserID:               "usr_1",
		SourceType:           domain.SourceMissionBase,
		SourceID:             "m_2",
		Delta:                50,
		LifetimeBalanceAfter: 150,
		PhaseBalanceAfter:    150,
		CreatedAt:            time.Now().Truncate(time.Second),
	}
	assert.NoError(t, ledgerRepo.InsertBatch(ctx, []domain.PointLedger{*l2}))

	// Test BadgeRepository and UserBadgeRepository
	_, err := pool.Exec(ctx, `INSERT INTO badges (badge_id, title, description, trigger_type, icon_url) VALUES ('badge_1', 'Badge 1', 'Desc', 'Speed', 'icon_1')`)
	require.NoError(t, err)

	badges, err := badgeRepo.FindByTriggerType(ctx, domain.TriggerSpeed)
	assert.NoError(t, err)
	assert.Len(t, badges, 1)

	eligibles, err := badgeRepo.FindEligible(ctx, "usr_1", domain.TriggerSpeed)
	assert.NoError(t, err)
	assert.Len(t, eligibles, 1)

	ub := &domain.UserBadge{
		UserBadgeID: "ub_1",
		UserID:      "usr_1",
		BadgeID:     "badge_1",
		SourceID:    "m_1",
		EarnedAt:    time.Now().Truncate(time.Second),
	}
	assert.NoError(t, ubRepo.Insert(ctx, ub))

	eligibles, err = badgeRepo.FindEligible(ctx, "usr_1", domain.TriggerSpeed)
	assert.NoError(t, err)
	assert.Len(t, eligibles, 0) // no longer eligible since they already earned it

	ubs, err := ubRepo.FindByUser(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Len(t, ubs, 1)

	// Test CreditScoreRepository
	c1 := &domain.CreditScore{
		CreditID:     "cred_1",
		UserID:       "usr_1",
		CurrentScore: 800,
		LastUpdated:  time.Now().Truncate(time.Second),
	}
	assert.NoError(t, creditRepo.Create(ctx, c1))

	foundCredit, err := creditRepo.FindByUserID(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, 800, foundCredit.CurrentScore)

	assert.NoError(t, creditRepo.Decrement(ctx, "usr_1", 50))
	foundCredit, err = creditRepo.FindByUserID(ctx, "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, 750, foundCredit.CurrentScore)

	_, err = creditRepo.FindByUserID(ctx, "usr_nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestRepositories_ErrNotFound(t *testing.T) {
	pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Repositories
	userRepo := postgres.NewUserRepository(pool)
	profileRepo := postgres.NewProfileRepository(pool)
	friendRepo := postgres.NewFriendshipRepository(pool)
	jobRepo := postgres.NewJobRepository(pool)
	cartRepo := postgres.NewUserCartRepository(pool)
	ratingRepo := postgres.NewJobOverallRatingRepository(pool)
	missionRepo := postgres.NewMissionRepository(pool)
	umRepo := postgres.NewUserMissionRepository(pool)
	phaseRepo := postgres.NewJourneyPhaseRepository(pool)
	historyRepo := postgres.NewUserPhaseHistoryRepository(pool)
	expenseRepo := postgres.NewExpenseRepository(pool)
	splitRepo := postgres.NewExpenseSplitRepository(pool)

	// Tests
	_, err := userRepo.FindByID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = userRepo.FindByEmail(ctx, "nonexistent@t.com")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = profileRepo.FindByUserID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = friendRepo.FindByID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = friendRepo.FindByCanonicalPair(ctx, "nonexistent1", "nonexistent2")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = jobRepo.FindByID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = cartRepo.FindByUserAndJob(ctx, "nonexistent", "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = cartRepo.FindByID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = ratingRepo.FindByJobID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = missionRepo.FindByID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = umRepo.FindByID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = phaseRepo.FindByID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = phaseRepo.FindByNumber(ctx, 9999)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = historyRepo.FindByUserAndPhase(ctx, "nonexistent", "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = expenseRepo.FindByID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = splitRepo.FindByID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
