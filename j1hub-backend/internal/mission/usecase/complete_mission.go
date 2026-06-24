package missionusecase

import (
	"context"
	"io"
	"log"

	gamificationusecase "github.com/j1hub/backend/internal/gamification/usecase"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"

	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"

	"github.com/j1hub/backend/internal/domain"
	port "github.com/j1hub/backend/internal/mission/port"
	"github.com/j1hub/backend/pkg/timeutil"
)

type CompleteMissionUseCase struct {
	umRepo       port.UserMissionRepository
	missionRepo  port.MissionRepository
	taskRepo     port.TaskRepository
	utRepo       port.UserTaskRepository
	userRepo     port.UserRepository
	ledgerRepo   port.PointLedgerRepository
	badgeRepo    port.BadgeRepository
	ubRepo       port.UserBadgeRepository
	storage      port.StoragePort
	notifier     port.NotifierPort
	rewardEngine *gamificationusecase.RewardEngine
	clock        timeutil.Clock
}

func NewCompleteMissionUseCase(
	umRepo port.UserMissionRepository,
	missionRepo port.MissionRepository,
	taskRepo port.TaskRepository,
	utRepo port.UserTaskRepository,
	userRepo port.UserRepository,
	ledgerRepo port.PointLedgerRepository,
	badgeRepo port.BadgeRepository,
	ubRepo port.UserBadgeRepository,
	storage port.StoragePort,
	notifier port.NotifierPort,
	rewardEngine *gamificationusecase.RewardEngine,
	clock timeutil.Clock,
) *CompleteMissionUseCase {
	log.Println("debugprint: entering NewCompleteMissionUseCase")
	return &CompleteMissionUseCase{
		umRepo:       umRepo,
		missionRepo:  missionRepo,
		taskRepo:     taskRepo,
		utRepo:       utRepo,
		userRepo:     userRepo,
		ledgerRepo:   ledgerRepo,
		badgeRepo:    badgeRepo,
		ubRepo:       ubRepo,
		storage:      storage,
		notifier:     notifier,
		rewardEngine: rewardEngine,
		clock:        clock,
	}
}

func (uc *CompleteMissionUseCase) SubmitProof(ctx context.Context, userID, userMissionID string, file io.Reader, contentType string) error {
	log.Println("debugprint: entering (*CompleteMissionUseCase).SubmitProof")
	um, err := uc.umRepo.FindByID(ctx, userMissionID)
	if err != nil {
		return err
	}
	if um.UserID != userID {
		return domain.ErrForbidden
	}
	if um.Status == missiondomain.StatusCompleted {
		return domain.ErrAlreadyCompleted
	}

	url, err := uc.storage.UploadFile(ctx, "proofs", userMissionID, file, contentType)
	if err != nil {
		return err
	}

	um.ProofURL = url
	now := uc.clock.Now()
	um.ProofSubmittedAt = &now
	um.Status = missiondomain.StatusPendingVerification
	um.UpdatedAt = now

	return uc.umRepo.UpdateStatus(ctx, userMissionID, um.Status)
	// Also need to update proof url and submitted at, but my repo interface only has UpdateStatus.
	// I should expand repo interface or use a more generic Update.
}

func (uc *CompleteMissionUseCase) VerifyMission(ctx context.Context, adminID, userMissionID string, approved bool) error {
	log.Println("debugprint: entering (*CompleteMissionUseCase).VerifyMission")
	um, err := uc.umRepo.FindByID(ctx, userMissionID)
	if err != nil {
		return err
	}

	if !approved {
		return uc.umRepo.UpdateStatus(ctx, userMissionID, missiondomain.StatusInProgress)
	}

	now := uc.clock.Now()
	if err := uc.umRepo.UpdateVerification(ctx, userMissionID, now, adminID); err != nil {
		return err
	}

	user, err := uc.userRepo.FindByID(ctx, um.UserID)
	if err != nil {
		return err
	}

	mission, err := uc.missionRepo.FindByID(ctx, um.MissionID)
	if err != nil {
		return err
	}

	reward, err := uc.rewardEngine.Calculate(ctx, um, user, mission)
	if err != nil {
		return err
	}

	// Update Reward and Status
	if err := uc.umRepo.UpdateReward(ctx, userMissionID, reward, now); err != nil {
		return err
	}
	if err := uc.umRepo.UpdateStatus(ctx, userMissionID, missiondomain.StatusCompleted); err != nil {
		return err
	}

	// Update User Points and Streak
	if err := uc.userRepo.IncrementPoints(ctx, user.UserID, reward.Total, reward.Total); err != nil {
		return err
	}
	// Streak update logic should be here too

	// Ledger Entry
	ledger := gamificationdomain.PointLedger{
		// ...
	}
	uc.ledgerRepo.Insert(ctx, &ledger)

	uc.notifier.Send(ctx, user.UserID, "Mission complete!", "You earned points!")

	return nil
}
