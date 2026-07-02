package jobusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/j1hub/backend/internal/domain"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
	jobusecase "github.com/j1hub/backend/internal/job/usecase"
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

	uc := jobusecase.NewManageJobUseCase(jobRepo, housingRepo, ratingRepo, reviewRepo, cartRepo, clock)

	ctx := context.Background()
	userID := "usr_123"
	jobID := "job_456"

	cartRepo.On("FindByUserAndJob", ctx, userID, jobID).Return((*jobdomain.UserCart)(nil), domain.ErrNotFound)
	cartRepo.On("Insert", ctx, mock.AnythingOfType("*jobdomain.UserCart")).Return(nil).Run(func(args mock.Arguments) {
		c := args.Get(1).(*jobdomain.UserCart)
		assert.Equal(t, userID, c.UserID)
		assert.Equal(t, jobID, c.JobID)
		assert.Equal(t, jobdomain.CartSaved, c.Status)
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

	uc := jobusecase.NewManageJobUseCase(jobRepo, housingRepo, ratingRepo, reviewRepo, cartRepo, clock)

	ctx := context.Background()
	userID := "usr_123"
	jobID := "job_456"

	cartRepo.On("FindByUserAndJob", ctx, userID, jobID).Return(&jobdomain.UserCart{CartID: "crt_1"}, nil)

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

	uc := jobusecase.NewManageJobUseCase(jobRepo, housingRepo, ratingRepo, reviewRepo, cartRepo, clock)

	ctx := context.Background()
	userID := "usr_123"
	jobID := "job_456"
	review := &jobdomain.JobReview{
		RatingStars: 4.5,
		ReviewText:  "Excellent place",
	}

	reviewRepo.On("Insert", ctx, review).Return(nil).Run(func(args mock.Arguments) {
		r := args.Get(1).(*jobdomain.JobReview)
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
	uc := jobusecase.NewManageJobUseCase(jobRepo, nil, nil, nil, nil, &MockClock{})

	ctx := context.Background()
	filters := map[string]interface{}{"agency": "InterExchange"}
	mockJobs := []jobdomain.JobPosting{{JobID: "job_1"}}

	jobRepo.On("FindWithFilters", ctx, filters, 10, 0).Return(mockJobs, 1, nil)

	res, totalCount, err := uc.ListJobs(ctx, filters, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, mockJobs, res)
	assert.Equal(t, 1, totalCount)
}

func TestManageJobUseCase_GetJobDetail_Success(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	housingRepo := new(MockJobHousingRepository)
	ratingRepo := new(MockJobOverallRatingRepository)
	uc := jobusecase.NewManageJobUseCase(jobRepo, housingRepo, ratingRepo, nil, nil, &MockClock{})

	ctx := context.Background()
	jobID := "job_1"
	mockJob := &jobdomain.JobPosting{JobID: jobID}
	mockHousing := []jobdomain.JobHousing{{HousingID: "h_1", JobID: jobID}}
	mockRating := &jobdomain.JobOverallRating{JobID: jobID, OverallRate: 4.2}

	jobRepo.On("FindByID", ctx, jobID).Return(mockJob, nil)
	housingRepo.On("FindByJobID", ctx, jobID).Return(mockHousing, nil)
	ratingRepo.On("FindByJobID", ctx, jobID).Return(mockRating, nil)

	job, housing, rating, err := uc.GetJobDetail(ctx, jobID)

	assert.NoError(t, err)
	assert.Equal(t, mockJob, job)
	assert.Equal(t, mockHousing, housing)
	assert.Equal(t, mockRating, rating)
}

func TestManageJobUseCase_ListCart_Success(t *testing.T) {
	cartRepo := new(MockUserCartRepository)
	uc := jobusecase.NewManageJobUseCase(nil, nil, nil, nil, cartRepo, &MockClock{})
	ctx := context.Background()
	userID := "usr_1"
	mockCart := []jobdomain.UserCart{{CartID: "crt_1", UserID: userID}}

	cartRepo.On("FindByUser", ctx, userID).Return(mockCart, nil)

	res, err := uc.ListCart(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, mockCart, res)
}

func TestManageJobUseCase_RemoveFromCart_Success(t *testing.T) {
	cartRepo := new(MockUserCartRepository)
	uc := jobusecase.NewManageJobUseCase(nil, nil, nil, nil, cartRepo, &MockClock{})
	ctx := context.Background()
	userID := "usr_1"
	cartID := "crt_1"

	cartRepo.On("FindByID", ctx, cartID).Return(&jobdomain.UserCart{CartID: cartID, UserID: userID}, nil)
	cartRepo.On("Delete", ctx, cartID).Return(nil)

	err := uc.RemoveFromCart(ctx, userID, cartID)
	assert.NoError(t, err)
}

func TestManageJobUseCase_RemoveFromCart_Forbidden(t *testing.T) {
	cartRepo := new(MockUserCartRepository)
	uc := jobusecase.NewManageJobUseCase(nil, nil, nil, nil, cartRepo, &MockClock{})
	ctx := context.Background()
	userID := "usr_1"
	cartID := "crt_1"

	cartRepo.On("FindByID", ctx, cartID).Return(&jobdomain.UserCart{CartID: cartID, UserID: "other_user"}, nil)

	err := uc.RemoveFromCart(ctx, userID, cartID)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)
}

func TestManageJobUseCase_ListReviews_Success(t *testing.T) {
	reviewRepo := new(MockJobReviewRepository)
	uc := jobusecase.NewManageJobUseCase(nil, nil, nil, reviewRepo, nil, &MockClock{})
	ctx := context.Background()
	jobID := "job_1"
	mockReviews := []jobdomain.JobReview{{ReviewID: "rev_1", JobID: jobID}}

	reviewRepo.On("FindByJobID", ctx, jobID).Return(mockReviews, nil)

	res, err := uc.ListReviews(ctx, jobID)
	assert.NoError(t, err)
	assert.Equal(t, mockReviews, res)
}

func TestManageJobUseCase_UpdateCartStatus_Success(t *testing.T) {
	cartRepo := new(MockUserCartRepository)
	uc := jobusecase.NewManageJobUseCase(nil, nil, nil, nil, cartRepo, &MockClock{})
	ctx := context.Background()
	userID := "usr_1"
	cartID := "crt_1"
	status := jobdomain.CartApplied

	cartRepo.On("FindByID", ctx, cartID).Return(&jobdomain.UserCart{CartID: cartID, UserID: userID}, nil)
	cartRepo.On("UpdateStatus", ctx, cartID, status).Return(nil)

	err := uc.UpdateCartStatus(ctx, userID, cartID, status)
	assert.NoError(t, err)
}

func TestManageJobUseCase_UpdateCartStatus_Forbidden(t *testing.T) {
	cartRepo := new(MockUserCartRepository)
	uc := jobusecase.NewManageJobUseCase(nil, nil, nil, nil, cartRepo, &MockClock{})
	ctx := context.Background()
	userID := "usr_1"
	cartID := "crt_1"
	status := jobdomain.CartApplied

	cartRepo.On("FindByID", ctx, cartID).Return(&jobdomain.UserCart{CartID: cartID, UserID: "other_user"}, nil)

	err := uc.UpdateCartStatus(ctx, userID, cartID, status)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)
}

func TestManageJobUseCase_UpdateCartStatus_InvalidStatus(t *testing.T) {
	cartRepo := new(MockUserCartRepository)
	uc := jobusecase.NewManageJobUseCase(nil, nil, nil, nil, cartRepo, &MockClock{})
	ctx := context.Background()
	userID := "usr_1"
	cartID := "crt_1"
	status := jobdomain.CartStatus("invalid_status")

	cartRepo.On("FindByID", ctx, cartID).Return(&jobdomain.UserCart{CartID: cartID, UserID: userID}, nil)

	err := uc.UpdateCartStatus(ctx, userID, cartID, status)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidInput, err)
}

