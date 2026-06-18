package scheduler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/j1hub/backend/internal/infrastructure/scheduler"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSplitRepo struct {
	mock.Mock
	port.ExpenseSplitRepository
}

func (m *MockSplitRepo) FindOverdue(ctx context.Context) ([]domain.ExpenseSplit, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ExpenseSplit), args.Error(1)
}

type MockUserMissionRepo struct {
	mock.Mock
	port.UserMissionRepository
}

func (m *MockUserMissionRepo) FindOverdue(ctx context.Context) ([]domain.UserMission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserMission), args.Error(1)
}

type MockJobPostingRepo struct {
	mock.Mock
	port.JobPostingRepository
}

func (m *MockJobPostingRepo) Upsert(ctx context.Context, job *domain.JobPosting) error {
	return m.Called(ctx, job).Error(0)
}

type MockJobHousingRepo struct {
	mock.Mock
	port.JobHousingRepository
}

func (m *MockJobHousingRepo) Upsert(ctx context.Context, housing *domain.JobHousing) error {
	return m.Called(ctx, housing).Error(0)
}

func TestNewScheduler(t *testing.T) {
	cfg := &config.Config{
		CronOverdueExpense: "*/5 * * * *",
		CronOverdueMission: "*/5 * * * *",
		CronScraper:        "*/5 * * * *",
	}

	splitRepo := new(MockSplitRepo)
	umRepo := new(MockUserMissionRepo)
	jobRepo := new(MockJobPostingRepo)
	housingRepo := new(MockJobHousingRepo)

	expenseJob := usecase.NewOverdueExpenseJob(splitRepo, nil, nil, nil)
	missionJob := usecase.NewOverdueMissionJob(umRepo, nil, nil, nil)
	scrapeJob := usecase.NewScrapeJobsUseCase(jobRepo, housingRepo)

	cronInstance := scheduler.NewScheduler(cfg, expenseJob, missionJob, scrapeJob)
	assert.NotNil(t, cronInstance)

	splitRepo.On("FindOverdue", mock.Anything).Return(nil, errors.New("err")).Once()
	umRepo.On("FindOverdue", mock.Anything).Return(nil, errors.New("err")).Once()
	jobRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil).Maybe()
	housingRepo.On("Upsert", mock.Anything, mock.Anything).Return(nil).Maybe()

	entries := cronInstance.Entries()
	assert.Len(t, entries, 3)
	for _, entry := range entries {
		entry.Job.Run()
	}
}
