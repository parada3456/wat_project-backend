package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUserUseCase_InitializeJourney_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	missionRepo := new(MockMissionRepository)
	umRepo := new(MockUserMissionRepository)
	hasher := new(MockHasher)
	issuer := new(MockIssuer)
	
	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}

	uc := usecase.NewRegisterUserUseCase(
		nil, // pool is not used by InitializeJourney
		userRepo,
		profileRepo,
		creditRepo,
		phaseRepo,
		historyRepo,
		missionRepo,
		umRepo,
		hasher,
		issuer,
		clock,
	)

	ctx := context.Background()
	userID := "usr_123"
	arrivalDate := nowTime.Add(24 * time.Hour)
	jobStartDate := nowTime.Add(48 * time.Hour)

	mockUser := &domain.User{
		UserID:    userID,
		Email:     "user@example.com",
		CreatedAt: nowTime,
	}

	mockPhase := &domain.JourneyPhase{
		PhaseID:     "phase_1",
		PhaseNumber: 1,
	}

	mockMissions := []domain.Mission{
		{
			MissionID:            "mission_1",
			PhaseID:              "phase_1",
			BasePoints:           100,
			RelativeTriggerEvent: "arrival_date",
			RelativeDaysOffset:   5,
		},
		{
			MissionID:            "mission_2",
			PhaseID:              "phase_1",
			BasePoints:           150,
			RelativeTriggerEvent: "job_start_date",
			RelativeDaysOffset:   10,
		},
	}

	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	phaseRepo.On("FindByNumber", ctx, 1).Return(mockPhase, nil)
	missionRepo.On("FindByPhase", ctx, "phase_1").Return(mockMissions, nil)
	
	// Expect user object update with phase and dates
	userRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*domain.User)
		assert.Equal(t, "phase_1", u.CurrentPhaseID)
		assert.Equal(t, arrivalDate, u.ArrivalDate)
		assert.Equal(t, jobStartDate, u.JobStartDate)
	})

	historyRepo.On("Insert", ctx, mock.AnythingOfType("*domain.UserPhaseHistory")).Return(nil).Run(func(args mock.Arguments) {
		h := args.Get(1).(*domain.UserPhaseHistory)
		assert.Equal(t, userID, h.UserID)
		assert.Equal(t, "phase_1", h.PhaseID)
		assert.Equal(t, nowTime, h.EnteredAt)
	})

	umRepo.On("BulkInsert", ctx, mock.AnythingOfType("[]domain.UserMission")).Return(nil).Run(func(args mock.Arguments) {
		ums := args.Get(1).([]domain.UserMission)
		assert.Len(t, ums, 2)
		assert.Equal(t, userID, ums[0].UserID)
		assert.Equal(t, "mission_1", ums[0].MissionID)
		assert.Equal(t, domain.StatusNotStarted, ums[0].Status)
	})

	err := uc.InitializeJourney(ctx, userID, usecase.InitJourneyCommand{
		ArrivalDate:  arrivalDate,
		JobStartDate: jobStartDate,
	})

	assert.NoError(t, err)

	userRepo.AssertExpectations(t)
	phaseRepo.AssertExpectations(t)
	missionRepo.AssertExpectations(t)
	historyRepo.AssertExpectations(t)
	umRepo.AssertExpectations(t)
}

func TestRegisterUserUseCase_Register_Success(t *testing.T) {
	poolMock := new(MockTxBeginner)
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	hasher := new(MockHasher)
	issuer := new(MockIssuer)
	
	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}

	uc := usecase.NewRegisterUserUseCase(
		poolMock,
		userRepo,
		profileRepo,
		creditRepo,
		nil,
		nil,
		nil,
		nil,
		hasher,
		issuer,
		clock,
	)

	ctx := context.Background()
	cmd := usecase.RegisterCommand{
		Email:     "john@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	txMock := new(MockTx)

	hasher.On("Hash", "password123").Return("hashed_password", nil)
	poolMock.On("Begin", ctx).Return(txMock, nil)
	userRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*domain.User)
		assert.Equal(t, "john@example.com", u.Email)
		assert.Equal(t, "hashed_password", u.PasswordHash)
		assert.Equal(t, "John", u.FirstName)
		assert.Equal(t, "Doe", u.LastName)
	})
	profileRepo.On("Create", ctx, mock.AnythingOfType("*domain.Profile")).Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(1).(*domain.Profile)
		assert.Equal(t, domain.VisibilityShowAnonymous, p.RadarVisibility)
	})
	creditRepo.On("Create", ctx, mock.AnythingOfType("*domain.CreditScore")).Return(nil).Run(func(args mock.Arguments) {
		c := args.Get(1).(*domain.CreditScore)
		assert.Equal(t, 100, c.CurrentScore)
	})
	txMock.On("Commit", ctx).Return(nil)
	txMock.On("Rollback", ctx).Return(nil) // Defer fallback

	tokens := &port.TokenPair{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}
	issuer.On("Issue", mock.AnythingOfType("string"), false).Return(tokens, nil)

	user, tokPair, err := uc.Register(ctx, cmd)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, tokens, tokPair)

	poolMock.AssertExpectations(t)
	txMock.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	profileRepo.AssertExpectations(t)
	creditRepo.AssertExpectations(t)
}

func TestRegisterUserUseCase_Register_HashError(t *testing.T) {
	hasher := new(MockHasher)
	uc := usecase.NewRegisterUserUseCase(nil, nil, nil, nil, nil, nil, nil, nil, hasher, nil, &MockClock{})

	ctx := context.Background()
	cmd := usecase.RegisterCommand{Password: "pwd"}
	hasher.On("Hash", "pwd").Return("", errors.New("hash error"))

	user, tokens, err := uc.Register(ctx, cmd)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, tokens)
	assert.Equal(t, "hash error", err.Error())
}

func TestRegisterUserUseCase_Register_BeginTxError(t *testing.T) {
	poolMock := new(MockTxBeginner)
	hasher := new(MockHasher)
	uc := usecase.NewRegisterUserUseCase(poolMock, nil, nil, nil, nil, nil, nil, nil, hasher, nil, &MockClock{})

	ctx := context.Background()
	cmd := usecase.RegisterCommand{Password: "pwd"}
	hasher.On("Hash", "pwd").Return("hash", nil)
	poolMock.On("Begin", ctx).Return((*MockTx)(nil), errors.New("db error"))

	user, tokens, err := uc.Register(ctx, cmd)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, tokens)
	assert.Equal(t, "db error", err.Error())
}

func TestRegisterUserUseCase_Register_CreateUserError(t *testing.T) {
	poolMock := new(MockTxBeginner)
	hasher := new(MockHasher)
	userRepo := new(MockUserRepository)
	uc := usecase.NewRegisterUserUseCase(poolMock, userRepo, nil, nil, nil, nil, nil, nil, hasher, nil, &MockClock{})

	ctx := context.Background()
	cmd := usecase.RegisterCommand{Password: "pwd", Email: "john@test.com"}
	txMock := new(MockTx)

	hasher.On("Hash", "pwd").Return("hash", nil)
	poolMock.On("Begin", ctx).Return(txMock, nil)
	userRepo.On("Create", ctx, mock.Anything).Return(errors.New("insert error"))
	txMock.On("Rollback", ctx).Return(nil)

	user, tokens, err := uc.Register(ctx, cmd)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, tokens)
	assert.Equal(t, "insert error", err.Error())
}
