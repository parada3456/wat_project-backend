package scheduler

import (
	"context"
	"errors"
	"testing"

	expensedomain "github.com/j1hub/backend/internal/expense/domain"
	expenseport "github.com/j1hub/backend/internal/expense/port"
	"github.com/j1hub/backend/internal/infrastructure/config"
	scraper "github.com/j1hub/backend/internal/infrastructure/outbound/scraper"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
	jobport "github.com/j1hub/backend/internal/job/port"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	missionport "github.com/j1hub/backend/internal/mission/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSplitRepo struct {
	mock.Mock
	expenseport.ExpenseSplitRepository
}

func (m *MockSplitRepo) FindOverdue(ctx context.Context) ([]expensedomain.ExpenseSplit, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]expensedomain.ExpenseSplit), args.Error(1)
}

type MockUserMissionRepo struct {
	mock.Mock
	missionport.UserMissionRepository
}

func (m *MockUserMissionRepo) FindOverdue(ctx context.Context) ([]missiondomain.UserMission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.UserMission), args.Error(1)
}

type MockJobPostingRepo struct {
	mock.Mock
	jobport.JobPostingRepository
}

func (m *MockJobPostingRepo) Upsert(ctx context.Context, job *jobdomain.JobPosting) error {
	return m.Called(ctx, job).Error(0)
}

type MockJobHousingRepo struct {
	mock.Mock
	jobport.JobHousingRepository
}

func (m *MockJobHousingRepo) Upsert(ctx context.Context, housing *jobdomain.JobHousing) error {
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

	expenseJob := NewOverdueExpenseJob(splitRepo, nil, nil, nil)
	missionJob := NewOverdueMissionJob(umRepo, nil, nil, nil)
	scrapeJob := scraper.NewScrapeJobsUseCase(jobRepo, housingRepo)

	cronInstance := NewScheduler(cfg, expenseJob, missionJob, scrapeJob)
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
