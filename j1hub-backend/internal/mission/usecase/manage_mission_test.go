package missionusecase_test

// import (
// 	"context"
// 	"testing"

// 	missionusecase "github.com/j1hub/backend/internal/mission/usecase"

// 	missiondomain "github.com/j1hub/backend/internal/mission/domain"
// 	userdomain "github.com/j1hub/backend/internal/user/domain"

// 	"github.com/j1hub/backend/internal/domain"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestMissionUseCase_ListAvailableMissions_Success(t *testing.T) {
// 	missionRepo := new(MockMissionRepository)
// 	umRepo := new(MockUserMissionRepository)
// 	taskRepo := new(MockTaskRepository)
// 	utRepo := new(MockUserTaskRepository)
// 	userRepo := new(MockUserRepository)

// 	uc := missionusecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)

// 	ctx := context.Background()
// 	userID := "usr_123"
// 	phaseID := "phase_1"

// 	mockUser := &userdomain.User{UserID: userID, CurrentPhaseID: phaseID}
// 	mockUserMissions := []missiondomain.UserMission{
// 		{UserMissionID: "ums_1", UserID: userID, MissionID: "m_1"},
// 	}

// 	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
// 	umRepo.On("FindByUserAndPhase", ctx, userID, phaseID).Return(mockUserMissions, nil)

// 	res, err := uc.ListAvailableMissions(ctx, userID)

// 	assert.NoError(t, err)
// 	assert.Equal(t, mockUserMissions, res)
// }

// func TestMissionUseCase_GetMissionDetail_Success(t *testing.T) {
// 	missionRepo := new(MockMissionRepository)
// 	umRepo := new(MockUserMissionRepository)
// 	taskRepo := new(MockTaskRepository)
// 	utRepo := new(MockUserTaskRepository)
// 	userRepo := new(MockUserRepository)

// 	uc := missionusecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)

// 	ctx := context.Background()
// 	userID := "usr_123"
// 	userMissionID := "ums_1"
// 	missionID := "m_1"

// 	mockUM := &missiondomain.UserMission{UserMissionID: userMissionID, UserID: userID, MissionID: missionID}
// 	mockMission := &missiondomain.Mission{MissionID: missionID, Title: "Test Mission"}
// 	mockTasks := []missiondomain.Task{{TaskID: "t_1", MissionID: missionID}}
// 	mockUserTasks := []missiondomain.UserTask{{UserTaskID: "ut_1", UserID: userID, TaskID: "t_1"}}

// 	umRepo.On("FindByID", ctx, userMissionID).Return(mockUM, nil)
// 	missionRepo.On("FindByID", ctx, missionID).Return(mockMission, nil)
// 	taskRepo.On("FindByMission", ctx, missionID).Return(mockTasks, nil)
// 	utRepo.On("FindByUserMission", ctx, userMissionID).Return(mockUserTasks, nil)

// 	res, err := uc.GetMissionDetail(ctx, userID, userMissionID)

// 	assert.NoError(t, err)
// 	assert.NotNil(t, res)
// 	assert.Equal(t, *mockMission, res.Mission)
// 	assert.Equal(t, *mockUM, res.UserMission)
// 	assert.Equal(t, mockTasks, res.Tasks)
// 	assert.Equal(t, mockUserTasks, res.UserTasks)
// }

// func TestMissionUseCase_GetMissionDetail_Forbidden(t *testing.T) {
// 	missionRepo := new(MockMissionRepository)
// 	umRepo := new(MockUserMissionRepository)
// 	taskRepo := new(MockTaskRepository)
// 	utRepo := new(MockUserTaskRepository)
// 	userRepo := new(MockUserRepository)

// 	uc := missionusecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)

// 	ctx := context.Background()
// 	userID := "usr_123"
// 	userMissionID := "ums_1"

// 	// Belongs to usr_999
// 	mockUM := &missiondomain.UserMission{UserMissionID: userMissionID, UserID: "usr_999", MissionID: "m_1"}

// 	umRepo.On("FindByID", ctx, userMissionID).Return(mockUM, nil)

// 	res, err := uc.GetMissionDetail(ctx, userID, userMissionID)

// 	assert.Error(t, err)
// 	assert.Nil(t, res)
// 	assert.Equal(t, domain.ErrForbidden, err)
// }

// func TestMissionUseCase_ToggleTask_Success(t *testing.T) {
// 	missionRepo := new(MockMissionRepository)
// 	umRepo := new(MockUserMissionRepository)
// 	taskRepo := new(MockTaskRepository)
// 	utRepo := new(MockUserTaskRepository)
// 	userRepo := new(MockUserRepository)

// 	uc := missionusecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)

// 	ctx := context.Background()
// 	userID := "usr_123"
// 	userTaskID := "ut_1"

// 	mockUT := &missiondomain.UserTask{
// 		UserTaskID:    userTaskID,
// 		UserID:        userID,
// 		TaskID:        "t_1",
// 		UserMissionID: "ums_1",
// 		IsCompleted:   false,
// 	}

// 	utRepo.On("FindByID", ctx, userTaskID).Return(mockUT, nil)
// 	utRepo.On("Upsert", ctx, mock.MatchedBy(func(ut *missiondomain.UserTask) bool {
// 		return ut.UserTaskID == userTaskID && ut.IsCompleted == true && ut.CompletedAt != nil
// 	})).Return(nil)

// 	err := uc.ToggleTask(ctx, userID, userTaskID, true)
// 	assert.NoError(t, err)
// }

// func TestMissionUseCase_ToggleTask_Forbidden(t *testing.T) {
// 	missionRepo := new(MockMissionRepository)
// 	umRepo := new(MockUserMissionRepository)
// 	taskRepo := new(MockTaskRepository)
// 	utRepo := new(MockUserTaskRepository)
// 	userRepo := new(MockUserRepository)

// 	uc := missionusecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)

// 	ctx := context.Background()
// 	userID := "usr_123"
// 	userTaskID := "ut_1"

// 	// Belongs to usr_999
// 	mockUT := &missiondomain.UserTask{
// 		UserTaskID:    userTaskID,
// 		UserID:        "usr_999",
// 		TaskID:        "t_1",
// 		UserMissionID: "ums_1",
// 		IsCompleted:   false,
// 	}

// 	utRepo.On("FindByID", ctx, userTaskID).Return(mockUT, nil)

// 	err := uc.ToggleTask(ctx, userID, userTaskID, true)
// 	assert.Equal(t, domain.ErrForbidden, err)
// }

// func TestMissionUseCase_ToggleTask_NotFound(t *testing.T) {
// 	missionRepo := new(MockMissionRepository)
// 	umRepo := new(MockUserMissionRepository)
// 	taskRepo := new(MockTaskRepository)
// 	utRepo := new(MockUserTaskRepository)
// 	userRepo := new(MockUserRepository)

// 	uc := missionusecase.NewMissionUseCase(missionRepo, umRepo, taskRepo, utRepo, userRepo)

// 	ctx := context.Background()
// 	userID := "usr_123"
// 	userTaskID := "ut_1"

// 	utRepo.On("FindByID", ctx, userTaskID).Return(nil, domain.ErrNotFound)

// 	err := uc.ToggleTask(ctx, userID, userTaskID, true)
// 	assert.Equal(t, domain.ErrNotFound, err)
// }
