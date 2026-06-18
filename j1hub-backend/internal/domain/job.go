package domain

import (
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
	switch s {
	case CartSaved, CartViewed, CartApplied, CartRemoved:
		return true
	}
	return false
}

type JobPosting struct {
	JobID          string
	AgencyName     string
	EmployerTitle  string
	Position       string
	PositionType   string
	LocationCity   string
	LocationState  string
	GroupLocation  string
	USSponsor      bool
	SalaryRangeMin float64
	SalaryRangeMax float64
	AvailableSlots int
	Description    string
	SourceURL      string
	ScrapeAt       time.Time
	PostedAt       time.Time
	UpdatedAt      time.Time
}

type JobHousing struct {
	HousingID         string
	JobID             string
	Description       string
	WeeklyRate        float64
	Deposit           float64
	Transportation    string
	RangeMinStartDate time.Time
	RangeMaxStartDate time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type UserCart struct {
	CartID    string
	UserID    string
	JobID     string
	Status    CartStatus
	AddedAt   time.Time
	UpdatedAt time.Time
}

type JobOverallRating struct {
	RatingSummaryID          string
	JobID                    string
	OverallRate              float64
	AgencyRate               float64
	JobRate                  float64
	CoworkersRate            float64
	TownRate                 float64
	HoursRate                float64
	HousingRate              float64
	SecondJobFeasibilityRate float64
	OvertimeAvailabilityRate float64
	ReviewCount              int
	UpdatedAt                time.Time
}

type JobReview struct {
	ReviewID                  string
	JobID                     string
	UserID                    string
	RatingStars               float64
	ReviewText                string
	TipsForNextGeneration     string
	ScoreAgency               float64
	ScoreJob                  float64
	ScoreCoworkers            float64
	ScoreTown                 float64
	ScoreHours                float64
	ScoreHousing              float64
	ScoreSecondJobFeasibility float64
	ScoreOvertimeAvailability float64
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

func (r *JobReview) ScoreMap() map[string]float64 {
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
