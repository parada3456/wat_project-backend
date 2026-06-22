package jobdomain

import (
	"log"
	"time"
)

type CartStatus string

const (
	CartSaved   CartStatus = "Saved"
	CartViewed  CartStatus = "Viewed"
	CartApplied CartStatus = "Applied"
	CartRemoved CartStatus = "Removed"
)

func (s CartStatus) Valid() bool {
	log.Println("debugprint: entering (CartStatus).Valid")
	switch s {
	case CartSaved, CartViewed, CartApplied, CartRemoved:
		return true
	}
	return false
}

type JobPosting struct {
	JobID          string    `json:"job_id"`
	AgencyName     string    `json:"agency_name"`
	EmployerTitle  string    `json:"employer_title"`
	Position       string    `json:"position"`
	PositionType   string    `json:"position_type"`
	LocationCity   string    `json:"location_city"`
	LocationState  string    `json:"location_state"`
	GroupLocation  string    `json:"group_location"`
	USSponsor      bool      `json:"us_sponsor"`
	SalaryRangeMin float64   `json:"salary_range_min"`
	SalaryRangeMax float64   `json:"salary_range_max"`
	AvailableSlots int       `json:"available_slots"`
	Description    string    `json:"description"`
	SourceURL      string    `json:"source_url"`
	ScrapeAt       time.Time `json:"scrape_at"`
	PostedAt       time.Time `json:"posted_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type JobHousing struct {
	HousingID         string    `json:"housing_id"`
	JobID             string    `json:"job_id"`
	Description       string    `json:"description"`
	WeeklyRate        float64   `json:"weekly_rate"`
	Deposit           float64   `json:"deposit"`
	Transportation    string    `json:"transportation"`
	RangeMinStartDate time.Time `json:"range_min_start_date"`
	RangeMaxStartDate time.Time `json:"range_max_start_date"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type UserCart struct {
	CartID    string     `json:"cart_id"`
	UserID    string     `json:"user_id"`
	JobID     string     `json:"job_id"`
	Status    CartStatus `json:"status"`
	AddedAt   time.Time  `json:"added_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type JobOverallRating struct {
	RatingSummaryID          string    `json:"rating_summary_id"`
	JobID                    string    `json:"job_id"`
	OverallRate              float64   `json:"overall_rate"`
	AgencyRate               float64   `json:"agency_rate"`
	JobRate                  float64   `json:"job_rate"`
	CoworkersRate            float64   `json:"coworkers_rate"`
	TownRate                 float64   `json:"town_rate"`
	HoursRate                float64   `json:"hours_rate"`
	HousingRate              float64   `json:"housing_rate"`
	SecondJobFeasibilityRate float64   `json:"second_job_feasibility_rate"`
	OvertimeAvailabilityRate float64   `json:"overtime_availability_rate"`
	ReviewCount              int       `json:"review_count"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type JobReview struct {
	ReviewID                  string    `json:"review_id"`
	JobID                     string    `json:"job_id"`
	UserID                    string    `json:"user_id"`
	RatingStars               float64   `json:"rating_stars"`
	ReviewText                string    `json:"review_text"`
	TipsForNextGeneration     string    `json:"tips_for_next_generation"`
	ScoreAgency               float64   `json:"score_agency"`
	ScoreJob                  float64   `json:"score_job"`
	ScoreCoworkers            float64   `json:"score_coworkers"`
	ScoreTown                 float64   `json:"score_town"`
	ScoreHours                float64   `json:"score_hours"`
	ScoreHousing              float64   `json:"score_housing"`
	ScoreSecondJobFeasibility float64   `json:"score_second_job_feasibility"`
	ScoreOvertimeAvailability float64   `json:"score_overtime_availability"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

func (r *JobReview) ScoreMap() map[string]float64 {
	log.Println("debugprint: entering (*JobReview).ScoreMap")
	return map[string]float64{
		"agency":                 r.ScoreAgency,
		"job":                    r.ScoreJob,
		"coworkers":              r.ScoreCoworkers,
		"town":                   r.ScoreTown,
		"hours":                  r.ScoreHours,
		"housing":                r.ScoreHousing,
		"second_job_feasibility": r.ScoreSecondJobFeasibility,
		"overtime_availability":  r.ScoreOvertimeAvailability,
	}
}
