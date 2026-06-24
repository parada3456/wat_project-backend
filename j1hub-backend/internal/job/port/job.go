package port

import (
	"context"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
)

type JobPostingRepository interface {
	FindWithFilters(ctx context.Context, filters map[string]interface{}) ([]jobdomain.JobPosting, error)
	FindByID(ctx context.Context, id string) (*jobdomain.JobPosting, error)
	Upsert(ctx context.Context, job *jobdomain.JobPosting) error
	Delete(ctx context.Context, id string) error
}

type JobHousingRepository interface {
	FindByJobID(ctx context.Context, jobID string) ([]jobdomain.JobHousing, error)
	Upsert(ctx context.Context, housing *jobdomain.JobHousing) error
}

type JobOverallRatingRepository interface {
	FindByJobID(ctx context.Context, jobID string) (*jobdomain.JobOverallRating, error)
	Recalculate(ctx context.Context, jobID string) error
}

type JobReviewRepository interface {
	Insert(ctx context.Context, r *jobdomain.JobReview) error
	FindByJobID(ctx context.Context, jobID string) ([]jobdomain.JobReview, error)
}

type UserCartRepository interface {
	Insert(ctx context.Context, c *jobdomain.UserCart) error
	FindByUserAndJob(ctx context.Context, userID, jobID string) (*jobdomain.UserCart, error)
	FindByID(ctx context.Context, id string) (*jobdomain.UserCart, error)
	UpdateStatus(ctx context.Context, id string, status jobdomain.CartStatus) error
	FindByUser(ctx context.Context, userID string) ([]jobdomain.UserCart, error)
	Delete(ctx context.Context, id string) error
}
