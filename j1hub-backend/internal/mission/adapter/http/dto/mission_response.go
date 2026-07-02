package dto

import (
	missiondomain "github.com/parada3456/wat_project-backend/internal/mission/domain"
	missionusecase "github.com/parada3456/wat_project-backend/internal/mission/usecase"
)

type MissionDetailResponse struct {
	Mission     *missiondomain.Mission     `json:"mission"`
	UserMission *missiondomain.UserMission `json:"user_mission"`
	Tasks       []missiondomain.Task       `json:"tasks"`
	UserTasks   []missiondomain.UserTask   `json:"user_tasks"`
}

func NewMissionDetailResponse(detail *missionusecase.MissionDetailResponse) *MissionDetailResponse {
	tasks := detail.Tasks
	if tasks == nil {
		tasks = []missiondomain.Task{}
	}
	userTasks := detail.UserTasks
	if userTasks == nil {
		userTasks = []missiondomain.UserTask{}
	}
	return &MissionDetailResponse{
		Mission:     &detail.Mission,
		UserMission: &detail.UserMission,
		Tasks:       tasks,
		UserTasks:   userTasks,
	}
}
