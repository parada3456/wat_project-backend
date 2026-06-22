package usecase

import (
	"time"

	"github.com/j1hub/backend/internal/port"
)

type UserDTO struct {
	UserID    string    `json:"userId"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
}

type ProfileDTO struct {
	ProfileID       string    `json:"profileId"`
	Bio             string    `json:"bio"`
	AvatarURL       string    `json:"avatarUrl"`
	RadarVisibility string    `json:"radarVisibility"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type CreditScoreDTO struct {
	CurrentScore int       `json:"currentScore"`
	LastUpdated  time.Time `json:"lastUpdated"`
}

type ProfileResponseDTO struct {
	User        UserDTO         `json:"user"`
	Profile     ProfileDTO      `json:"profile"`
	CreditScore *CreditScoreDTO `json:"creditScore,omitempty"`
}

type PublicProfileDTO struct {
	User    UserDTO    `json:"user"`
	Profile ProfileDTO `json:"profile"`
}

type AuthResponseDTO struct {
	User   UserDTO         `json:"user"`
	Tokens *port.TokenPair `json:"tokens"`
}

type JobPostingDTO struct {
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

type JobHousingDTO struct {
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

type JobOverallRatingDTO struct {
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

type JobDetailDTO struct {
	Job     JobPostingDTO       `json:"job"`
	Housing []JobHousingDTO     `json:"housing"`
	Rating  JobOverallRatingDTO `json:"rating"`
}

type UserCartDTO struct {
	CartID    string    `json:"cart_id"`
	UserID    string    `json:"user_id"`
	JobID     string    `json:"job_id"`
	Status    string    `json:"status"`
	AddedAt   time.Time `json:"added_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type JobReviewDTO struct {
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

type MissionDTO struct {
	MissionID            string     `json:"mission_id"`
	PhaseID              string     `json:"phase_id"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	Location             string     `json:"location"`
	BasePoints           int        `json:"base_points"`
	IsMandatory          bool       `json:"is_mandatory"`
	VerificationType     string     `json:"verification_type"`
	DueDateType          string     `json:"due_date_type"`
	FixedDueDate         *time.Time `json:"fixed_due_date"`
	RelativeTriggerEvent string     `json:"relative_trigger_event"`
	RelativeDaysOffset   int        `json:"relative_days_offset"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type UserMissionDTO struct {
	UserMissionID             string     `json:"user_mission_id"`
	UserID                    string     `json:"user_id"`
	MissionID                 string     `json:"mission_id"`
	Status                    string     `json:"status"`
	CalculatedDueDate         time.Time  `json:"calculated_due_date"`
	ProofURL                  string     `json:"proof_url"`
	ProofSubmittedAt          *time.Time `json:"proof_submitted_at"`
	VerifiedAt                *time.Time `json:"verified_at"`
	VerifiedBy                string     `json:"verified_by"`
	BasePointsEarned          int        `json:"base_points_earned"`
	SpeedBonusPoints          int        `json:"speed_bonus_points"`
	StreakBonusPoints         int        `json:"streak_bonus_points"`
	FirstCompleterBonusPoints int        `json:"first_completer_bonus_points"`
	TotalPointsEarned         int        `json:"total_points_earned"`
	RewardedAt                *time.Time `json:"rewarded_at"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}

type TaskDTO struct {
	TaskID      string    `json:"task_id"`
	MissionID   string    `json:"mission_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserTaskDTO struct {
	UserTaskID    string     `json:"user_task_id"`
	UserID        string     `json:"user_id"`
	TaskID        string     `json:"task_id"`
	UserMissionID string     `json:"user_mission_id"`
	IsCompleted   bool       `json:"is_completed"`
	CompletedAt   *time.Time `json:"completed_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type MissionDetailDTO struct {
	Mission     MissionDTO     `json:"mission"`
	UserMission UserMissionDTO `json:"user_mission"`
	Tasks       []TaskDTO      `json:"tasks"`
	UserTasks   []UserTaskDTO  `json:"user_tasks"`
}

type JourneyPhaseDTO struct {
	PhaseID     string    `json:"phase_id"`
	PhaseNumber int       `json:"phase_number"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserPhaseHistoryDTO struct {
	HistoryID         string     `json:"history_id"`
	UserID            string     `json:"user_id"`
	PhaseID           string     `json:"phase_id"`
	PhasePointsEarned int        `json:"phase_points_earned"`
	EnteredAt         time.Time  `json:"entered_at"`
	CompletedAt       *time.Time `json:"completed_at"`
}

type UserBadgeDTO struct {
	UserBadgeID string    `json:"user_badge_id"`
	UserID      string    `json:"user_id"`
	BadgeID     string    `json:"badge_id"`
	SourceID    string    `json:"source_id"`
	EarnedAt    time.Time `json:"earned_at"`
}

type PointLedgerDTO struct {
	LedgerID             string    `json:"ledger_id"`
	UserID               string    `json:"user_id"`
	SourceType           string    `json:"source_type"`
	SourceID             string    `json:"source_id"`
	Delta                int       `json:"delta"`
	LifetimeBalanceAfter int       `json:"lifetime_balance_after"`
	PhaseBalanceAfter    int       `json:"phase_balance_after"`
	Note                 string    `json:"note"`
	CreatedAt            time.Time `json:"created_at"`
}

type FriendshipDTO struct {
	FriendshipID string    `json:"friendship_id"`
	UserID1      string    `json:"user_id1"`
	UserID2      string    `json:"user_id2"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ExpenseTransactionDTO struct {
	TransactionID   string    `json:"transaction_id"`
	PaidByUserID    string    `json:"paid_by_user_id"`
	Title           string    `json:"title"`
	TotalAmount     float64   `json:"total_amount"`
	Currency        string    `json:"currency"`
	Memo            string    `json:"memo"`
	TransactionDate time.Time `json:"transaction_date"`
	DueDate         time.Time `json:"due_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ExpenseSplitDTO struct {
	SplitID        string     `json:"split_id"`
	TransactionID  string     `json:"transaction_id"`
	UserID         string     `json:"user_id"`
	OweAmount      float64    `json:"owe_amount"`
	PaymentStatus  string     `json:"payment_status"`
	PaymentMethod  string     `json:"payment_method"`
	PayslipURL     string     `json:"payslip_url"`
	ApprovalStatus string     `json:"approval_status"`
	SettledAt      *time.Time `json:"settled_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ExpenseDetailDTO struct {
	Transaction ExpenseTransactionDTO `json:"transaction"`
	Splits      []ExpenseSplitDTO     `json:"splits"`
}
