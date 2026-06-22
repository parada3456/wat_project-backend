package dto

type VerifyMissionRequest struct {
	Approved        bool    `json:"approved"`
	RejectionReason *string `json:"rejectionReason"`
}

type AdjustPointsRequest struct {
	PointsDelta int    `json:"pointsDelta"`
	Reason      string `json:"reason" validate:"required"`
}

