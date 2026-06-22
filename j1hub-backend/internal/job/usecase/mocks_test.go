package usecase_test

import (
	"context"
	"time"

	jobdomain "github.com/j1hub/backend/internal/job/domain"
	"github.com/stretchr/testify/mock"
)

// MockJobPostingRepository
type MockJobPostingRepository struct{ mock.Mock }

func (m *MockJobPostingRepository) FindWithFilters(ctx context.Context, filters map[string]interface{}) ([]jobdomain.JobPosting, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jobdomain.JobPosting), args.Error(1)
}
func (m *MockJobPostingRepository) FindByID(ctx context.Context, id string) (*jobdomain.JobPosting, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jobdomain.JobPosting), args.Error(1)
}
func (m *MockJobPostingRepository) Upsert(ctx context.Context, job *jobdomain.JobPosting) error {
	return m.Called(ctx, job).Error(0)
}

// MockJobHousingRepository
type MockJobHousingRepository struct{ mock.Mock }

func (m *MockJobHousingRepository) FindByJobID(ctx context.Context, jobID string) ([]jobdomain.JobHousing, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jobdomain.JobHousing), args.Error(1)
}
func (m *MockJobHousingRepository) Upsert(ctx context.Context, housing *jobdomain.JobHousing) error {
	return m.Called(ctx, housing).Error(0)
}

// MockJobOverallRatingRepository
type MockJobOverallRatingRepository struct{ mock.Mock }

func (m *MockJobOverallRatingRepository) FindByJobID(ctx context.Context, jobID string) (*jobdomain.JobOverallRating, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jobdomain.JobOverallRating), args.Error(1)
}
func (m *MockJobOverallRatingRepository) Recalculate(ctx context.Context, jobID string) error {
	return m.Called(ctx, jobID).Error(0)
}

// MockJobReviewRepository
type MockJobReviewRepository struct{ mock.Mock }

func (m *MockJobReviewRepository) Insert(ctx context.Context, r *jobdomain.JobReview) error {
	return m.Called(ctx, r).Error(0)
}
func (m *MockJobReviewRepository) FindByJobID(ctx context.Context, jobID string) ([]jobdomain.JobReview, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jobdomain.JobReview), args.Error(1)
}

// MockUserCartRepository
type MockUserCartRepository struct{ mock.Mock }

func (m *MockUserCartRepository) Insert(ctx context.Context, c *jobdomain.UserCart) error {
	return m.Called(ctx, c).Error(0)
}
func (m *MockUserCartRepository) FindByUserAndJob(ctx context.Context, userID, jobID string) (*jobdomain.UserCart, error) {
	args := m.Called(ctx, userID, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jobdomain.UserCart), args.Error(1)
}
func (m *MockUserCartRepository) FindByID(ctx context.Context, id string) (*jobdomain.UserCart, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jobdomain.UserCart), args.Error(1)
}
func (m *MockUserCartRepository) UpdateStatus(ctx context.Context, id string, status jobdomain.CartStatus) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *MockUserCartRepository) FindByUser(ctx context.Context, userID string) ([]jobdomain.UserCart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jobdomain.UserCart), args.Error(1)
}
func (m *MockUserCartRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// MockClock
type MockClock struct {
	NowTime time.Time
}

func (m *MockClock) Now() time.Time {
	return m.NowTime
}
