package missionusecase_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	gamificationusecase "github.com/j1hub/backend/internal/gamification/usecase"
	missionusecase "github.com/j1hub/backend/internal/mission/usecase"

	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCompleteMissionUseCase_SubmitProof_Success(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	missionRepo := new(MockMissionRepository)
	taskRepo := new(MockTaskRepository)
	utRepo := new(MockUserTaskRepository)
	userRepo := new(MockUserRepository)
	ledgerRepo := new(MockPointLedgerRepository)
	badgeRepo := new(MockBadgeRepository)
	ubRepo := new(MockUserBadgeRepository)
	storage := new(MockStoragePort)
	notifier := new(MockNotifierPort)

	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}

	cfg := &config.Config{}
	rewardEngine := gamificationusecase.NewRewardEngine(cfg, userRepo, umRepo)

	uc := missionusecase.NewCompleteMissionUseCase(
		umRepo, missionRepo, taskRepo, utRepo, userRepo, ledgerRepo,
		badgeRepo, ubRepo, storage, notifier, rewardEngine, clock,
	)

	ctx := context.Background()
	userID := "usr_123"
	userMissionID := "ums_1"
	fileContent := []byte("fake_image_bytes")
	fileReader := bytes.NewReader(fileContent)
	contentType := "image/png"

	mockUM := &missiondomain.UserMission{
		UserMissionID: userMissionID,
		UserID:        userID,
		Status:        missiondomain.StatusInProgress,
	}

	umRepo.On("FindByID", ctx, userMissionID).Return(mockUM, nil)
	storage.On("UploadFile", ctx, "proofs", userMissionID, fileReader, contentType).Return("https://supabase.com/proof.png", nil)
	umRepo.On("UpdateStatus", ctx, userMissionID, missiondomain.StatusPendingVerification).Return(nil)

	err := uc.SubmitProof(ctx, userID, userMissionID, fileReader, contentType)

	assert.NoError(t, err)
	assert.Equal(t, "https://supabase.com/proof.png", mockUM.ProofURL)
	assert.Equal(t, nowTime, *mockUM.ProofSubmittedAt)
	assert.Equal(t, missiondomain.StatusPendingVerification, mockUM.Status)
}

func TestCompleteMissionUseCase_VerifyMission_Approve(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	missionRepo := new(MockMissionRepository)
	taskRepo := new(MockTaskRepository)
	utRepo := new(MockUserTaskRepository)
	userRepo := new(MockUserRepository)
	ledgerRepo := new(MockPointLedgerRepository)
	badgeRepo := new(MockBadgeRepository)
	ubRepo := new(MockUserBadgeRepository)
	storage := new(MockStoragePort)
	notifier := new(MockNotifierPort)

	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}

	cfg := &config.Config{
		Reward: config.RewardConfig{
			SpeedBonus7dPct: 10,
		},
	}
	rewardEngine := gamificationusecase.NewRewardEngine(cfg, userRepo, umRepo)

	uc := missionusecase.NewCompleteMissionUseCase(
		umRepo, missionRepo, taskRepo, utRepo, userRepo, ledgerRepo,
		badgeRepo, ubRepo, storage, notifier, rewardEngine, clock,
	)

	ctx := context.Background()
	adminID := "adm_999"
	userMissionID := "ums_1"
	userID := "usr_123"
	missionID := "m_1"

	submitTime := nowTime.Add(-2 * time.Hour)
	mockUM := &missiondomain.UserMission{
		UserMissionID:     userMissionID,
		UserID:            userID,
		MissionID:         missionID,
		Status:            missiondomain.StatusPendingVerification,
		ProofSubmittedAt:  &submitTime,
		CalculatedDueDate: nowTime.Add(24 * time.Hour),
	}

	mockUser := &userdomain.User{
		UserID:        userID,
		MissionStreak: 2,
	}

	mockMission := &missiondomain.Mission{
		MissionID:  missionID,
		BasePoints: 100,
	}

	umRepo.On("FindByID", ctx, userMissionID).Return(mockUM, nil)
	umRepo.On("UpdateVerification", ctx, userMissionID, nowTime, adminID).Return(nil)
	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	missionRepo.On("FindByID", ctx, missionID).Return(mockMission, nil)

	umRepo.On("UpdateReward", ctx, userMissionID, mock.AnythingOfType("*gamificationdomain.PointReward"), nowTime).Return(nil)
	umRepo.On("UpdateStatus", ctx, userMissionID, missiondomain.StatusCompleted).Return(nil)
	userRepo.On("IncrementPoints", ctx, userID, 100, 100).Return(nil)

	ledgerRepo.On("Insert", ctx, mock.Anything).Return(nil)
	notifier.On("Send", ctx, userID, "Mission complete!", "You earned points!").Return(nil)

	err := uc.VerifyMission(ctx, adminID, userMissionID, true)

	assert.NoError(t, err)
}

