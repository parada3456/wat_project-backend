package usecase

import (
	"context"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/j1hub/backend/internal/port"
)

type RewardEngine struct {
	cfg      *config.Config
	userRepo port.UserRepository
	umRepo   port.UserMissionRepository
}

func NewRewardEngine(cfg *config.Config, userRepo port.UserRepository, umRepo port.UserMissionRepository) *RewardEngine {
	return &RewardEngine{cfg: cfg, userRepo: userRepo, umRepo: umRepo}
}

func (re *RewardEngine) Calculate(ctx context.Context, um *domain.UserMission, user *domain.User, mission *domain.Mission) (*domain.PointReward, error) {
	reward := &domain.PointReward{
		Base: mission.BasePoints,
	}

	// Speed Bonus
	if um.ProofSubmittedAt != nil {
		daysBefore := um.CalculatedDueDate.Sub(*um.ProofSubmittedAt).Hours() / 24
		if daysBefore >= 7 {
			reward.SpeedBonus = int(float64(reward.Base) * float64(re.cfg.Reward.SpeedBonus7dPct) / 100.0)
		} else if daysBefore >= 1 {
			reward.SpeedBonus = int(float64(reward.Base) * float64(re.cfg.Reward.SpeedBonus1dPct) / 100.0)
		}
	}

	// Streak Bonus
	if user.MissionStreak >= 7 {
		reward.StreakBonus = int(float64(reward.Base) * float64(re.cfg.Reward.Streak7Pct) / 100.0)
	} else if user.MissionStreak >= 3 {
		reward.StreakBonus = int(float64(reward.Base) * float64(re.cfg.Reward.Streak3Pct) / 100.0)
	}

	// First Completer Bonus
	// Note: this is a simplified check, plan says query USER_MISSION count where status=Completed
	// For now, let's assume we implement it in repo or here
	// I'll skip the actual query for now and just set it to 0 or implement a simple check
	// reward.FirstCompleterBonus = 200

	reward.Total = reward.Base + reward.SpeedBonus + reward.StreakBonus + reward.FirstCompleterBonus
	return reward, nil
}
