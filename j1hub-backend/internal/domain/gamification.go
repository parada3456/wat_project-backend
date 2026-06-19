package domain

import (
	"log"
	"time"
)

type SourceType string

const (
	SourceMissionBase    SourceType = "Mission_Base"
	SourceSpeedBonus     SourceType = "Speed_Bonus"
	SourceStreakBonus    SourceType = "Streak_Bonus"
	SourceFirstCompleter SourceType = "First_Completer"
	SourceExpensePenalty SourceType = "Expense_Penalty"
	SourceAdminAdjust    SourceType = "Admin_Adjust"
)

func (s SourceType) Valid() bool {
	log.Println("debugprint: entering (SourceType).Valid")
	switch s {
	case SourceMissionBase, SourceSpeedBonus, SourceStreakBonus, SourceFirstCompleter, SourceExpensePenalty, SourceAdminAdjust:
		return true
	}
	return false
}

type TriggerType string

const (
	TriggerSpeed          TriggerType = "Speed"
	TriggerStreak         TriggerType = "Streak"
	TriggerFirstCompleter TriggerType = "First_Completer"
	TriggerPhaseComplete  TriggerType = "Phase_Complete"
	TriggerManual         TriggerType = "Manual"
)

func (t TriggerType) Valid() bool {
	log.Println("debugprint: entering (TriggerType).Valid")
	switch t {
	case TriggerSpeed, TriggerStreak, TriggerFirstCompleter, TriggerPhaseComplete, TriggerManual:
		return true
	}
	return false
}

type PointLedger struct {
	LedgerID             string
	UserID               string
	SourceType           SourceType
	SourceID             string
	Delta                int
	LifetimeBalanceAfter int
	PhaseBalanceAfter    int
	Note                 string
	CreatedAt            time.Time
}

type Badge struct {
	BadgeID     string
	Title       string
	Description string
	TriggerType TriggerType
	IconURL     string
	CreatedAt   time.Time
}

type UserBadge struct {
	UserBadgeID string
	UserID      string
	BadgeID     string
	SourceID    string
	EarnedAt    time.Time
}

type CreditScore struct {
	CreditID     string
	UserID       string
	CurrentScore int
	LastUpdated  time.Time
}

type PointReward struct {
	Base                int
	SpeedBonus          int
	StreakBonus         int
	FirstCompleterBonus int
	Total               int
}
