package missionusecase

import (
	"context"
	"log"
	"time"

	missiondomain "github.com/j1hub/backend/internal/mission/domain"

	"github.com/j1hub/backend/internal/domain"
	port "github.com/j1hub/backend/internal/mission/port"
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

func (uc *MissionUseCase) ListAvailableMissions(ctx context.Context, userID string, ids []string) ([]missiondomain.UserMission, error) {
	log.Println("debugprint: entering (*MissionUseCase).ListAvailableMissions")
	if len(ids) > 0 {
		return uc.umRepo.FindByIDs(ctx, ids)
	}
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		print("error ")
		return nil, err
	}
	if user.CurrentPhaseID == "" {
		return []missiondomain.UserMission{}, nil
	}

	return uc.umRepo.FindByUserAndPhase(ctx, userID, user.CurrentPhaseID)
}

func (uc *MissionUseCase) ListStaticMissions(ctx context.Context, userID string, ids []string) ([]missiondomain.Mission, error) {
	log.Println("debugprint: entering (*MissionUseCase).ListStaticMissions")
	if len(ids) > 0 {
		return uc.missionRepo.FindByIDs(ctx, ids)
	}
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.CurrentPhaseID == "" {
		return []missiondomain.Mission{}, nil
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
	log.Println("debugprint: entering (*MissionUseCase).ToggleTask")

	ut, err := uc.utRepo.FindByID(ctx, userTaskID)
	if err != nil {
		return err
	}

	if ut.UserID != userID {
		return domain.ErrForbidden
	}

	ut.IsCompleted = completed
	if completed {
		now := time.Now()
		ut.CompletedAt = &now
	} else {
		ut.CompletedAt = nil
	}
	ut.UpdatedAt = time.Now()

	return uc.utRepo.Upsert(ctx, ut)
}

func (uc *MissionUseCase) ListTasks(ctx context.Context, ids []string) ([]missiondomain.Task, error) {
	log.Println("debugprint: entering (*MissionUseCase).ListTasks")
	if len(ids) > 0 {
		return uc.taskRepo.FindByIDs(ctx, ids)
	}
	return uc.taskRepo.ListAll(ctx)
}

func (uc *MissionUseCase) ListUserTasks(ctx context.Context, ids []string) ([]missiondomain.UserTask, error) {
	log.Println("debugprint: entering (*MissionUseCase).ListUserTasks")
	if len(ids) > 0 {
		return uc.utRepo.FindByIDs(ctx, ids)
	}
	return uc.utRepo.ListAll(ctx)
}
