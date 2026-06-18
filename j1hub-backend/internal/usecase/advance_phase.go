package usecase

import (
	"context"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/pkg/timeutil"
	"github.com/j1hub/backend/pkg/uid"
)

type AdvancePhaseUseCase struct {
	userRepo    port.UserRepository
	umRepo      port.UserMissionRepository
	phaseRepo   port.JourneyPhaseRepository
	historyRepo port.UserPhaseHistoryRepository
	missionRepo port.MissionRepository
	notifier    port.NotifierPort
	clock       timeutil.Clock
}

func NewAdvancePhaseUseCase(
	userRepo port.UserRepository,
	umRepo port.UserMissionRepository,
	phaseRepo port.JourneyPhaseRepository,
	historyRepo port.UserPhaseHistoryRepository,
	missionRepo port.MissionRepository,
	notifier port.NotifierPort,
	clock timeutil.Clock,
) *AdvancePhaseUseCase {
	return &AdvancePhaseUseCase{
		userRepo:    userRepo,
		umRepo:      umRepo,
		phaseRepo:   phaseRepo,
		historyRepo: historyRepo,
		missionRepo: missionRepo,
		notifier:    notifier,
		clock:       clock,
	}
}

func (uc *AdvancePhaseUseCase) TryAdvancePhase(ctx context.Context, userID string) (bool, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	ums, err := uc.umRepo.FindByUserAndPhase(ctx, userID, user.CurrentPhaseID)
	if err != nil {
		return false, err
	}

	if !domain.CanAdvancePhase(ums) {
		return false, domain.ErrPhaseNotComplete
	}

	currentPhase, err := uc.phaseRepo.FindByID(ctx, user.CurrentPhaseID)
	if err != nil {
		return false, err
	}

	nextPhase, err := uc.phaseRepo.FindByNumber(ctx, currentPhase.PhaseNumber+1)
	if err != nil {
		return false, err
	}

	now := uc.clock.Now()

	// Snapshot history
	if err := uc.historyRepo.CompleteCurrentPhase(ctx, userID, user.CurrentPhasePoints, now); err != nil {
		return false, err
	}

	// Update user
	if err := uc.userRepo.SetPhase(ctx, userID, nextPhase.PhaseID); err != nil {
		return false, err
	}
	// Also reset phase points to 0
	// ...

	// Insert new missions
	missions, err := uc.missionRepo.FindByPhase(ctx, nextPhase.PhaseID)
	if err != nil {
		return false, err
	}

	var newUMs []domain.UserMission
	for _, m := range missions {
		triggerDate := user.ArrivalDate
		if m.RelativeTriggerEvent == "job_start_date" {
			triggerDate = user.JobStartDate
		}
		newUMs = append(newUMs, domain.UserMission{
			UserMissionID:     uid.New("ums_"),
			UserID:            userID,
			MissionID:         m.MissionID,
			Status:            domain.StatusNotStarted,
			CalculatedDueDate: m.CalculateDueDate(triggerDate),
			CreatedAt:         now,
			UpdatedAt:         now,
		})
	}
	if err := uc.umRepo.BulkInsert(ctx, newUMs); err != nil {
		return false, err
	}

	uc.notifier.Send(ctx, userID, "Phase unlocked!", "New missions await!")

	return true, nil
}
