package dto

import (
	jobdomain "github.com/parada3456/wat_project-backend/internal/job/domain"
)

type JobDetailResponse struct {
	Job     *jobdomain.JobPosting       `json:"job"`
	Housing *jobdomain.JobHousing       `json:"housing"`
	Rating  *jobdomain.JobOverallRating `json:"rating"`
}

func NewJobDetailResponse(job *jobdomain.JobPosting, housing []jobdomain.JobHousing, rating *jobdomain.JobOverallRating) *JobDetailResponse {
	var h *jobdomain.JobHousing
	if len(housing) > 0 {
		h = &housing[0]
	}
	return &JobDetailResponse{
		Job:     job,
		Housing: h,
		Rating:  rating,
	}
}
