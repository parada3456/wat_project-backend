package usecase

import (
	"context"
	"log"
	"time"

	missiondomain "github.com/j1hub/backend/internal/mission/domain"

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
	log.Println("debugprint: entering NewAdvancePhaseUseCase")
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

type PhaseTransitionResponse struct {
	Transitioned    bool      `json:"transitioned"`
	PreviousPhaseID string    `json:"previousPhaseId"`
	NewPhaseID      string    `json:"newPhaseId"`
	PointsRewarded  int       `json:"pointsRewarded"`
	CompletedAt     time.Time `json:"completedAt"`
}

func (uc *AdvancePhaseUseCase) TryAdvancePhase(ctx context.Context, userID string) (*PhaseTransitionResponse, error) {
	log.Println("debugprint: entering (*AdvancePhaseUseCase).TryAdvancePhase")
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	ums, err := uc.umRepo.FindByUserAndPhase(ctx, userID, user.CurrentPhaseID)
	if err != nil {
		return nil, err
	}

	if !domain.CanAdvancePhase(ums) {
		return nil, domain.ErrPhaseNotComplete
	}

	currentPhase, err := uc.phaseRepo.FindByID(ctx, user.CurrentPhaseID)
	if err != nil {
		return nil, err
	}

	nextPhase, err := uc.phaseRepo.FindByNumber(ctx, currentPhase.PhaseNumber+1)
	if err != nil {
		return nil, err
	}

	now := uc.clock.Now()

	// Snapshot history
	if err := uc.historyRepo.CompleteCurrentPhase(ctx, userID, user.CurrentPhasePoints, now); err != nil {
		return nil, err
	}

	// Update user
	if err := uc.userRepo.SetPhase(ctx, userID, nextPhase.PhaseID); err != nil {
		return nil, err
	}

	// Insert new missions
	missions, err := uc.missionRepo.FindByPhase(ctx, nextPhase.PhaseID)
	if err != nil {
		return nil, err
	}

	var newUMs []missiondomain.UserMission
	for _, m := range missions {
		triggerDate := user.ArrivalDate
		if m.RelativeTriggerEvent == "job_start_date" {
			triggerDate = user.JobStartDate
		}
		newUMs = append(newUMs, missiondomain.UserMission{
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
		return nil, err
	}

	uc.notifier.Send(ctx, userID, "Phase unlocked!", "New missions await!")

	return &PhaseTransitionResponse{
		Transitioned:    true,
		PreviousPhaseID: currentPhase.PhaseID,
		NewPhaseID:      nextPhase.PhaseID,
		PointsRewarded:  200,
		CompletedAt:     now,
	}, nil
}
