package dto

import (
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
)

type MissionDetailResponse struct {
	UserMissionID string             `json:"user_mission_id"`
	MissionID     string             `json:"mission_id"`
	Status        domain.UserMissionStatus `json:"status"`
	Tasks         []domain.UserTask  `json:"tasks"`
}

func NewMissionDetailResponse(detail *usecase.MissionDetailResponse) *MissionDetailResponse {
	return &MissionDetailResponse{
		UserMissionID: detail.UserMission.UserMissionID,
		MissionID:     detail.UserMission.MissionID,
		Status:        detail.UserMission.Status,
		Tasks:         detail.UserTasks,
	}
}
