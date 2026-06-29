package adminusecase

import (
	"context"
	"fmt"
	"log"

	gamificationusecase "github.com/j1hub/backend/internal/gamification/usecase"

	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/j1hub/backend/internal/domain"
	port "github.com/j1hub/backend/internal/admin/port"
	"github.com/j1hub/backend/pkg/timeutil"
	"github.com/j1hub/backend/pkg/uid"
	"github.com/jackc/pgx/v5"
)

type adminUseCase struct {
	pool         port.TxBeginner
	adminRepo    port.AdminRepository
	userRepo     port.UserRepository
	profileRepo  port.ProfileRepository
	umRepo       port.UserMissionRepository
	missionRepo  port.MissionRepository
	ledgerRepo   port.PointLedgerRepository
	notifier     port.NotifierPort
	rewardEngine *gamificationusecase.RewardEngine
	clock        timeutil.Clock
}

func NewAdminUseCase(
	pool port.TxBeginner,
	adminRepo port.AdminRepository,
	userRepo port.UserRepository,
	profileRepo port.ProfileRepository,
	umRepo port.UserMissionRepository,
	missionRepo port.MissionRepository,
	ledgerRepo port.PointLedgerRepository,
	notifier port.NotifierPort,
	rewardEngine *gamificationusecase.RewardEngine,
	clock timeutil.Clock,
) port.AdminUseCase {
	log.Println("debugprint: entering NewAdminUseCase")
	return &adminUseCase{
		pool:         pool,
		adminRepo:    adminRepo,
		userRepo:     userRepo,
		profileRepo:  profileRepo,
		umRepo:       umRepo,
		missionRepo:  missionRepo,
		ledgerRepo:   ledgerRepo,
		notifier:     notifier,
		rewardEngine: rewardEngine,
		clock:        clock,
	}
}

func (u *adminUseCase) GetDashboardStats(ctx context.Context) (*port.AdminStats, error) {
	log.Println("debugprint: entering (*adminUseCase).GetDashboardStats")
	return u.adminRepo.GetStats(ctx)
}

func (u *adminUseCase) ListPendingVerifications(ctx context.Context, page, pageSize int) ([]missiondomain.UserMission, int, error) {
	log.Println("debugprint: entering (*adminUseCase).ListPendingVerifications")
	limit := pageSize
	offset := (page - 1) * pageSize
	return u.adminRepo.ListPendingVerifications(ctx, limit, offset)
}

func (u *adminUseCase) ListUsers(ctx context.Context, search string, page, pageSize int) ([]port.UserWithProfile, int, error) {
	log.Println("debugprint: entering (*adminUseCase).ListUsers")
	limit := pageSize
	offset := (page - 1) * pageSize
	return u.adminRepo.SearchUsers(ctx, search, limit, offset)
}

func (u *adminUseCase) GetUserDetail(ctx context.Context, id string) (*port.UserWithProfile, error) {
	log.Println("debugprint: entering (*adminUseCase).GetUserDetail")
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	profile, err := u.profileRepo.FindByUserID(ctx, id)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	res := &port.UserWithProfile{
		User: *user,
	}
	if profile != nil {
		res.Profile = *profile
	}
	return res, nil
}

