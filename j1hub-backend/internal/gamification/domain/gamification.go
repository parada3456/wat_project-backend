package gamificationdomain

import (
	"log"
	"time"
)

type SourceType string

const (
	SourceMissionBase    SourceType = "mission_base"
	SourceSpeedBonus     SourceType = "speed_bonus"
	SourceStreakBonus    SourceType = "streak_bonus"
	SourceFirstCompleter SourceType = "first_completer"
	SourceExpensePenalty SourceType = "expense_penalty"
	SourceAdminAdjust    SourceType = "admin_adjust"
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
	TriggerSpeed          TriggerType = "speed"
	TriggerStreak         TriggerType = "streak"
	TriggerFirstCompleter TriggerType = "first_completer"
	TriggerPhaseComplete  TriggerType = "phase_complete"
	TriggerManual         TriggerType = "manual"
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
	LedgerID             string     `json:"ledger_id"`
	UserID               string     `json:"user_id"`
	SourceType           SourceType `json:"source_type"`
	SourceID             string     `json:"source_id"`
	Delta                int        `json:"delta"`
	LifetimeBalanceAfter int        `json:"lifetime_balance_after"`
	PhaseBalanceAfter    int        `json:"phase_balance_after"`
	Note                 string     `json:"note"`
	CreatedAt            time.Time  `json:"created_at"`
}

type Badge struct {
	BadgeID     string      `json:"badge_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	TriggerType TriggerType `json:"trigger_type"`
	IconURL     string      `json:"icon_url"`
	CreatedAt   time.Time   `json:"created_at"`
}

type UserBadge struct {
	UserBadgeID string    `json:"user_badge_id"`
	UserID      string    `json:"user_id"`
	BadgeID     string    `json:"badge_id"`
	SourceID    string    `json:"source_id"`
	EarnedAt    time.Time `json:"earned_at"`
	Badge       *Badge    `json:"badge,omitempty"`
}

type CreditScore struct {
	CreditID     string    `json:"credit_id"`
	UserID       string    `json:"user_id"`
	CurrentScore int       `json:"current_score"`
	LastUpdated  time.Time `json:"last_updated"`
}

type PointReward struct {
	Base                int `json:"base"`
	SpeedBonus          int `json:"speed_bonus"`
	StreakBonus         int `json:"streak_bonus"`
	FirstCompleterBonus int `json:"first_completer_bonus"`
	Total               int `json:"total"`
}
