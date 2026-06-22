package scraper

import (
	"context"

	jobdomain "github.com/j1hub/backend/internal/job/domain"
)

type JobSource interface {
	GetJobLinks(ctx context.Context, listURL string) ([]string, error)
	GetJobDetails(ctx context.Context, detailURL string) (*jobdomain.JobPosting, error)
}
