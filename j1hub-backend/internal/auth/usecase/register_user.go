package authusecase

import (
	"context"
	"log"
	"time"

	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	"net/http"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/pkg/apperror"
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
	log.Println("debugprint: entering NewRegisterUserUseCase")
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

func (uc *RegisterUserUseCase) Register(ctx context.Context, cmd RegisterCommand) (*userdomain.User, *port.TokenPair, error) {
	log.Println("debugprint: entering (*RegisterUserUseCase).Register")
	hash, err := uc.hasher.Hash(cmd.Password)
	if err != nil {
		return nil, nil, &apperror.ProblemDetails{
			Type:   "https://example.com/probs/internal-error",
			Title:  "Password Hashing Failed",
			Status: http.StatusInternalServerError,
			Detail: "Could not hash the password.",
		}
	}

	user := &userdomain.User{
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
		return nil, nil, &apperror.ProblemDetails{
			Type:   "https://example.com/probs/internal-error",
			Title:  "Transaction Error",
			Status: http.StatusInternalServerError,
			Detail: "Could not begin database transaction.",
		}
	}
	defer tx.Rollback(ctx)

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, nil, &apperror.ProblemDetails{
			Type:   "https://example.com/probs/user-creation-failed",
			Title:  "User Creation Failed",
			Status: http.StatusInternalServerError,
			Detail: "Failed to create user in the database.",
		}
	}

	profile := &userdomain.Profile{
		ProfileID:       uid.New("prf_"),
		UserID:          user.UserID,
		RadarVisibility: domain.VisibilityShowAnonymous,
		UpdatedAt:       uc.clock.Now(),
	}
	if err := uc.profileRepo.Create(ctx, profile); err != nil {
		return nil, nil, &apperror.ProblemDetails{
			Type:   "https://example.com/probs/profile-creation-failed",
			Title:  "Profile Creation Failed",
			Status: http.StatusInternalServerError,
			Detail: "Failed to create user profile.",
		}
	}

	credit := &userdomain.CreditScore{
		CreditID:     uid.New("crd_"),
		UserID:       user.UserID,
		CurrentScore: 100,
		LastUpdated:  uc.clock.Now(),
	}
	if err := uc.creditRepo.Create(ctx, credit); err != nil {
		return nil, nil, &apperror.ProblemDetails{
			Type:   "https://example.com/probs/credit-creation-failed",
			Title:  "Credit Creation Failed",
			Status: http.StatusInternalServerError,
			Detail: "Failed to create credit score record.",
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, &apperror.ProblemDetails{
			Type:   "https://example.com/probs/transaction-commit-failed",
			Title:  "Transaction Commit Failed",
			Status: http.StatusInternalServerError,
			Detail: "Could not commit the database transaction.",
		}
	}

	tokens, err := uc.tokenIssuer.Issue(user.UserID, false)
	if err != nil {
		return nil, nil, &apperror.ProblemDetails{
			Type:   "https://example.com/probs/token-issue-failed",
			Title:  "Token Issue Failed",
			Status: http.StatusInternalServerError,
			Detail: "Failed to issue authentication tokens.",
		}
	}

	return user, tokens, nil
}

type InitJourneyCommand struct {
	ArrivalDate  time.Time
	JobStartDate time.Time
}

func (uc *RegisterUserUseCase) InitializeJourney(ctx context.Context, userID string, cmd InitJourneyCommand) error {
	log.Println("debugprint: entering (*RegisterUserUseCase).InitializeJourney")
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

	var userMissions []missiondomain.UserMission
	for _, m := range missions {
		triggerDate := user.ArrivalDate
		if m.RelativeTriggerEvent == "job_start_date" {
			triggerDate = user.JobStartDate
		}

		um := missiondomain.UserMission{
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

	history := &gamificationdomain.UserPhaseHistory{
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
