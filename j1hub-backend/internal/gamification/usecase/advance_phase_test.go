package gamificationusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	missiondomain "github.com/parada3456/wat_project-backend/internal/mission/domain"
	userdomain "github.com/parada3456/wat_project-backend/internal/user/domain"

	"github.com/parada3456/wat_project-backend/internal/domain"
	gamificationusecase "github.com/parada3456/wat_project-backend/internal/gamification/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdvancePhaseUseCase_TryAdvancePhase_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	umRepo := new(MockUserMissionRepository)
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	missionRepo := new(MockMissionRepository)
	notifier := new(MockNotifierPort)

	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}

	uc := gamificationusecase.NewAdvancePhaseUseCase(userRepo, umRepo, phaseRepo, historyRepo, missionRepo, notifier, clock)

	ctx := context.Background()
	userID := "usr_123"

	mockUser := &userdomain.User{
		UserID:             userID,
		CurrentPhaseID:     "phase_1",
		CurrentPhasePoints: 300,
		ArrivalDate:        nowTime.Add(-10 * 24 * time.Hour),
		JobStartDate:       nowTime.Add(-5 * 24 * time.Hour),
	}

	mockUserMissions := []missiondomain.UserMission{
		{UserMissionID: "ums_1", Status: missiondomain.StatusCompleted},
	}

	mockCurrentPhase := &missiondomain.JourneyPhase{
		PhaseID:     "phase_1",
		PhaseNumber: 1,
	}

	mockNextPhase := &missiondomain.JourneyPhase{
		PhaseID:     "phase_2",
		PhaseNumber: 2,
	}

	mockNewMissions := []missiondomain.Mission{
		{
			MissionID:            "m_next_1",
			PhaseID:              "phase_2",
			BasePoints:           200,
			RelativeTriggerEvent: "arrival_date",
			RelativeDaysOffset:   15,
		},
	}

	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	umRepo.On("FindByUserAndPhase", ctx, userID, "phase_1").Return(mockUserMissions, nil)
	phaseRepo.On("FindByID", ctx, "phase_1").Return(mockCurrentPhase, nil)
	phaseRepo.On("FindByNumber", ctx, 2).Return(mockNextPhase, nil)
	historyRepo.On("CompleteCurrentPhase", ctx, userID, 300, nowTime).Return(nil)
	userRepo.On("SetPhase", ctx, userID, "phase_2").Return(nil)
	missionRepo.On("FindByPhase", ctx, "phase_2").Return(mockNewMissions, nil)

	umRepo.On("BulkInsert", ctx, mock.AnythingOfType("[]missiondomain.UserMission")).Return(nil).Run(func(args mock.Arguments) {
		ums := args.Get(1).([]missiondomain.UserMission)
		assert.Len(t, ums, 1)
		assert.Equal(t, userID, ums[0].UserID)
		assert.Equal(t, "m_next_1", ums[0].MissionID)
	})

	notifier.On("Send", ctx, userID, "Phase unlocked!", "New missions await!").Return(nil)

	resp, err := uc.TryAdvancePhase(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Transitioned)
}

func TestAdvancePhaseUseCase_TryAdvancePhase_IncompleteMissions(t *testing.T) {
	userRepo := new(MockUserRepository)
	umRepo := new(MockUserMissionRepository)
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	missionRepo := new(MockMissionRepository)
	notifier := new(MockNotifierPort)
	clock := &MockClock{}

	uc := gamificationusecase.NewAdvancePhaseUseCase(userRepo, umRepo, phaseRepo, historyRepo, missionRepo, notifier, clock)

	ctx := context.Background()
	userID := "usr_123"

	mockUser := &userdomain.User{
		UserID:         userID,
		CurrentPhaseID: "phase_1",
	}

	// Contains an incomplete mission
	mockUserMissions := []missiondomain.UserMission{
		{UserMissionID: "ums_1", Status: missiondomain.StatusCompleted},
		{UserMissionID: "ums_2", Status: missiondomain.StatusInProgress},
	}

	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	umRepo.On("FindByUserAndPhase", ctx, userID, "phase_1").Return(mockUserMissions, nil)

	resp, err := uc.TryAdvancePhase(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, domain.ErrPhaseNotComplete, err)
}

