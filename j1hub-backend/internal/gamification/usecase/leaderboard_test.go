package gamificationusecase_test

import (
	"context"
	"testing"

	gamificationusecase "github.com/parada3456/wat_project-backend/internal/gamification/usecase"

	gamificationdomain "github.com/parada3456/wat_project-backend/internal/gamification/domain"
	userdomain "github.com/parada3456/wat_project-backend/internal/user/domain"

	"github.com/stretchr/testify/assert"
)

func TestLeaderboardUseCase_GetLeaderboard_Success(t *testing.T) {
	leaderRepo := new(MockLeaderboardRepository)
	profileRepo := new(MockProfileRepository)
	ubRepo := new(MockUserBadgeRepository)

	uc := gamificationusecase.NewLeaderboardUseCase(leaderRepo, profileRepo, ubRepo)

	ctx := context.Background()
	scope := "global"
	jobID := ""

	mockUsers := []userdomain.User{
		{
			UserID:             "usr_1",
			CurrentPhasePoints: 150,
			MissionStreak:      5,
		},
		{
			UserID:             "usr_2",
			CurrentPhasePoints: 100,
			MissionStreak:      2,
		},
		{
			UserID:             "usr_3",
			CurrentPhasePoints: 50,
			MissionStreak:      1,
		},
	}

	mockProfile1 := &userdomain.Profile{
		UserID:    "usr_1",
		FirstName: "John",
		LastName:  "Doe",
	}
	mockProfile2 := &userdomain.Profile{
		UserID:    "usr_2",
		FirstName: "Alice",
		LastName:  "",
	}
	mockProfileHidden := &userdomain.Profile{
		UserID:          "usr_3",
		FirstName:       "Bob",
		LastName:        "Smith",
		RadarVisibility: "hidden",
	}

	mockUserBadges := []gamificationdomain.UserBadge{
		{UserID: "usr_1", BadgeID: "badge_gold"},
	}
	mockAliceBadges := []gamificationdomain.UserBadge{
		{UserID: "usr_2", BadgeID: "badge_silver"},
		{UserID: "usr_2", BadgeID: "badge_bronze"},
	}

	leaderRepo.On("FindByScope", ctx, scope, jobID).Return(mockUsers, nil)
	profileRepo.On("FindByUserID", ctx, "usr_1").Return(mockProfile1, nil)
	profileRepo.On("FindByUserID", ctx, "usr_2").Return(mockProfile2, nil)
	profileRepo.On("FindByUserID", ctx, "usr_3").Return(mockProfileHidden, nil)

	ubRepo.On("FindByUser", ctx, "usr_1").Return(mockUserBadges, nil)
	ubRepo.On("FindByUser", ctx, "usr_2").Return(mockAliceBadges, nil)
	ubRepo.On("FindByUser", ctx, "usr_3").Return([]gamificationdomain.UserBadge{}, nil)

	results, err := uc.GetLeaderboard(ctx, scope, jobID)

	assert.NoError(t, err)
	assert.Len(t, results, 3)

	// User 1 details
	assert.Equal(t, 1, results[0].Rank)
	assert.Equal(t, "John D.", results[0].Name)
	assert.Equal(t, []string{"badge_gold"}, results[0].Badges)

	// User 2 details (empty last name)
	assert.Equal(t, 2, results[1].Rank)
	assert.Equal(t, "Alice", results[1].Name)
	assert.Equal(t, []string{"badge_silver", "badge_bronze"}, results[1].Badges)

	// User 3 details (profile visibility Hidden)
	assert.Equal(t, 3, results[2].Rank)
	assert.Equal(t, "J1 Student #3", results[2].Name)
}
