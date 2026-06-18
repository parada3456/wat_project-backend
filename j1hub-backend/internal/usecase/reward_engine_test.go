package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestRewardEngine_Calculate(t *testing.T) {
	cfg := &config.Config{
		Reward: config.RewardConfig{
			SpeedBonus7dPct:    20,
			SpeedBonus1dPct:    10,
			Streak3Pct:         10,
			Streak7Pct:         25,
			FirstCompleterFlat: 200,
		},
	}

	userRepo := new(MockUserRepository)
	umRepo := new(MockUserMissionRepository)
	engine := usecase.NewRewardEngine(cfg, userRepo, umRepo)

	ctx := context.Background()

	t.Run("Base points only, no bonuses", func(t *testing.T) {
		dueDate := time.Now().Add(1 * time.Hour)
		submitDate := time.Now()
		um := &domain.UserMission{
			CalculatedDueDate: dueDate,
			ProofSubmittedAt:  &submitDate, // less than 1 day before
		}
		user := &domain.User{
			MissionStreak: 0,
		}
		mission := &domain.Mission{
			BasePoints: 100,
		}

		reward, err := engine.Calculate(ctx, um, user, mission)
		assert.NoError(t, err)
		assert.Equal(t, 100, reward.Base)
		assert.Equal(t, 0, reward.SpeedBonus)
		assert.Equal(t, 0, reward.StreakBonus)
		assert.Equal(t, 100, reward.Total)
	})

	t.Run("7d Speed Bonus and 7+ Streak Bonus", func(t *testing.T) {
		dueDate := time.Now().Add(8 * 24 * time.Hour)
		submitDate := time.Now()
		um := &domain.UserMission{
			CalculatedDueDate: dueDate,
			ProofSubmittedAt:  &submitDate, // 8 days before -> speed bonus 7d (20%)
		}
		user := &domain.User{
			MissionStreak: 7, // streak bonus 7d (25%)
		}
		mission := &domain.Mission{
			BasePoints: 100,
		}

		reward, err := engine.Calculate(ctx, um, user, mission)
		assert.NoError(t, err)
		assert.Equal(t, 100, reward.Base)
		assert.Equal(t, 20, reward.SpeedBonus)  // 20% of 100
		assert.Equal(t, 25, reward.StreakBonus) // 25% of 100
		assert.Equal(t, 145, reward.Total)
	})

	t.Run("1d Speed Bonus and 3+ Streak Bonus", func(t *testing.T) {
		dueDate := time.Now().Add(2 * 24 * time.Hour)
		submitDate := time.Now()
		um := &domain.UserMission{
			CalculatedDueDate: dueDate,
			ProofSubmittedAt:  &submitDate, // 2 days before -> speed bonus 1d (10%)
		}
		user := &domain.User{
			MissionStreak: 4, // streak bonus 3d (10%)
		}
		mission := &domain.Mission{
			BasePoints: 200,
		}

		reward, err := engine.Calculate(ctx, um, user, mission)
		assert.NoError(t, err)
		assert.Equal(t, 200, reward.Base)
		assert.Equal(t, 20, reward.SpeedBonus)  // 10% of 200
		assert.Equal(t, 20, reward.StreakBonus) // 10% of 200
		assert.Equal(t, 240, reward.Total)
	})
}
