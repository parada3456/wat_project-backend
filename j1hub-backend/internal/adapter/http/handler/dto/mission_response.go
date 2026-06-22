package dto

import (
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	"github.com/j1hub/backend/internal/usecase"
)

type MissionDetailResponse struct {
	UserMissionID string                          `json:"user_mission_id"`
	MissionID     string                          `json:"mission_id"`
	Status        missiondomain.UserMissionStatus `json:"status"`
	Tasks         []missiondomain.UserTask        `json:"tasks"`
}

func NewMissionDetailResponse(detail *usecase.MissionDetailResponse) *MissionDetailResponse {
	return &MissionDetailResponse{
		UserMissionID: detail.UserMission.UserMissionID,
		MissionID:     detail.UserMission.MissionID,
		Status:        detail.UserMission.Status,
		Tasks:         detail.UserTasks,
	}
}
