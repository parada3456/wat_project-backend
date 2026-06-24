package jobusecase

import (
	"context"
	"log"

	"github.com/j1hub/backend/internal/domain"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
	port "github.com/j1hub/backend/internal/job/port"
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
	log.Println("debugprint: entering NewManageJobUseCase")
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
	log.Println("debugprint: entering (*ManageJobUseCase).AddToCart")
	existing, err := uc.cartRepo.FindByUserAndJob(ctx, userID, jobID)
	if err == nil && existing != nil {
		return domain.ErrConflict
	}

	cart := &jobdomain.UserCart{
		CartID:    uid.New("crt_"),
		UserID:    userID,
		JobID:     jobID,
		Status:    jobdomain.CartSaved,
		AddedAt:   uc.clock.Now(),
		UpdatedAt: uc.clock.Now(),
	}

	return uc.cartRepo.Insert(ctx, cart)
}

func (uc *ManageJobUseCase) WriteReview(ctx context.Context, userID, jobID string, rv *jobdomain.JobReview) error {
	log.Println("debugprint: entering (*ManageJobUseCase).WriteReview")
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

func (uc *ManageJobUseCase) ListJobs(ctx context.Context, filters map[string]interface{}) ([]jobdomain.JobPosting, error) {
	log.Println("debugprint: entering (*ManageJobUseCase).ListJobs")
	return uc.jobRepo.FindWithFilters(ctx, filters)
}

func (uc *ManageJobUseCase) GetJobDetail(ctx context.Context, jobID string) (*jobdomain.JobPosting, []jobdomain.JobHousing, *jobdomain.JobOverallRating, error) {
	log.Println("debugprint: entering (*ManageJobUseCase).GetJobDetail")
	job, err := uc.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, nil, nil, err
	}
	housing, _ := uc.housingRepo.FindByJobID(ctx, jobID)
	rating, _ := uc.ratingRepo.FindByJobID(ctx, jobID)
	return job, housing, rating, nil
}

func (uc *ManageJobUseCase) ListCart(ctx context.Context, userID string) ([]jobdomain.UserCart, error) {
	log.Println("debugprint: entering (*ManageJobUseCase).ListCart")
	return uc.cartRepo.FindByUser(ctx, userID)
}

func (uc *ManageJobUseCase) RemoveFromCart(ctx context.Context, userID, cartID string) error {
	log.Println("debugprint: entering (*ManageJobUseCase).RemoveFromCart")
	cart, err := uc.cartRepo.FindByID(ctx, cartID)
	if err != nil {
		return err
	}
	if cart.UserID != userID {
		return domain.ErrForbidden
	}
	return uc.cartRepo.Delete(ctx, cartID)
}

func (uc *ManageJobUseCase) ListReviews(ctx context.Context, jobID string) ([]jobdomain.JobReview, error) {
	log.Println("debugprint: entering (*ManageJobUseCase).ListReviews")
	return uc.reviewRepo.FindByJobID(ctx, jobID)
}

func (uc *ManageJobUseCase) UpdateCartStatus(ctx context.Context, userID, cartID string, status jobdomain.CartStatus) error {
	log.Println("debugprint: entering (*ManageJobUseCase).UpdateCartStatus")
	cart, err := uc.cartRepo.FindByID(ctx, cartID)
	if err != nil {
		return err
	}
	if cart.UserID != userID {
		return domain.ErrForbidden
	}
	if !status.Valid() {
		return domain.ErrInvalidInput
	}
	return uc.cartRepo.UpdateStatus(ctx, cart.CartID, status)
}

func (uc *ManageJobUseCase) CreateJob(ctx context.Context, job *jobdomain.JobPosting) error {
	log.Println("debugprint: entering (*ManageJobUseCase).CreateJob")
	if job.JobID == "" {
		job.JobID = uid.New("job_")
	}
	job.PostedAt = uc.clock.Now()
	job.UpdatedAt = uc.clock.Now()
	if job.ScrapeAt.IsZero() {
		job.ScrapeAt = uc.clock.Now()
	}
	return uc.jobRepo.Upsert(ctx, job)
}

func (uc *ManageJobUseCase) UpdateJob(ctx context.Context, job *jobdomain.JobPosting) error {
	log.Println("debugprint: entering (*ManageJobUseCase).UpdateJob")
	existing, err := uc.jobRepo.FindByID(ctx, job.JobID)
	if err != nil {
		return err
	}
	if job.PostedAt.IsZero() {
		job.PostedAt = existing.PostedAt
	}
	if job.ScrapeAt.IsZero() {
		job.ScrapeAt = existing.ScrapeAt
	}
	job.UpdatedAt = uc.clock.Now()
	return uc.jobRepo.Upsert(ctx, job)
}

func (uc *ManageJobUseCase) PatchJob(ctx context.Context, jobID string, updates map[string]interface{}) error {
	log.Println("debugprint: entering (*ManageJobUseCase).PatchJob")
	existing, err := uc.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return err
	}

	for k, v := range updates {
		switch k {
		case "agency_name":
			if val, ok := v.(string); ok {
				existing.AgencyName = val
			}
		case "employer_title":
			if val, ok := v.(string); ok {
				existing.EmployerTitle = val
			}
		case "position":
			if val, ok := v.(string); ok {
				existing.Position = val
			}
		case "position_type":
			if val, ok := v.(string); ok {
				existing.PositionType = val
			}
		case "location_city":
			if val, ok := v.(string); ok {
				existing.LocationCity = val
			}
		case "location_state":
			if val, ok := v.(string); ok {
				existing.LocationState = val
			}
		case "group_location":
			if val, ok := v.(string); ok {
				existing.GroupLocation = val
			}
		case "us_sponsor":
			if val, ok := v.(bool); ok {
				existing.USSponsor = val
			}
		case "salary_range_min":
			if val, ok := v.(float64); ok {
				existing.SalaryRangeMin = val
			}
		case "salary_range_max":
			if val, ok := v.(float64); ok {
				existing.SalaryRangeMax = val
			}
		case "available_slots":
			if val, ok := v.(float64); ok {
				existing.AvailableSlots = int(val)
			} else if val, ok := v.(int); ok {
				existing.AvailableSlots = val
			}
		case "description":
			if val, ok := v.(string); ok {
				existing.Description = val
			}
		case "source_url":
			if val, ok := v.(string); ok {
				existing.SourceURL = val
			}
		}
	}
	existing.UpdatedAt = uc.clock.Now()
	return uc.jobRepo.Upsert(ctx, existing)
}

func (uc *ManageJobUseCase) DeleteJob(ctx context.Context, jobID string) error {
	log.Println("debugprint: entering (*ManageJobUseCase).DeleteJob")
	return uc.jobRepo.Delete(ctx, jobID)
}