func TestAdvancePhaseUseCase_TryAdvancePhase_UserRepoError(t *testing.T) {
	userRepo := new(MockUserRepository)
	uc := gamificationusecase.NewAdvancePhaseUseCase(userRepo, nil, nil, nil, nil, nil, &MockClock{})

	ctx := context.Background()
	userRepo.On("FindByID", ctx, "usr_1").Return((*userdomain.User)(nil), errors.New("db error"))

	resp, err := uc.TryAdvancePhase(ctx, "usr_1")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAdvancePhaseUseCase_TryAdvancePhase_UserMissionsRepoError(t *testing.T) {
	userRepo := new(MockUserRepository)
	umRepo := new(MockUserMissionRepository)
	uc := gamificationusecase.NewAdvancePhaseUseCase(userRepo, umRepo, nil, nil, nil, nil, &MockClock{})

	ctx := context.Background()
	userRepo.On("FindByID", ctx, "usr_1").Return(&userdomain.User{CurrentPhaseID: "phase_1"}, nil)
	umRepo.On("FindByUserAndPhase", ctx, "usr_1", "phase_1").Return([]missiondomain.UserMission{}, errors.New("db error"))

	resp, err := uc.TryAdvancePhase(ctx, "usr_1")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAdvancePhaseUseCase_TryAdvancePhase_CurrentPhaseRepoError(t *testing.T) {
	userRepo := new(MockUserRepository)
	umRepo := new(MockUserMissionRepository)
	phaseRepo := new(MockJourneyPhaseRepository)
	uc := gamificationusecase.NewAdvancePhaseUseCase(userRepo, umRepo, phaseRepo, nil, nil, nil, &MockClock{})

	ctx := context.Background()
	userRepo.On("FindByID", ctx, "usr_1").Return(&userdomain.User{CurrentPhaseID: "phase_1"}, nil)
	umRepo.On("FindByUserAndPhase", ctx, "usr_1", "phase_1").Return([]missiondomain.UserMission{{Status: missiondomain.StatusCompleted}}, nil)
	phaseRepo.On("FindByID", ctx, "phase_1").Return((*missiondomain.JourneyPhase)(nil), errors.New("db error"))

	resp, err := uc.TryAdvancePhase(ctx, "usr_1")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAdvancePhaseUseCase_TryAdvancePhase_NextPhaseRepoError(t *testing.T) {
	userRepo := new(MockUserRepository)
	umRepo := new(MockUserMissionRepository)
	phaseRepo := new(MockJourneyPhaseRepository)
	uc := gamificationusecase.NewAdvancePhaseUseCase(userRepo, umRepo, phaseRepo, nil, nil, nil, &MockClock{})

	ctx := context.Background()
	userRepo.On("FindByID", ctx, "usr_1").Return(&userdomain.User{CurrentPhaseID: "phase_1"}, nil)
	umRepo.On("FindByUserAndPhase", ctx, "usr_1", "phase_1").Return([]missiondomain.UserMission{{Status: missiondomain.StatusCompleted}}, nil)
	phaseRepo.On("FindByID", ctx, "phase_1").Return(&missiondomain.JourneyPhase{PhaseNumber: 1}, nil)
	phaseRepo.On("FindByNumber", ctx, 2).Return((*missiondomain.JourneyPhase)(nil), errors.New("db error"))

	resp, err := uc.TryAdvancePhase(ctx, "usr_1")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAdvancePhaseUseCase_TryAdvancePhase_CompletePhaseError(t *testing.T) {
	userRepo := new(MockUserRepository)
	umRepo := new(MockUserMissionRepository)
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	nowTime := time.Now()
	clock := &MockClock{NowTime: nowTime}
	uc := gamificationusecase.NewAdvancePhaseUseCase(userRepo, umRepo, phaseRepo, historyRepo, nil, nil, clock)

	ctx := context.Background()
	userRepo.On("FindByID", ctx, "usr_1").Return(&userdomain.User{CurrentPhaseID: "phase_1", CurrentPhasePoints: 100}, nil)
	umRepo.On("FindByUserAndPhase", ctx, "usr_1", "phase_1").Return([]missiondomain.UserMission{{Status: missiondomain.StatusCompleted}}, nil)
	phaseRepo.On("FindByID", ctx, "phase_1").Return(&missiondomain.JourneyPhase{PhaseNumber: 1}, nil)
	phaseRepo.On("FindByNumber", ctx, 2).Return(&missiondomain.JourneyPhase{PhaseID: "phase_2"}, nil)
	historyRepo.On("CompleteCurrentPhase", ctx, "usr_1", 100, nowTime).Return(errors.New("db error"))

	resp, err := uc.TryAdvancePhase(ctx, "usr_1")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAdvancePhaseUseCase_TryAdvancePhase_SetPhaseError(t *testing.T) {
	userRepo := new(MockUserRepository)
	umRepo := new(MockUserMissionRepository)
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	nowTime := time.Now()
	clock := &MockClock{NowTime: nowTime}
	uc := gamificationusecase.NewAdvancePhaseUseCase(userRepo, umRepo, phaseRepo, historyRepo, nil, nil, clock)

	ctx := context.Background()
	userRepo.On("FindByID", ctx, "usr_1").Return(&userdomain.User{CurrentPhaseID: "phase_1", CurrentPhasePoints: 100}, nil)
	umRepo.On("FindByUserAndPhase", ctx, "usr_1", "phase_1").Return([]missiondomain.UserMission{{Status: missiondomain.StatusCompleted}}, nil)
	phaseRepo.On("FindByID", ctx, "phase_1").Return(&missiondomain.JourneyPhase{PhaseNumber: 1}, nil)
	phaseRepo.On("FindByNumber", ctx, 2).Return(&missiondomain.JourneyPhase{PhaseID: "phase_2"}, nil)
	historyRepo.On("CompleteCurrentPhase", ctx, "usr_1", 100, nowTime).Return(nil)
	userRepo.On("SetPhase", ctx, "usr_1", "phase_2").Return(errors.New("db error"))

	resp, err := uc.TryAdvancePhase(ctx, "usr_1")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAdvancePhaseUseCase_TryAdvancePhase_FindByPhaseError(t *testing.T) {
	userRepo := new(MockUserRepository)
	umRepo := new(MockUserMissionRepository)
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	missionRepo := new(MockMissionRepository)
	nowTime := time.Now()
	clock := &MockClock{NowTime: nowTime}
	uc := gamificationusecase.NewAdvancePhaseUseCase(userRepo, umRepo, phaseRepo, historyRepo, missionRepo, nil, clock)

	ctx := context.Background()
	userRepo.On("FindByID", ctx, "usr_1").Return(&userdomain.User{CurrentPhaseID: "phase_1", CurrentPhasePoints: 100}, nil)
	umRepo.On("FindByUserAndPhase", ctx, "usr_1", "phase_1").Return([]missiondomain.UserMission{{Status: missiondomain.StatusCompleted}}, nil)
	phaseRepo.On("FindByID", ctx, "phase_1").Return(&missiondomain.JourneyPhase{PhaseNumber: 1}, nil)
	phaseRepo.On("FindByNumber", ctx, 2).Return(&missiondomain.JourneyPhase{PhaseID: "phase_2"}, nil)
	historyRepo.On("CompleteCurrentPhase", ctx, "usr_1", 100, nowTime).Return(nil)
	userRepo.On("SetPhase", ctx, "usr_1", "phase_2").Return(nil)
	missionRepo.On("FindByPhase", ctx, "phase_2").Return([]missiondomain.Mission{}, errors.New("db error"))

	resp, err := uc.TryAdvancePhase(ctx, "usr_1")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAdvancePhaseUseCase_TryAdvancePhase_BulkInsertError(t *testing.T) {
	userRepo := new(MockUserRepository)
	umRepo := new(MockUserMissionRepository)
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	missionRepo := new(MockMissionRepository)
	nowTime := time.Now()
	clock := &MockClock{NowTime: nowTime}
	uc := gamificationusecase.NewAdvancePhaseUseCase(userRepo, umRepo, phaseRepo, historyRepo, missionRepo, nil, clock)

	ctx := context.Background()
	userRepo.On("FindByID", ctx, "usr_1").Return(&userdomain.User{CurrentPhaseID: "phase_1", CurrentPhasePoints: 100}, nil)
	umRepo.On("FindByUserAndPhase", ctx, "usr_1", "phase_1").Return([]missiondomain.UserMission{{Status: missiondomain.StatusCompleted}}, nil)
	phaseRepo.On("FindByID", ctx, "phase_1").Return(&missiondomain.JourneyPhase{PhaseNumber: 1}, nil)
	phaseRepo.On("FindByNumber", ctx, 2).Return(&missiondomain.JourneyPhase{PhaseID: "phase_2"}, nil)
	historyRepo.On("CompleteCurrentPhase", ctx, "usr_1", 100, nowTime).Return(nil)
	userRepo.On("SetPhase", ctx, "usr_1", "phase_2").Return(nil)
	missionRepo.On("FindByPhase", ctx, "phase_2").Return([]missiondomain.Mission{
		{MissionID: "m_next_1", PhaseID: "phase_2"},
	}, nil)
	umRepo.On("BulkInsert", ctx, mock.Anything).Return(errors.New("db error"))

	resp, err := uc.TryAdvancePhase(ctx, "usr_1")
	assert.Error(t, err)
	assert.Nil(t, resp)
}