func TestCompleteMissionUseCase_VerifyMission_Reject(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	missionRepo := new(MockMissionRepository)
	taskRepo := new(MockTaskRepository)
	utRepo := new(MockUserTaskRepository)
	userRepo := new(MockUserRepository)
	ledgerRepo := new(MockPointLedgerRepository)
	badgeRepo := new(MockBadgeRepository)
	ubRepo := new(MockUserBadgeRepository)
	storage := new(MockStoragePort)
	notifier := new(MockNotifierPort)
	clock := &MockClock{}

	uc := missionusecase.NewCompleteMissionUseCase(
		umRepo, missionRepo, taskRepo, utRepo, userRepo, ledgerRepo,
		badgeRepo, ubRepo, storage, notifier, nil, clock,
	)

	ctx := context.Background()
	adminID := "adm_999"
	userMissionID := "ums_1"

	mockUM := &missiondomain.UserMission{
		UserMissionID: userMissionID,
		UserID:        "usr_123",
		Status:        missiondomain.StatusPendingVerification,
	}

	umRepo.On("FindByID", ctx, userMissionID).Return(mockUM, nil)
	umRepo.On("UpdateStatus", ctx, userMissionID, missiondomain.StatusInProgress).Return(nil)

	err := uc.VerifyMission(ctx, adminID, userMissionID, false)

	assert.NoError(t, err)
}

func TestCompleteMissionUseCase_SubmitProof_FindByID_Error(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	ctx := context.Background()
	umRepo.On("FindByID", ctx, "ums_1").Return((*missiondomain.UserMission)(nil), errors.New("db error"))

	err := uc.SubmitProof(ctx, "usr_123", "ums_1", nil, "image/png")
	assert.Error(t, err)
}

func TestCompleteMissionUseCase_SubmitProof_Forbidden(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	ctx := context.Background()
	mockUM := &missiondomain.UserMission{UserID: "usr_456"}
	umRepo.On("FindByID", ctx, "ums_1").Return(mockUM, nil)

	err := uc.SubmitProof(ctx, "usr_123", "ums_1", nil, "image/png")
	assert.Equal(t, domain.ErrForbidden, err)
}

func TestCompleteMissionUseCase_SubmitProof_AlreadyCompleted(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	ctx := context.Background()
	mockUM := &missiondomain.UserMission{UserID: "usr_123", Status: missiondomain.StatusCompleted}
	umRepo.On("FindByID", ctx, "ums_1").Return(mockUM, nil)

	err := uc.SubmitProof(ctx, "usr_123", "ums_1", nil, "image/png")
	assert.Equal(t, domain.ErrAlreadyCompleted, err)
}

func TestCompleteMissionUseCase_SubmitProof_UploadError(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	storage := new(MockStoragePort)
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, nil, nil, nil, nil, nil, nil, nil, storage, nil, nil, nil)
	ctx := context.Background()
	mockUM := &missiondomain.UserMission{UserID: "usr_123", Status: missiondomain.StatusInProgress}
	umRepo.On("FindByID", ctx, "ums_1").Return(mockUM, nil)
	storage.On("UploadFile", ctx, "proofs", "ums_1", mock.Anything, "image/png").Return("", errors.New("upload fail"))

	err := uc.SubmitProof(ctx, "usr_123", "ums_1", nil, "image/png")
	assert.Error(t, err)
}

