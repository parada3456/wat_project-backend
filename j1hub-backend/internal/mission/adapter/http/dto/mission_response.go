package dto

import (
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	missionusecase "github.com/j1hub/backend/internal/mission/usecase"
)

type MissionDetailResponse struct {
	Mission     *missiondomain.Mission     `json:"mission"`
	UserMission *missiondomain.UserMission `json:"user_mission"`
	Tasks       []string                   `json:"tasks"`
	UserTasks   []string                   `json:"user_tasks"`
}

func NewMissionDetailResponse(detail *missionusecase.MissionDetailResponse) *MissionDetailResponse {
	taskIDs := make([]string, len(detail.Tasks))
	for i, t := range detail.Tasks {
		taskIDs[i] = t.TaskID
	}
	userTaskIDs := make([]string, len(detail.UserTasks))
	for i, ut := range detail.UserTasks {
		userTaskIDs[i] = ut.UserTaskID
	}
	return &MissionDetailResponse{
		Mission:     &detail.Mission,
		UserMission: &detail.UserMission,
		Tasks:       taskIDs,
		UserTasks:   userTaskIDs,
	}
}
