package dto

import (
	jobdomain "github.com/j1hub/backend/internal/job/domain"
)

type JobDetailResponse struct {
	Job     *jobdomain.JobPosting       `json:"job"`
	Housing []jobdomain.JobHousing      `json:"housing"`
	Rating  *jobdomain.JobOverallRating `json:"rating"`
}

func NewJobDetailResponse(job *jobdomain.JobPosting, housing []jobdomain.JobHousing, rating *jobdomain.JobOverallRating) *JobDetailResponse {
	return &JobDetailResponse{
		Job:     job,
		Housing: housing,
		Rating:  rating,
	}
}