func TestCompleteMissionUseCase_SubmitProof_UpdateStatusError(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	storage := new(MockStoragePort)
	clock := &MockClock{NowTime: time.Now()}
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, nil, nil, nil, nil, nil, nil, nil, storage, nil, nil, clock)
	ctx := context.Background()
	mockUM := &missiondomain.UserMission{UserID: "usr_123", Status: missiondomain.StatusInProgress}
	umRepo.On("FindByID", ctx, "ums_1").Return(mockUM, nil)
	storage.On("UploadFile", ctx, "proofs", "ums_1", mock.Anything, "image/png").Return("http://url", nil)
	umRepo.On("UpdateStatus", ctx, "ums_1", missiondomain.StatusPendingVerification).Return(errors.New("update fail"))

	err := uc.SubmitProof(ctx, "usr_123", "ums_1", nil, "image/png")
	assert.Error(t, err)
}

func TestCompleteMissionUseCase_VerifyMission_FindByID_Error(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	ctx := context.Background()
	umRepo.On("FindByID", ctx, "ums_1").Return((*missiondomain.UserMission)(nil), errors.New("db error"))

	err := uc.VerifyMission(ctx, "adm_1", "ums_1", true)
	assert.Error(t, err)
}

func TestCompleteMissionUseCase_VerifyMission_Reject_UpdateStatusError(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	ctx := context.Background()
	mockUM := &missiondomain.UserMission{Status: missiondomain.StatusPendingVerification}
	umRepo.On("FindByID", ctx, "ums_1").Return(mockUM, nil)
	umRepo.On("UpdateStatus", ctx, "ums_1", missiondomain.StatusInProgress).Return(errors.New("db error"))

	err := uc.VerifyMission(ctx, "adm_1", "ums_1", false)
	assert.Error(t, err)
}

func TestCompleteMissionUseCase_VerifyMission_UpdateVerificationError(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	now := time.Now()
	clock := &MockClock{NowTime: now}
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, clock)
	ctx := context.Background()
	mockUM := &missiondomain.UserMission{Status: missiondomain.StatusPendingVerification}
	umRepo.On("FindByID", ctx, "ums_1").Return(mockUM, nil)
	umRepo.On("UpdateVerification", ctx, "ums_1", now, "adm_1").Return(errors.New("db error"))

	err := uc.VerifyMission(ctx, "adm_1", "ums_1", true)
	assert.Error(t, err)
}

func TestCompleteMissionUseCase_VerifyMission_UserFindByID_Error(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	userRepo := new(MockUserRepository)
	now := time.Now()
	clock := &MockClock{NowTime: now}
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, nil, nil, nil, userRepo, nil, nil, nil, nil, nil, nil, clock)
	ctx := context.Background()
	mockUM := &missiondomain.UserMission{Status: missiondomain.StatusPendingVerification, UserID: "usr_1"}
	umRepo.On("FindByID", ctx, "ums_1").Return(mockUM, nil)
	umRepo.On("UpdateVerification", ctx, "ums_1", now, "adm_1").Return(nil)
	userRepo.On("FindByID", ctx, "usr_1").Return((*userdomain.User)(nil), errors.New("db error"))

	err := uc.VerifyMission(ctx, "adm_1", "ums_1", true)
	assert.Error(t, err)
}

func TestCompleteMissionUseCase_VerifyMission_MissionFindByID_Error(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	userRepo := new(MockUserRepository)
	missionRepo := new(MockMissionRepository)
	now := time.Now()
	clock := &MockClock{NowTime: now}
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, missionRepo, nil, nil, userRepo, nil, nil, nil, nil, nil, nil, clock)
	ctx := context.Background()
	mockUM := &missiondomain.UserMission{Status: missiondomain.StatusPendingVerification, UserID: "usr_1", MissionID: "m_1"}
	umRepo.On("FindByID", ctx, "ums_1").Return(mockUM, nil)
	umRepo.On("UpdateVerification", ctx, "ums_1", now, "adm_1").Return(nil)
	userRepo.On("FindByID", ctx, "usr_1").Return(&userdomain.User{}, nil)
	missionRepo.On("FindByID", ctx, "m_1").Return((*missiondomain.Mission)(nil), errors.New("db error"))

	err := uc.VerifyMission(ctx, "adm_1", "ums_1", true)
	assert.Error(t, err)
}

