package domain

import (
	"log"
	"time"
)

type UserMissionStatus string

const (
	StatusNotStarted          UserMissionStatus = "Not_Started"
	StatusInProgress          UserMissionStatus = "In_Progress"
	StatusPendingVerification UserMissionStatus = "Pending_Verification"
	StatusCompleted           UserMissionStatus = "Completed"
	StatusOverdue             UserMissionStatus = "Overdue"
)

func (s UserMissionStatus) Valid() bool {
	log.Println("debugprint: entering (UserMissionStatus).Valid")
	switch s {
	case StatusNotStarted, StatusInProgress, StatusPendingVerification, StatusCompleted, StatusOverdue:
		return true
	}
	return false
}

type VerificationType string

const (
	VerificationNone   VerificationType = "None"
	VerificationUpload VerificationType = "Upload"
	VerificationAdmin  VerificationType = "Admin"
)

func (v VerificationType) Valid() bool {
	log.Println("debugprint: entering (VerificationType).Valid")
	switch v {
	case VerificationNone, VerificationUpload, VerificationAdmin:
		return true
	}
	return false
}

type JourneyPhase struct {
	PhaseID     string
	PhaseNumber int
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserPhaseHistory struct {
	HistoryID         string
	UserID            string
	PhaseID           string
	PhasePointsEarned int
	EnteredAt         time.Time
	CompletedAt       *time.Time
}

type Mission struct {
	MissionID            string
	PhaseID              string
	Title                string
	Description          string
	Location             string
	BasePoints           int
	IsMandatory          bool
	VerificationType     VerificationType
	DueDateType          string // "Relative" or "Fixed"
	FixedDueDate         *time.Time
	RelativeTriggerEvent string // "arrival_date" or "job_start_date"
	RelativeDaysOffset   int
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (m *Mission) CalculateDueDate(triggerDate time.Time) time.Time {
	log.Println("debugprint: entering (*Mission).CalculateDueDate")
	if m.DueDateType == "Fixed" && m.FixedDueDate != nil {
		return *m.FixedDueDate
	}
	return triggerDate.AddDate(0, 0, m.RelativeDaysOffset)
}

type UserMission struct {
	UserMissionID             string
	UserID                    string
	MissionID                 string
	Status                    UserMissionStatus
	CalculatedDueDate         time.Time
	ProofURL                  string
	ProofSubmittedAt          *time.Time
	VerifiedAt                *time.Time
	VerifiedBy                string
	BasePointsEarned          int
	SpeedBonusPoints          int
	StreakBonusPoints         int
	FirstCompleterBonusPoints int
	TotalPointsEarned         int
	RewardedAt                *time.Time
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

type Task struct {
	TaskID      string
	MissionID   string
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserTask struct {
	UserTaskID    string
	UserID        string
	TaskID        string
	UserMissionID string
	IsCompleted   bool
	CompletedAt   *time.Time
	UpdatedAt     time.Time
}

func CanAdvancePhase(missions []UserMission) bool {
	log.
		// Plan says: returns true if all mandatory missions are Completed
		// But UserMission doesn't have IsMandatory, so we'd need to join or pass them in.
		// For now, assume the caller passes only mandatory missions if they want this check.
		// Or we change the signature to include IsMandatory.
		// Let's stick to the plan's suggestion but assume the input list is what we check.
		Println("debugprint: entering CanAdvancePhase")

	for _, m := range missions {
		if m.Status != StatusCompleted {
			return false
		}
	}
	return true
}
