package usecase_test

import (
	"context"
	"testing"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestRadarUseCase_GetRadar_Success(t *testing.T) {
	profileRepo := new(MockProfileRepository)
	radarRepo := new(MockRadarRepository)
	friendRepo := new(MockFriendshipRepository)

	cfg := &config.Config{
		RadarRadiusMeters: 5000,
		RadarStaleMinutes: 30,
	}

	uc := usecase.NewRadarUseCase(cfg, profileRepo, radarRepo, friendRepo)

	ctx := context.Background()
	requesterID := "usr_req"

	mockProfile := &domain.Profile{
		UserID: requesterID,
		Lat:    40.7128,
		Lng:    -74.0060,
	}

	nearbyProfiles := []domain.Profile{
		{
			UserID:          "usr_friend",
			Lat:             40.7130,
			Lng:             -74.0058,
			RadarVisibility: domain.VisibilityShowFriends,
			AvatarURL:       "friend_avatar.png",
		},
		{
			UserID:          "usr_anon",
			Lat:             40.7140,
			Lng:             -74.0050,
			RadarVisibility: domain.VisibilityShowAnonymous,
			AvatarURL:       "anon_avatar.png",
		},
		{
			UserID:          "usr_hidden",
			Lat:             40.7150,
			Lng:             -74.0040,
			RadarVisibility: domain.VisibilityHidden,
		},
	}

	profileRepo.On("FindByUserID", ctx, requesterID).Return(mockProfile, nil)
	radarRepo.On("FindNearby", ctx, mockProfile.Lat, mockProfile.Lng, 5000.0, 30).Return(nearbyProfiles, nil)

	// Friend friendship check: mock as accepted friend
	friendRepo.On("FindByCanonicalPair", ctx, "usr_friend", "usr_req").Return(&domain.Friendship{
		Status: domain.FriendshipAccepted,
	}, nil)

	// Anon friendship check: mock as not a friend
	friendRepo.On("FindByCanonicalPair", ctx, "usr_anon", "usr_req").Return((*domain.Friendship)(nil), domain.ErrNotFound)

	results, err := uc.GetRadar(ctx, requesterID)

	assert.NoError(t, err)
	// usr_friend and usr_anon should be returned, usr_hidden is excluded
	assert.Len(t, results, 2)

	assert.Equal(t, "usr_friend", results[0].UserID)
	assert.Equal(t, "Real Name", results[0].Name)
	assert.Equal(t, "friend_avatar.png", results[0].AvatarURL)

	assert.Equal(t, "usr_anon", results[1].UserID)
	assert.Equal(t, "Real Name", results[1].Name)
	assert.Equal(t, "anon_avatar.png", results[1].AvatarURL)
}
