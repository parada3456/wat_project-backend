package dto

type CartReq struct {
	JobID string `json:"job_id" validate:"required"`
}
