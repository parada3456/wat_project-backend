package dto

import (
	"time"

	missiondomain "github.com/j1hub/backend/internal/mission/domain"
)

type DashboardStatsResponse struct {
	// Add fields as necessary
}

type VerifyMissionResponse struct {
	UserMissionID string     `json:"userMissionId"`
	Status        string     `json:"status"`
	VerifiedAt    *time.Time `json:"verifiedAt"`
	VerifiedBy    string     `json:"verifiedBy"`
}

func NewVerifyMissionResponse(um *missiondomain.UserMission) *VerifyMissionResponse {
	return &VerifyMissionResponse{
		UserMissionID: um.UserMissionID,
		Status:        string(um.Status),
		VerifiedAt:    um.VerifiedAt,
		VerifiedBy:    um.VerifiedBy,
	}
}
