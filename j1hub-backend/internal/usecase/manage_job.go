package usecase

import (
	"context"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/pkg/timeutil"
	"github.com/j1hub/backend/pkg/uid"
)

type ManageJobUseCase struct {
	jobRepo     port.JobPostingRepository
	housingRepo port.JobHousingRepository
	ratingRepo  port.JobOverallRatingRepository
	reviewRepo  port.JobReviewRepository
	cartRepo    port.UserCartRepository
	clock       timeutil.Clock
}

func NewManageJobUseCase(
	jobRepo port.JobPostingRepository,
	housingRepo port.JobHousingRepository,
	ratingRepo port.JobOverallRatingRepository,
	reviewRepo port.JobReviewRepository,
	cartRepo port.UserCartRepository,
	clock timeutil.Clock,
) *ManageJobUseCase {
	return &ManageJobUseCase{
		jobRepo:     jobRepo,
		housingRepo: housingRepo,
		ratingRepo:  ratingRepo,
		reviewRepo:  reviewRepo,
		cartRepo:    cartRepo,
		clock:       clock,
	}
}

func (uc *ManageJobUseCase) AddToCart(ctx context.Context, userID, jobID string) error {
	existing, err := uc.cartRepo.FindByUserAndJob(ctx, userID, jobID)
	if err == nil && existing != nil {
		return domain.ErrConflict
	}

	cart := &domain.UserCart{
		CartID:    uid.New("crt_"),
		UserID:    userID,
		JobID:     jobID,
		Status:    domain.CartSaved,
		AddedAt:   uc.clock.Now(),
		UpdatedAt: uc.clock.Now(),
	}

	return uc.cartRepo.Insert(ctx, cart)
}

func (uc *ManageJobUseCase) WriteReview(ctx context.Context, userID, jobID string, rv *domain.JobReview) error {
	rv.ReviewID = uid.New("rvw_")
	rv.UserID = userID
	rv.JobID = jobID
	rv.CreatedAt = uc.clock.Now()
	rv.UpdatedAt = uc.clock.Now()

	if err := uc.reviewRepo.Insert(ctx, rv); err != nil {
		return err
	}

	return uc.ratingRepo.Recalculate(ctx, jobID)
}

func (uc *ManageJobUseCase) ListJobs(ctx context.Context, filters map[string]interface{}) ([]domain.JobPosting, error) {
	return uc.jobRepo.FindWithFilters(ctx, filters)
}

func (uc *ManageJobUseCase) GetJobDetail(ctx context.Context, jobID string) (*domain.JobPosting, []domain.JobHousing, *domain.JobOverallRating, error) {
	job, err := uc.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, nil, nil, err
	}
	housing, _ := uc.housingRepo.FindByJobID(ctx, jobID)
	rating, _ := uc.ratingRepo.FindByJobID(ctx, jobID)
	return job, housing, rating, nil
}

func (uc *ManageJobUseCase) ListCart(ctx context.Context, userID string) ([]domain.UserCart, error) {
	// Need FindByUser in repo
	return nil, nil
}

func (uc *ManageJobUseCase) RemoveFromCart(ctx context.Context, userID, cartID string) error {
	// Check ownership and delete
	return nil
}
