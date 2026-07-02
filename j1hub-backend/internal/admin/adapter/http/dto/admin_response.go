package dto

import (
	"time"

	missiondomain "github.com/j1hub/backend/internal/mission/domain"
)

//	{
//	  "users": {
//	    "total": 1250,
//	    "active_this_month": 430,
//	    "new_this_week": 52
//	  },
//	  "missions": {
//	    "pending_verification": 18,
//	    "completed_today": 35
//	  },
//	  "expenses": {
//	    "overdue_splits": 7,
//	    "total_volume_usd": "48250.00"
//	  },
//	  "jobs": {
//	    "active_listings": 92,
//	    "total_reviews": 340
//	  }
//	}
type DashboardStatsResponse struct {
	// TODO: implement like json
}

type VerifyMissionResponse struct {
	UserMissionID string     `json:"user_mission_id"`
	Status        string     `json:"status"`
	VerifiedAt    *time.Time `json:"verified_at"`
	VerifiedBy    string     `json:"verified_by"`
}

func NewVerifyMissionResponse(um *missiondomain.UserMission) *VerifyMissionResponse {
	return &VerifyMissionResponse{
		UserMissionID: um.UserMissionID,
		Status:        string(um.Status),
		VerifiedAt:    um.VerifiedAt,
		VerifiedBy:    um.VerifiedBy,
	}
}