func (u *adminUseCase) AdjustPoints(ctx context.Context, userID string, delta int, reason string) (*port.PointsAdjustmentResult, error) {
	log.Println("debugprint: entering (*adminUseCase).AdjustPoints")

	// Execute within database transaction
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Fetch user inside transaction
	var user userdomain.User
	var currentPhaseID *string
	err = tx.QueryRow(ctx, `
		SELECT user_id, total_lifetime_points, current_phase_points, current_phase_id
		FROM users WHERE user_id = $1 FOR UPDATE`, userID).Scan(&user.UserID, &user.TotalLifetimePoints, &user.CurrentPhasePoints, &currentPhaseID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	newLifetime := user.TotalLifetimePoints + delta
	newPhasePoints := user.CurrentPhasePoints + delta
	if newLifetime < 0 {
		newLifetime = 0
	}
	if newPhasePoints < 0 {
		newPhasePoints = 0
	}

	// Update user points
	_, err = tx.Exec(ctx, `
		UPDATE users 
		SET total_lifetime_points = $1, current_phase_points = $2, updated_at = NOW() 
		WHERE user_id = $3`, newLifetime, newPhasePoints, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user points: %w", err)
	}

	// Insert point ledger
	ledgerID := uid.New("ldg")
	_, err = tx.Exec(ctx, `
		INSERT INTO point_ledger (ledger_id, user_id, points, source_type, description, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())`,
		ledgerID, userID, delta, "Admin_Adjust", reason,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert point ledger: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	u.notifier.Send(ctx, userID, "Points adjusted by Admin", fmt.Sprintf("Your points balance was updated by %d", delta))

	return &port.PointsAdjustmentResult{
		UserID:               userID,
		LifetimeBalanceAfter: newLifetime,
		PhaseBalanceAfter:    newPhasePoints,
		LedgerID:             ledgerID,
	}, nil
}

func (u *adminUseCase) VerifyMission(ctx context.Context, adminID, userMissionID string, approved bool, rejectionReason *string) (*missiondomain.UserMission, error) {
	log.Println("debugprint: entering (*adminUseCase).VerifyMission")

	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Fetch user mission to verify
	var um missiondomain.UserMission
	var proofURL *string
	var verifiedBy *string
	query := `
		SELECT 
			user_mission_id, user_id, mission_id, status, calculated_due_date, 
			proof_url, proof_submitted_at, verified_at, verified_by, 
			base_points_earned, speed_bonus_points, streak_bonus_points, 
			first_completer_bonus_points, total_points_earned, rewarded_at, 
			created_at, updated_at 
		FROM user_missions 
		WHERE user_mission_id = $1 FOR UPDATE`

	err = tx.QueryRow(ctx, query, userMissionID).Scan(
		&um.UserMissionID, &um.UserID, &um.MissionID, &um.Status, &um.CalculatedDueDate,
		&proofURL, &um.ProofSubmittedAt, &um.VerifiedAt, &verifiedBy,
		&um.BasePointsEarned, &um.SpeedBonusPoints, &um.StreakBonusPoints,
		&um.FirstCompleterBonusPoints, &um.TotalPointsEarned, &um.RewardedAt,
		&um.CreatedAt, &um.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch user mission: %w", err)
	}

	// Idempotency: Reject calls if user mission is already Completed
	if um.Status == missiondomain.StatusCompleted {
		return nil, fmt.Errorf("user mission is already completed")
	}

	now := u.clock.Now()

	if !approved {
		// Rejection flow
		um.Status = missiondomain.StatusInProgress
		_, err = tx.Exec(ctx, `
			UPDATE user_missions 
			SET status = $1, updated_at = $2 
			WHERE user_mission_id = $3`, um.Status, now, userMissionID)
		if err != nil {
			return nil, fmt.Errorf("failed to reject user mission: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}

		msg := "Your mission submission was rejected."
		if rejectionReason != nil {
			msg = fmt.Sprintf("Your mission submission was rejected. Reason: %s", *rejectionReason)
		}
		u.notifier.Send(ctx, um.UserID, "Mission Submission Rejected", msg)

		return &um, nil
	}

	// Approval flow: Calculate points/streak
	var user userdomain.User
	var curPhaseID *string
	err = tx.QueryRow(ctx, `
		SELECT user_id, email, password_hash, current_phase_id, total_lifetime_points, current_phase_points, mission_streak
		FROM users WHERE user_id = $1 FOR UPDATE`, um.UserID).Scan(
		&user.UserID, &user.Email, &user.PasswordHash,
		&curPhaseID, &user.TotalLifetimePoints, &user.CurrentPhasePoints, &user.MissionStreak,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user for reward calculation: %w", err)
	}
	if curPhaseID != nil {
		user.CurrentPhaseID = *curPhaseID
	}

	var mission missiondomain.Mission
	var mDesc *string
	var mLoc *string
	var mTrigger *string
	var mOffset *int
	err = tx.QueryRow(ctx, `
		SELECT mission_id, phase_id, title, description, location, base_points, is_mandatory, verification_type, due_date_type, fixed_due_date, relative_trigger_event, relative_days_offset
		FROM missions WHERE mission_id = $1`, um.MissionID).Scan(
		&mission.MissionID, &mission.PhaseID, &mission.Title, &mDesc, &mLoc, &mission.BasePoints,
		&mission.IsMandatory, &mission.VerificationType, &mission.DueDateType, &mission.FixedDueDate,
		&mTrigger, &mOffset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mission details: %w", err)
	}

	reward, err := u.rewardEngine.Calculate(ctx, &um, &user, &mission)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate rewards: %w", err)
	}

	um.Status = missiondomain.StatusCompleted
	um.BasePointsEarned = reward.Base
	um.SpeedBonusPoints = reward.SpeedBonus
	um.StreakBonusPoints = reward.StreakBonus
	um.FirstCompleterBonusPoints = reward.FirstCompleterBonus
	um.TotalPointsEarned = reward.Total
	um.VerifiedAt = &now
	um.VerifiedBy = adminID
	um.RewardedAt = &now

	// Update user mission record
	_, err = tx.Exec(ctx, `
		UPDATE user_missions 
		SET status = $1, base_points_earned = $2, speed_bonus_points = $3, 
			streak_bonus_points = $4, first_completer_bonus_points = $5, 
			total_points_earned = $6, verified_at = $7, verified_by = $8, 
			rewarded_at = $9, updated_at = $10
		WHERE user_mission_id = $11`,
		um.Status, um.BasePointsEarned, um.SpeedBonusPoints,
		um.StreakBonusPoints, um.FirstCompleterBonusPoints,
		um.TotalPointsEarned, um.VerifiedAt, um.VerifiedBy,
		um.RewardedAt, now, userMissionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user mission state: %w", err)
	}

	// Increment User Points
	newLifetime := user.TotalLifetimePoints + reward.Total
	newPhasePoints := user.CurrentPhasePoints + reward.Total
	_, err = tx.Exec(ctx, `
		UPDATE users 
		SET total_lifetime_points = $1, current_phase_points = $2, updated_at = $3 
		WHERE user_id = $4`,
		newLifetime, newPhasePoints, now, um.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user points balance: %w", err)
	}

	// Insert point ledger record
	ledgerID := uid.New("ldg")
	ledgerQuery := `
		INSERT INTO point_ledger (ledger_id, user_id, points, source_type, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(ctx, ledgerQuery, ledgerID, um.UserID, reward.Total, "Mission_Base", "Completed Mission: "+mission.Title, now)
	if err != nil {
		return nil, fmt.Errorf("failed to insert points transaction log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	u.notifier.Send(ctx, um.UserID, "Mission verified and completed!", fmt.Sprintf("You earned %d total points!", reward.Total))

	return &um, nil
}