func TestManageJobUseCase_WriteReview_InsertError(t *testing.T) {
	reviewRepo := new(MockJobReviewRepository)
	uc := jobusecase.NewManageJobUseCase(nil, nil, nil, reviewRepo, nil, &MockClock{})
	ctx := context.Background()

	review := &jobdomain.JobReview{}
	reviewRepo.On("Insert", ctx, mock.AnythingOfType("*jobdomain.JobReview")).Return(errors.New("db error"))

	err := uc.WriteReview(ctx, "usr_1", "job_1", review)
	assert.Error(t, err)
}

func TestManageJobUseCase_GetJobDetail_NotFound(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	uc := jobusecase.NewManageJobUseCase(jobRepo, nil, nil, nil, nil, &MockClock{})
	ctx := context.Background()

	jobRepo.On("FindByID", ctx, "job_1").Return((*jobdomain.JobPosting)(nil), domain.ErrNotFound)

	_, _, _, err := uc.GetJobDetail(ctx, "job_1")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestManageJobUseCase_RemoveFromCart_NotFound(t *testing.T) {
	cartRepo := new(MockUserCartRepository)
	uc := jobusecase.NewManageJobUseCase(nil, nil, nil, nil, cartRepo, &MockClock{})
	ctx := context.Background()

	cartRepo.On("FindByID", ctx, "crt_1").Return((*jobdomain.UserCart)(nil), domain.ErrNotFound)

	err := uc.RemoveFromCart(ctx, "usr_1", "crt_1")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestManageJobUseCase_UpdateCartStatus_NotFound(t *testing.T) {
	cartRepo := new(MockUserCartRepository)
	uc := jobusecase.NewManageJobUseCase(nil, nil, nil, nil, cartRepo, &MockClock{})
	ctx := context.Background()

	cartRepo.On("FindByID", ctx, "crt_1").Return((*jobdomain.UserCart)(nil), domain.ErrNotFound)

	err := uc.UpdateCartStatus(ctx, "usr_1", "crt_1", jobdomain.CartApplied)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestManageJobUseCase_CreateJob_Success(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}
	uc := jobusecase.NewManageJobUseCase(jobRepo, nil, nil, nil, nil, clock)
	ctx := context.Background()

	job := &jobdomain.JobPosting{
		AgencyName: "Test Agency",
	}

	jobRepo.On("Upsert", ctx, job).Return(nil).Run(func(args mock.Arguments) {
		j := args.Get(1).(*jobdomain.JobPosting)
		assert.NotEmpty(t, j.JobID)
		assert.Equal(t, nowTime, j.PostedAt)
		assert.Equal(t, nowTime, j.UpdatedAt)
		assert.Equal(t, nowTime, j.ScrapeAt)
	})

	err := uc.CreateJob(ctx, job)
	assert.NoError(t, err)
}

func TestManageJobUseCase_UpdateJob_Success(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}
	uc := jobusecase.NewManageJobUseCase(jobRepo, nil, nil, nil, nil, clock)
	ctx := context.Background()

	existing := &jobdomain.JobPosting{
		JobID:    "job_123",
		PostedAt: nowTime.Add(-1 * time.Hour),
		ScrapeAt: nowTime.Add(-1 * time.Hour),
	}
	job := &jobdomain.JobPosting{
		JobID:      "job_123",
		AgencyName: "Updated Agency",
	}

	jobRepo.On("FindByID", ctx, "job_123").Return(existing, nil)
	jobRepo.On("Upsert", ctx, job).Return(nil).Run(func(args mock.Arguments) {
		j := args.Get(1).(*jobdomain.JobPosting)
		assert.Equal(t, existing.PostedAt, j.PostedAt)
		assert.Equal(t, existing.ScrapeAt, j.ScrapeAt)
		assert.Equal(t, nowTime, j.UpdatedAt)
	})

	err := uc.UpdateJob(ctx, job)
	assert.NoError(t, err)
}

func TestManageJobUseCase_UpdateJob_NotFound(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	uc := jobusecase.NewManageJobUseCase(jobRepo, nil, nil, nil, nil, &MockClock{})
	ctx := context.Background()

	job := &jobdomain.JobPosting{JobID: "job_123"}
	jobRepo.On("FindByID", ctx, "job_123").Return((*jobdomain.JobPosting)(nil), domain.ErrNotFound)

	err := uc.UpdateJob(ctx, job)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestManageJobUseCase_PatchJob_Success(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}
	uc := jobusecase.NewManageJobUseCase(jobRepo, nil, nil, nil, nil, clock)
	ctx := context.Background()

	existing := &jobdomain.JobPosting{
		JobID:      "job_123",
		AgencyName: "Original Agency",
		Position:   "Original Position",
	}

	jobRepo.On("FindByID", ctx, "job_123").Return(existing, nil)
	jobRepo.On("Upsert", ctx, existing).Return(nil).Run(func(args mock.Arguments) {
		j := args.Get(1).(*jobdomain.JobPosting)
		assert.Equal(t, "Patched Agency", j.AgencyName)
		assert.Equal(t, "Original Position", j.Position)
		assert.Equal(t, nowTime, j.UpdatedAt)
	})

	err := uc.PatchJob(ctx, "job_123", map[string]interface{}{"agency_name": "Patched Agency"})
	assert.NoError(t, err)
}

func TestManageJobUseCase_DeleteJob_Success(t *testing.T) {
	jobRepo := new(MockJobPostingRepository)
	uc := jobusecase.NewManageJobUseCase(jobRepo, nil, nil, nil, nil, &MockClock{})
	ctx := context.Background()

	jobRepo.On("Delete", ctx, "job_123").Return(nil)

	err := uc.DeleteJob(ctx, "job_123")
	assert.NoError(t, err)
}
