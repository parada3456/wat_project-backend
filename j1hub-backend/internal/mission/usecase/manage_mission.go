package missionusecase

import (
	"context"
	"log"

	missiondomain "github.com/j1hub/backend/internal/mission/domain"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
)

type MissionUseCase struct {
	missionRepo port.MissionRepository
	umRepo      port.UserMissionRepository
	taskRepo    port.TaskRepository
	utRepo      port.UserTaskRepository
	userRepo    port.UserRepository
}

func NewMissionUseCase(
	missionRepo port.MissionRepository,
	umRepo port.UserMissionRepository,
	taskRepo port.TaskRepository,
	utRepo port.UserTaskRepository,
	userRepo port.UserRepository,
) *MissionUseCase {
	log.Println("debugprint: entering NewMissionUseCase")
	return &MissionUseCase{
		missionRepo: missionRepo,
		umRepo:      umRepo,
		taskRepo:    taskRepo,
		utRepo:      utRepo,
		userRepo:    userRepo,
	}
}

type MissionDetailResponse struct {
	Mission     missiondomain.Mission     `json:"mission"`
	UserMission missiondomain.UserMission `json:"user_mission"`
	Tasks       []missiondomain.Task      `json:"tasks"`
	UserTasks   []missiondomain.UserTask  `json:"user_tasks"`
}

func (uc *MissionUseCase) ListAvailableMissions(ctx context.Context, userID string) ([]missiondomain.UserMission, error) {
	log.Println("debugprint: entering (*MissionUseCase).ListAvailableMissions")
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return uc.umRepo.FindByUserAndPhase(ctx, userID, user.CurrentPhaseID)
}

func (uc *MissionUseCase) ListStaticMissions(ctx context.Context, userID string) ([]missiondomain.Mission, error) {
	log.Println("debugprint: entering (*MissionUseCase).ListStaticMissions")
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return uc.missionRepo.FindByPhase(ctx, user.CurrentPhaseID)
}

func (uc *MissionUseCase) GetMissionDetail(ctx context.Context, userID, userMissionID string) (*MissionDetailResponse, error) {
	log.Println("debugprint: entering (*MissionUseCase).GetMissionDetail")
	um, err := uc.umRepo.FindByID(ctx, userMissionID)
	if err != nil {
		return nil, err
	}
	if um.UserID != userID {
		return nil, domain.ErrForbidden
	}

	mission, err := uc.missionRepo.FindByID(ctx, um.MissionID)
	if err != nil {
		return nil, err
	}

	tasks, err := uc.taskRepo.FindByMission(ctx, um.MissionID)
	if err != nil {
		return nil, err
	}

	userTasks, err := uc.utRepo.FindByUserMission(ctx, userMissionID)
	if err != nil {
		return nil, err
	}

	return &MissionDetailResponse{
		Mission:     *mission,
		UserMission: *um,
		Tasks:       tasks,
		UserTasks:   userTasks,
	}, nil
}

func (uc *MissionUseCase) ToggleTask(ctx context.Context, userID, userTaskID string, completed bool) error {
	log.
		// Need to check ownership and update user task
		Println("debugprint: entering (*MissionUseCase).ToggleTask")

	return nil
}
