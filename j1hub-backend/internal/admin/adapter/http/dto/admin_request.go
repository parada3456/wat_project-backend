package dto

type VerifyMissionRequest struct {
	Approved        bool    `json:"approved"`
	RejectionReason *string `json:"rejection_reason"`
}

type AdjustPointsRequest struct {
	PointsDelta int    `json:"points_delta"`
	Reason      string `json:"reason" validate:"required"`
}
