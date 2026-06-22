package dto

import "github.com/j1hub/backend/internal/domain"

type JobDetailResponse struct {
	Job     *domain.JobPosting        `json:"job"`
	Housing []domain.JobHousing       `json:"housing"`
	Rating  *domain.JobOverallRating  `json:"rating"`
}

func NewJobDetailResponse(job *domain.JobPosting, housing []domain.JobHousing, rating *domain.JobOverallRating) *JobDetailResponse {
	return &JobDetailResponse{
		Job:     job,
		Housing: housing,
		Rating:  rating,
	}
}
