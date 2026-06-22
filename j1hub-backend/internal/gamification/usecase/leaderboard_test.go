package gamificationusecase_test

import (
	"context"
	"testing"

	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/stretchr/testify/assert"
)

func TestLeaderboardUseCase_GetLeaderboard_Success(t *testing.T) {
	leaderRepo := new(MockLeaderboardRepository)
	profileRepo := new(MockProfileRepository)
	ubRepo := new(MockUserBadgeRepository)

	uc := gamificationgamificationusecase.NewLeaderboardUseCase(leaderRepo, profileRepo, ubRepo)

	ctx := context.Background()
	scope := "global"
	jobID := ""

	mockUsers := []userdomain.User{
		{
			UserID:             "usr_1",
			FirstName:          "John",
			LastName:           "Doe",
			CurrentPhasePoints: 150,
			MissionStreak:      5,
		},
		{
			UserID:             "usr_2",
			FirstName:          "Alice",
			LastName:           "",
			CurrentPhasePoints: 100,
			MissionStreak:      2,
		},
		{
			UserID:             "usr_3",
			FirstName:          "Bob",
			LastName:           "Smith",
			CurrentPhasePoints: 50,
			MissionStreak:      1,
		},
	}

	mockProfileHidden := &userdomain.Profile{
		UserID:          "usr_3",
		RadarVisibility: "Hidden",
	}

	mockUserBadges := []gamificationdomain.UserBadge{
		{UserID: "usr_1", BadgeID: "badge_gold"},
	}

	leaderRepo.On("FindByScope", ctx, scope, jobID).Return(mockUsers, nil)
	profileRepo.On("FindByUserID", ctx, "usr_1").Return((*userdomain.Profile)(nil), nil)
	profileRepo.On("FindByUserID", ctx, "usr_2").Return((*userdomain.Profile)(nil), nil)
	profileRepo.On("FindByUserID", ctx, "usr_3").Return(mockProfileHidden, nil)

	ubRepo.On("FindByUser", ctx, "usr_1").Return(mockUserBadges, nil)
	ubRepo.On("FindByUser", ctx, "usr_2").Return([]gamificationdomain.UserBadge{}, nil)
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

	// User 3 details (profile visibility Hidden)
	assert.Equal(t, 3, results[2].Rank)
	assert.Equal(t, "J1 Student #3", results[2].Name)
}
