package usecase_test

import (
	"context"
	"testing"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestMissionUseCase_ListAvailableMissions_Success(t *testing.T) {
	missionRepo := new(MockMissionRepository)
	umRepo := new(MockUserMissionRepository)
	taskRepo := new(MockTaskRepository)
	utRepo := new(MockUserTaskRepository)
	userRepo := new(MockUserRepository)

	uc := usecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)

	ctx := context.Background()
	userID := "usr_123"
	phaseID := "phase_1"

	mockUser := &domain.User{UserID: userID, CurrentPhaseID: phaseID}
	mockUserMissions := []domain.UserMission{
		{UserMissionID: "ums_1", UserID: userID, MissionID: "m_1"},
	}

	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	umRepo.On("FindByUserAndPhase", ctx, userID, phaseID).Return(mockUserMissions, nil)

	res, err := uc.ListAvailableMissions(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockUserMissions, res)
}

func TestMissionUseCase_GetMissionDetail_Success(t *testing.T) {
	missionRepo := new(MockMissionRepository)
	umRepo := new(MockUserMissionRepository)
	taskRepo := new(MockTaskRepository)
	utRepo := new(MockUserTaskRepository)
	userRepo := new(MockUserRepository)

	uc := usecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)

	ctx := context.Background()
	userID := "usr_123"
	userMissionID := "ums_1"
	missionID := "m_1"

	mockUM := &domain.UserMission{UserMissionID: userMissionID, UserID: userID, MissionID: missionID}
	mockMission := &domain.Mission{MissionID: missionID, Title: "Test Mission"}
	mockTasks := []domain.Task{{TaskID: "t_1", MissionID: missionID}}
	mockUserTasks := []domain.UserTask{{UserTaskID: "ut_1", UserID: userID, TaskID: "t_1"}}

	umRepo.On("FindByID", ctx, userMissionID).Return(mockUM, nil)
	missionRepo.On("FindByID", ctx, missionID).Return(mockMission, nil)
	taskRepo.On("FindByMission", ctx, missionID).Return(mockTasks, nil)
	utRepo.On("FindByUserMission", ctx, userMissionID).Return(mockUserTasks, nil)

	res, err := uc.GetMissionDetail(ctx, userID, userMissionID)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, *mockMission, res.Mission)
	assert.Equal(t, *mockUM, res.UserMission)
	assert.Equal(t, mockTasks, res.Tasks)
	assert.Equal(t, mockUserTasks, res.UserTasks)
}

func TestMissionUseCase_GetMissionDetail_Forbidden(t *testing.T) {
	missionRepo := new(MockMissionRepository)
	umRepo := new(MockUserMissionRepository)
	taskRepo := new(MockTaskRepository)
	utRepo := new(MockUserTaskRepository)
	userRepo := new(MockUserRepository)

	uc := usecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)

	ctx := context.Background()
	userID := "usr_123"
	userMissionID := "ums_1"

	// Belongs to usr_999
	mockUM := &domain.UserMission{UserMissionID: userMissionID, UserID: "usr_999", MissionID: "m_1"}

	umRepo.On("FindByID", ctx, userMissionID).Return(mockUM, nil)

	res, err := uc.GetMissionDetail(ctx, userID, userMissionID)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, domain.ErrForbidden, err)
}

func TestMissionUseCase_ToggleTask_Stub(t *testing.T) {
	uc := usecase.NewMissionUseCase(nil, nil, nil, nil, nil)
	err := uc.ToggleTask(context.Background(), "usr_123", "ut_1", true)
	assert.NoError(t, err)
}
