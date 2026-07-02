package scraper

import (
	"context"
	"testing"
	"time"

	jobdomain "github.com/parada3456/wat_project-backend/internal/job/domain"
	jobport "github.com/parada3456/wat_project-backend/internal/job/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJobPostingRepository struct {
	mock.Mock
	jobport.JobPostingRepository
}

func (m *MockJobPostingRepository) Upsert(ctx context.Context, job *jobdomain.JobPosting) error {
	return m.Called(ctx, job).Error(0)
}

type MockJobHousingRepository struct {
	mock.Mock
	jobport.JobHousingRepository
}

func (m *MockJobHousingRepository) Upsert(ctx context.Context, housing *jobdomain.JobHousing) error {
	return m.Called(ctx, housing).Error(0)
}

func TestScrapeJobsUseCase_Run(t *testing.T) {
	jobRepo := &MockJobPostingRepository{}
	housingRepo := &MockJobHousingRepository{}

	// Expect Upsert calls for scraped jobs/housings
	jobRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil).Maybe()
	housingRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil).Maybe()

	uc := NewScrapeJobsUseCase(jobRepo, housingRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := uc.Run(ctx)
	assert.NoError(t, err)
}
