package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScrapeJobsUseCase_Run(t *testing.T) {
	jobRepo := &MockJobPostingRepository{}
	housingRepo := &MockJobHousingRepository{}

	// Expect Upsert calls for scraped jobs/housings
	jobRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil).Maybe()
	housingRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil).Maybe()

	uc := usecase.NewScrapeJobsUseCase(jobRepo, housingRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := uc.Run(ctx)
	assert.NoError(t, err)
}
