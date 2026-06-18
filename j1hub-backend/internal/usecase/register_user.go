package usecase

import (
	"context"
	"time"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/pkg/timeutil"
	"github.com/j1hub/backend/pkg/uid"
	"github.com/jackc/pgx/v5"
)

type TxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type RegisterUserUseCase struct {
	pool        TxBeginner
	userRepo    port.UserRepository
	profileRepo port.ProfileRepository
	creditRepo  port.CreditScoreRepository
	phaseRepo   port.JourneyPhaseRepository
	historyRepo port.UserPhaseHistoryRepository
	missionRepo port.MissionRepository
	umRepo      port.UserMissionRepository
	hasher      port.PasswordHasher
	tokenIssuer port.TokenIssuer
	clock       timeutil.Clock
}

func NewRegisterUserUseCase(
	pool TxBeginner,
	userRepo port.UserRepository,
	profileRepo port.ProfileRepository,
	creditRepo port.CreditScoreRepository,
	phaseRepo port.JourneyPhaseRepository,
	historyRepo port.UserPhaseHistoryRepository,
	missionRepo port.MissionRepository,
	umRepo port.UserMissionRepository,
	hasher port.PasswordHasher,
	tokenIssuer port.TokenIssuer,
	clock timeutil.Clock,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		pool:        pool,
		userRepo:    userRepo,
		profileRepo: profileRepo,
		creditRepo:  creditRepo,
		phaseRepo:   phaseRepo,
		historyRepo: historyRepo,
		missionRepo: missionRepo,
		umRepo:      umRepo,
		hasher:      hasher,
		tokenIssuer: tokenIssuer,
		clock:       clock,
	}
}

type RegisterCommand struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func (uc *RegisterUserUseCase) Register(ctx context.Context, cmd RegisterCommand) (*domain.User, *port.TokenPair, error) {
	hash, err := uc.hasher.Hash(cmd.Password)
	if err != nil {
		return nil, nil, err
	}

	user := &domain.User{
		UserID:       uid.New("usr_"),
		Email:        cmd.Email,
		PasswordHash: hash,
		FirstName:    cmd.FirstName,
		LastName:     cmd.LastName,
		CreatedAt:    uc.clock.Now(),
		UpdatedAt:    uc.clock.Now(),
	}

	// Transactional insert
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	profile := &domain.Profile{
		ProfileID:       uid.New("prf_"),
		UserID:          user.UserID,
		RadarVisibility: domain.VisibilityShowAnonymous,
		UpdatedAt:       uc.clock.Now(),
	}
	if err := uc.profileRepo.Create(ctx, profile); err != nil {
		return nil, nil, err
	}

	credit := &domain.CreditScore{
		CreditID:     uid.New("crd_"),
		UserID:       user.UserID,
		CurrentScore: 100,
		LastUpdated:  uc.clock.Now(),
	}
	if err := uc.creditRepo.Create(ctx, credit); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	tokens, err := uc.tokenIssuer.Issue(user.UserID, false)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

type InitJourneyCommand struct {
	ArrivalDate  time.Time
	JobStartDate time.Time
}

func (uc *RegisterUserUseCase) InitializeJourney(ctx context.Context, userID string, cmd InitJourneyCommand) error {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.ArrivalDate = cmd.ArrivalDate
	user.JobStartDate = cmd.JobStartDate
	user.UpdatedAt = uc.clock.Now()

	phase, err := uc.phaseRepo.FindByNumber(ctx, 1)
	if err != nil {
		return err
	}

	user.CurrentPhaseID = phase.PhaseID

	missions, err := uc.missionRepo.FindByPhase(ctx, phase.PhaseID)
	if err != nil {
		return err
	}

	var userMissions []domain.UserMission
	for _, m := range missions {
		triggerDate := user.ArrivalDate
		if m.RelativeTriggerEvent == "job_start_date" {
			triggerDate = user.JobStartDate
		}

		um := domain.UserMission{
			UserMissionID:     uid.New("ums_"),
			UserID:            user.UserID,
			MissionID:         m.MissionID,
			Status:            domain.StatusNotStarted,
			CalculatedDueDate: m.CalculateDueDate(triggerDate),
			CreatedAt:         uc.clock.Now(),
			UpdatedAt:         uc.clock.Now(),
		}
		userMissions = append(userMissions, um)
	}

	// Transaction
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	history := &domain.UserPhaseHistory{
		HistoryID: uid.New("uph_"),
		UserID:    user.UserID,
		PhaseID:   phase.PhaseID,
		EnteredAt: uc.clock.Now(),
	}
	if err := uc.historyRepo.Insert(ctx, history); err != nil {
		return err
	}

	if err := uc.umRepo.BulkInsert(ctx, userMissions); err != nil {
		return err
	}

	return nil
}