func TestCompleteMissionUseCase_VerifyMission_UpdateRewardError(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	userRepo := new(MockUserRepository)
	missionRepo := new(MockMissionRepository)
	now := time.Now()
	clock := &MockClock{NowTime: now}
	rewardEngine := gamificationusecase.NewRewardEngine(&config.Config{}, userRepo, umRepo)
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, missionRepo, nil, nil, userRepo, nil, nil, nil, nil, nil, rewardEngine, clock)
	ctx := context.Background()
	mockUM := &missiondomain.UserMission{Status: missiondomain.StatusPendingVerification, UserID: "usr_1", MissionID: "m_1"}
	umRepo.On("FindByID", ctx, "ums_1").Return(mockUM, nil)
	umRepo.On("UpdateVerification", ctx, "ums_1", now, "adm_1").Return(nil)
	userRepo.On("FindByID", ctx, "usr_1").Return(&userdomain.User{}, nil)
	missionRepo.On("FindByID", ctx, "m_1").Return(&missiondomain.Mission{}, nil)
	umRepo.On("UpdateReward", ctx, "ums_1", mock.Anything, now).Return(errors.New("db error"))

	err := uc.VerifyMission(ctx, "adm_1", "ums_1", true)
	assert.Error(t, err)
}

func TestCompleteMissionUseCase_VerifyMission_UpdateStatusError(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	userRepo := new(MockUserRepository)
	missionRepo := new(MockMissionRepository)
	now := time.Now()
	clock := &MockClock{NowTime: now}
	rewardEngine := gamificationusecase.NewRewardEngine(&config.Config{}, userRepo, umRepo)
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, missionRepo, nil, nil, userRepo, nil, nil, nil, nil, nil, rewardEngine, clock)
	ctx := context.Background()
	mockUM := &missiondomain.UserMission{Status: missiondomain.StatusPendingVerification, UserID: "usr_1", MissionID: "m_1"}
	umRepo.On("FindByID", ctx, "ums_1").Return(mockUM, nil)
	umRepo.On("UpdateVerification", ctx, "ums_1", now, "adm_1").Return(nil)
	userRepo.On("FindByID", ctx, "usr_1").Return(&userdomain.User{}, nil)
	missionRepo.On("FindByID", ctx, "m_1").Return(&missiondomain.Mission{}, nil)
	umRepo.On("UpdateReward", ctx, "ums_1", mock.Anything, now).Return(nil)
	umRepo.On("UpdateStatus", ctx, "ums_1", missiondomain.StatusCompleted).Return(errors.New("db error"))

	err := uc.VerifyMission(ctx, "adm_1", "ums_1", true)
	assert.Error(t, err)
}

func TestCompleteMissionUseCase_VerifyMission_IncrementPointsError(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	userRepo := new(MockUserRepository)
	missionRepo := new(MockMissionRepository)
	now := time.Now()
	clock := &MockClock{NowTime: now}
	rewardEngine := gamificationusecase.NewRewardEngine(&config.Config{}, userRepo, umRepo)
	uc := missionusecase.NewCompleteMissionUseCase(umRepo, missionRepo, nil, nil, userRepo, nil, nil, nil, nil, nil, rewardEngine, clock)
	ctx := context.Background()
	mockUM := &missiondomain.UserMission{Status: missiondomain.StatusPendingVerification, UserID: "usr_1", MissionID: "m_1"}
	umRepo.On("FindByID", ctx, "ums_1").Return(mockUM, nil)
	umRepo.On("UpdateVerification", ctx, "ums_1", now, "adm_1").Return(nil)
	userRepo.On("FindByID", ctx, "usr_1").Return(&userdomain.User{UserID: "usr_1"}, nil)
	missionRepo.On("FindByID", ctx, "m_1").Return(&missiondomain.Mission{}, nil)
	umRepo.On("UpdateReward", ctx, "ums_1", mock.Anything, now).Return(nil)
	umRepo.On("UpdateStatus", ctx, "ums_1", missiondomain.StatusCompleted).Return(nil)
	userRepo.On("IncrementPoints", ctx, "usr_1", mock.Anything, mock.Anything).Return(errors.New("db error"))

	err := uc.VerifyMission(ctx, "adm_1", "ums_1", true)
	assert.Error(t, err)
}
