package postgres

import (
	"context"
	"fmt"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type jobRepo struct {
	pool *pgxpool.Pool
}

func NewJobRepository(pool *pgxpool.Pool) port.JobPostingRepository {
	return &jobRepo{pool: pool}
}

func (r *jobRepo) FindWithFilters(ctx context.Context, filters map[string]interface{}) ([]domain.JobPosting, error) {
	// Simple implementation with static filters for now, can be expanded to dynamic
	query := `SELECT job_id, agency_name, employer_title, position, position_type, location_city, location_state, group_location, us_sponsor, salary_range_min, salary_range_max, available_slots, description, source_url, scrape_at, posted_at, updated_at FROM job_postings WHERE 1=1`
	var args []interface{}
	i := 1
	if v, ok := filters["position_type"]; ok {
		query += fmt.Sprintf(" AND position_type = $%d", i)
		args = append(args, v)
		i++
	}
	if v, ok := filters["location_state"]; ok {
		query += fmt.Sprintf(" AND location_state = $%d", i)
		args = append(args, v)
		i++
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []domain.JobPosting
	for rows.Next() {
		var j domain.JobPosting
		if err := rows.Scan(&j.JobID, &j.AgencyName, &j.EmployerTitle, &j.Position, &j.PositionType, &j.LocationCity, &j.LocationState, &j.GroupLocation, &j.USSponsor, &j.SalaryRangeMin, &j.SalaryRangeMax, &j.AvailableSlots, &j.Description, &j.SourceURL, &j.ScrapeAt, &j.PostedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (r *jobRepo) FindByID(ctx context.Context, id string) (*domain.JobPosting, error) {
	var j domain.JobPosting
	err := r.pool.QueryRow(ctx, `SELECT job_id, agency_name, employer_title, position, position_type, location_city, location_state, group_location, us_sponsor, salary_range_min, salary_range_max, available_slots, description, source_url, scrape_at, posted_at, updated_at FROM job_postings WHERE job_id = $1`, id).Scan(&j.JobID, &j.AgencyName, &j.EmployerTitle, &j.Position, &j.PositionType, &j.LocationCity, &j.LocationState, &j.GroupLocation, &j.USSponsor, &j.SalaryRangeMin, &j.SalaryRangeMax, &j.AvailableSlots, &j.Description, &j.SourceURL, &j.ScrapeAt, &j.PostedAt, &j.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &j, err
}

func (r *jobRepo) Upsert(ctx context.Context, job *domain.JobPosting) error {
	query := `
		INSERT INTO job_postings (
			job_id, agency_name, employer_title, position, position_type,
			location_city, location_state, group_location, us_sponsor,
			salary_range_min, salary_range_max, available_slots, description,
			source_url, scrape_at, posted_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (job_id) DO UPDATE SET
			employer_title = EXCLUDED.employer_title,
			position = EXCLUDED.position,
			position_type = EXCLUDED.position_type,
			location_city = EXCLUDED.location_city,
			location_state = EXCLUDED.location_state,
			group_location = EXCLUDED.group_location,
			us_sponsor = EXCLUDED.us_sponsor,
			salary_range_min = EXCLUDED.salary_range_min,
			salary_range_max = EXCLUDED.salary_range_max,
			available_slots = EXCLUDED.available_slots,
			description = EXCLUDED.description,
			source_url = EXCLUDED.source_url,
			scrape_at = EXCLUDED.scrape_at,
			updated_at = NOW()`
	_, err := r.pool.Exec(ctx, query,
		job.JobID, job.AgencyName, job.EmployerTitle, job.Position, job.PositionType,
		job.LocationCity, job.LocationState, job.GroupLocation, job.USSponsor,
		job.SalaryRangeMin, job.SalaryRangeMax, job.AvailableSlots, job.Description,
		job.SourceURL, job.ScrapeAt, job.PostedAt, job.UpdatedAt)
	return err
}

type jobHousingRepo struct {
	pool *pgxpool.Pool
}

func NewJobHousingRepository(pool *pgxpool.Pool) port.JobHousingRepository {
	return &jobHousingRepo{pool: pool}
}

func (r *jobHousingRepo) FindByJobID(ctx context.Context, jobID string) ([]domain.JobHousing, error) {
	rows, err := r.pool.Query(ctx, `SELECT housing_id, job_id, description, weekly_rate, deposit, transportation, range_min_start_date, range_max_start_date, created_at, updated_at FROM job_housings WHERE job_id = $1`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var housings []domain.JobHousing
	for rows.Next() {
		var h domain.JobHousing
		if err := rows.Scan(&h.HousingID, &h.JobID, &h.Description, &h.WeeklyRate, &h.Deposit, &h.Transportation, &h.RangeMinStartDate, &h.RangeMaxStartDate, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		housings = append(housings, h)
	}
	return housings, nil
}

func (r *jobHousingRepo) Upsert(ctx context.Context, h *domain.JobHousing) error {
	query := `
		INSERT INTO job_housings (
			housing_id, job_id, description, weekly_rate, deposit,
			transportation, range_min_start_date, range_max_start_date, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (housing_id) DO UPDATE SET
			weekly_rate = EXCLUDED.weekly_rate,
			deposit = EXCLUDED.deposit,
			transportation = EXCLUDED.transportation,
			range_min_start_date = EXCLUDED.range_min_start_date,
			range_max_start_date = EXCLUDED.range_max_start_date,
			updated_at = NOW()`
	_, err := r.pool.Exec(ctx, query,
		h.HousingID, h.JobID, h.Description, h.WeeklyRate, h.Deposit,
		h.Transportation, h.RangeMinStartDate, h.RangeMaxStartDate, h.CreatedAt, h.UpdatedAt)
	return err
}

type jobRatingRepo struct {
	pool *pgxpool.Pool
}

func NewJobOverallRatingRepository(pool *pgxpool.Pool) port.JobOverallRatingRepository {
	return &jobRatingRepo{pool: pool}
}

func (r *jobRatingRepo) FindByJobID(ctx context.Context, jobID string) (*domain.JobOverallRating, error) {
	var rating domain.JobOverallRating
	err := r.pool.QueryRow(ctx, `SELECT rating_summary_id, job_id, overall_rate, agency_rate, job_rate, coworkers_rate, town_rate, hours_rate, housing_rate, second_job_feasibility_rate, overtime_availability_rate, review_count, updated_at FROM job_overall_ratings WHERE job_id = $1`, jobID).Scan(&rating.RatingSummaryID, &rating.JobID, &rating.OverallRate, &rating.AgencyRate, &rating.JobRate, &rating.CoworkersRate, &rating.TownRate, &rating.HoursRate, &rating.HousingRate, &rating.SecondJobFeasibilityRate, &rating.OvertimeAvailabilityRate, &rating.ReviewCount, &rating.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &rating, err
}

func (r *jobRatingRepo) Recalculate(ctx context.Context, jobID string) error {
	query := `
		INSERT INTO job_overall_ratings (rating_summary_id, job_id, overall_rate, agency_rate, job_rate, coworkers_rate, town_rate, hours_rate, housing_rate, second_job_feasibility_rate, overtime_availability_rate, review_count, updated_at)
		SELECT 
			'smr_' || $1, $1,
			AVG(rating_stars), AVG(score_agency), AVG(score_job), AVG(score_coworkers), AVG(score_town), AVG(score_hours), AVG(score_housing), AVG(score_second_job_feasibility), AVG(score_overtime_availability),
			COUNT(*), NOW()
		FROM job_reviews WHERE job_id = $1
		ON CONFLICT (job_id) DO UPDATE SET
			overall_rate = EXCLUDED.overall_rate,
			agency_rate = EXCLUDED.agency_rate,
			job_rate = EXCLUDED.job_rate,
			coworkers_rate = EXCLUDED.coworkers_rate,
			town_rate = EXCLUDED.town_rate,
			hours_rate = EXCLUDED.hours_rate,
			housing_rate = EXCLUDED.housing_rate,
			second_job_feasibility_rate = EXCLUDED.second_job_feasibility_rate,
			overtime_availability_rate = EXCLUDED.overtime_availability_rate,
			review_count = EXCLUDED.review_count,
			updated_at = EXCLUDED.updated_at`
	_, err := r.pool.Exec(ctx, query, jobID)
	return err
}

type jobReviewRepo struct {
	pool *pgxpool.Pool
}

func NewJobReviewRepository(pool *pgxpool.Pool) port.JobReviewRepository {
	return &jobReviewRepo{pool: pool}
}

func (r *jobReviewRepo) Insert(ctx context.Context, rv *domain.JobReview) error {
	query := `INSERT INTO job_reviews (review_id, job_id, user_id, rating_stars, review_text, tips_for_next_generation, score_agency, score_job, score_coworkers, score_town, score_hours, score_housing, score_second_job_feasibility, score_overtime_availability, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
	_, err := r.pool.Exec(ctx, query, rv.ReviewID, rv.JobID, rv.UserID, rv.RatingStars, rv.ReviewText, rv.TipsForNextGeneration, rv.ScoreAgency, rv.ScoreJob, rv.ScoreCoworkers, rv.ScoreTown, rv.ScoreHours, rv.ScoreHousing, rv.ScoreSecondJobFeasibility, rv.ScoreOvertimeAvailability, rv.CreatedAt, rv.UpdatedAt)
	return err
}

func (r *jobReviewRepo) FindByJobID(ctx context.Context, jobID string) ([]domain.JobReview, error) {
	rows, err := r.pool.Query(ctx, `SELECT review_id, job_id, user_id, rating_stars, review_text, tips_for_next_generation, score_agency, score_job, score_coworkers, score_town, score_hours, score_housing, score_second_job_feasibility, score_overtime_availability, created_at, updated_at FROM job_reviews WHERE job_id = $1 ORDER BY created_at DESC`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []domain.JobReview
	for rows.Next() {
		var rv domain.JobReview
		if err := rows.Scan(&rv.ReviewID, &rv.JobID, &rv.UserID, &rv.RatingStars, &rv.ReviewText, &rv.TipsForNextGeneration, &rv.ScoreAgency, &rv.ScoreJob, &rv.ScoreCoworkers, &rv.ScoreTown, &rv.ScoreHours, &rv.ScoreHousing, &rv.ScoreSecondJobFeasibility, &rv.ScoreOvertimeAvailability, &rv.CreatedAt, &rv.UpdatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rv)
	}
	return reviews, nil
}

type userCartRepo struct {
	pool *pgxpool.Pool
}

func NewUserCartRepository(pool *pgxpool.Pool) port.UserCartRepository {
	return &userCartRepo{pool: pool}
}

func (r *userCartRepo) Insert(ctx context.Context, c *domain.UserCart) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO user_carts (cart_id, user_id, job_id, status, added_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		c.CartID, c.UserID, c.JobID, c.Status, c.AddedAt, c.UpdatedAt)
	return err
}

func (r *userCartRepo) FindByUserAndJob(ctx context.Context, userID, jobID string) (*domain.UserCart, error) {
	var c domain.UserCart
	err := r.pool.QueryRow(ctx, `SELECT cart_id, user_id, job_id, status, added_at, updated_at FROM user_carts WHERE user_id = $1 AND job_id = $2`, userID, jobID).Scan(&c.CartID, &c.UserID, &c.JobID, &c.Status, &c.AddedAt, &c.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &c, err
}

func (r *userCartRepo) FindByID(ctx context.Context, id string) (*domain.UserCart, error) {
	var c domain.UserCart
	err := r.pool.QueryRow(ctx, `SELECT cart_id, user_id, job_id, status, added_at, updated_at FROM user_carts WHERE cart_id = $1`, id).Scan(&c.CartID, &c.UserID, &c.JobID, &c.Status, &c.AddedAt, &c.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &c, err
}

func (r *userCartRepo) UpdateStatus(ctx context.Context, id string, status domain.CartStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE user_carts SET status = $1, updated_at = NOW() WHERE cart_id = $2`, status, id)
	return err
}
