package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestManageJobUseCase_AddToCart_Success(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	housingRepo := new(MockJobHousingRepository)
	ratingRepo := new(MockJobOverallRatingRepository)
	reviewRepo := new(MockJobReviewRepository)
	cartRepo := new(MockUserCartRepository)
	
	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}

	uc := usecase.NewManageJobUseCase(jobRepo, housingRepo, ratingRepo, reviewRepo, cartRepo, clock)

	ctx := context.Background()
	userID := "usr_123"
	jobID := "job_456"

	cartRepo.On("FindByUserAndJob", ctx, userID, jobID).Return((*domain.UserCart)(nil), domain.ErrNotFound)
	cartRepo.On("Insert", ctx, mock.AnythingOfType("*domain.UserCart")).Return(nil).Run(func(args mock.Arguments) {
		c := args.Get(1).(*domain.UserCart)
		assert.Equal(t, userID, c.UserID)
		assert.Equal(t, jobID, c.JobID)
		assert.Equal(t, domain.CartSaved, c.Status)
		assert.Equal(t, nowTime, c.AddedAt)
	})

	err := uc.AddToCart(ctx, userID, jobID)

	assert.NoError(t, err)
	cartRepo.AssertExpectations(t)
}

func TestManageJobUseCase_AddToCart_Conflict(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	housingRepo := new(MockJobHousingRepository)
	ratingRepo := new(MockJobOverallRatingRepository)
	reviewRepo := new(MockJobReviewRepository)
	cartRepo := new(MockUserCartRepository)
	clock := &MockClock{}

	uc := usecase.NewManageJobUseCase(jobRepo, housingRepo, ratingRepo, reviewRepo, cartRepo, clock)

	ctx := context.Background()
	userID := "usr_123"
	jobID := "job_456"

	cartRepo.On("FindByUserAndJob", ctx, userID, jobID).Return(&domain.UserCart{CartID: "crt_1"}, nil)

	err := uc.AddToCart(ctx, userID, jobID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrConflict, err)
}

func TestManageJobUseCase_WriteReview_Success(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	housingRepo := new(MockJobHousingRepository)
	ratingRepo := new(MockJobOverallRatingRepository)
	reviewRepo := new(MockJobReviewRepository)
	cartRepo := new(MockUserCartRepository)
	
	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}

	uc := usecase.NewManageJobUseCase(jobRepo, housingRepo, ratingRepo, reviewRepo, cartRepo, clock)

	ctx := context.Background()
	userID := "usr_123"
	jobID := "job_456"
	review := &domain.JobReview{
		RatingStars: 4.5,
		ReviewText:  "Excellent place",
	}

	reviewRepo.On("Insert", ctx, review).Return(nil).Run(func(args mock.Arguments) {
		r := args.Get(1).(*domain.JobReview)
		assert.Equal(t, userID, r.UserID)
		assert.Equal(t, jobID, r.JobID)
		assert.Equal(t, nowTime, r.CreatedAt)
	})

	ratingRepo.On("Recalculate", ctx, jobID).Return(nil)

	err := uc.WriteReview(ctx, userID, jobID, review)

	assert.NoError(t, err)
	reviewRepo.AssertExpectations(t)
	ratingRepo.AssertExpectations(t)
}

func TestManageJobUseCase_ListJobs_Success(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	uc := usecase.NewManageJobUseCase(jobRepo, nil, nil, nil, nil, &MockClock{})

	ctx := context.Background()
	filters := map[string]interface{}{"agency": "InterExchange"}
	mockJobs := []domain.JobPosting{{JobID: "job_1"}}

	jobRepo.On("FindWithFilters", ctx, filters).Return(mockJobs, nil)

	res, err := uc.ListJobs(ctx, filters)

	assert.NoError(t, err)
	assert.Equal(t, mockJobs, res)
}

func TestManageJobUseCase_GetJobDetail_Success(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	housingRepo := new(MockJobHousingRepository)
	ratingRepo := new(MockJobOverallRatingRepository)
	uc := usecase.NewManageJobUseCase(jobRepo, housingRepo, ratingRepo, nil, nil, &MockClock{})

	ctx := context.Background()
	jobID := "job_1"
	mockJob := &domain.JobPosting{JobID: jobID}
	mockHousing := []domain.JobHousing{{HousingID: "h_1", JobID: jobID}}
	mockRating := &domain.JobOverallRating{JobID: jobID, OverallRate: 4.2}

	jobRepo.On("FindByID", ctx, jobID).Return(mockJob, nil)
	housingRepo.On("FindByJobID", ctx, jobID).Return(mockHousing, nil)
	ratingRepo.On("FindByJobID", ctx, jobID).Return(mockRating, nil)

	job, housing, rating, err := uc.GetJobDetail(ctx, jobID)

	assert.NoError(t, err)
	assert.Equal(t, mockJob, job)
	assert.Equal(t, mockHousing, housing)
	assert.Equal(t, mockRating, rating)
}

func TestManageJobUseCase_ListCart_Stub(t *testing.T) {
	uc := usecase.NewManageJobUseCase(nil, nil, nil, nil, nil, &MockClock{})
	res, err := uc.ListCart(context.Background(), "usr_1")
	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestManageJobUseCase_RemoveFromCart_Stub(t *testing.T) {
	uc := usecase.NewManageJobUseCase(nil, nil, nil, nil, nil, &MockClock{})
	err := uc.RemoveFromCart(context.Background(), "usr_1", "crt_1")
	assert.NoError(t, err)
}
