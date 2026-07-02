package dto

import "time"

type VerifyMissionRequest struct {
	Approved        bool    `json:"approved"`
	RejectionReason *string `json:"rejection_reason"`
}

type AdjustPointsRequest struct {
	PointsDelta int    `json:"points_delta"`
	Reason      string `json:"reason" validate:"required"`
}

type CreateMissionRequest struct {
	PhaseID              string              `json:"phase_id" validate:"required"`
	Title                string              `json:"title" validate:"required"`
	Description          string              `json:"description"`
	Location             string              `json:"location"`
	BasePoints           int                 `json:"base_points"`
	IsMandatory          bool                `json:"is_mandatory"`
	VerificationType     string              `json:"verification_type"`
	DueDateType          string              `json:"due_date_type"`
	FixedDueDate         *time.Time          `json:"fixed_due_date"`
	RelativeTriggerEvent string              `json:"relative_trigger_event"`
	RelativeDaysOffset   int                 `json:"relative_days_offset"`
	Tasks                []CreateTaskRequest `json:"tasks"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}
