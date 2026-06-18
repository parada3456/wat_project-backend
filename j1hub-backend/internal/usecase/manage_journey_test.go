package usecase_test

import (
	"context"
	"testing"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestJourneyUseCase_ListUserBadges_Success(t *testing.T) {
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	badgeRepo := new(MockBadgeRepository)
	ubRepo := new(MockUserBadgeRepository)
	creditRepo := new(MockCreditScoreRepository)

	uc := usecase.NewJourneyUseCase(phaseRepo, historyRepo, badgeRepo, ubRepo, creditRepo)

	ctx := context.Background()
	userID := "usr_123"
	mockBadges := []domain.UserBadge{
		{UserID: userID, BadgeID: "badge_1"},
		{UserID: userID, BadgeID: "badge_2"},
	}

	ubRepo.On("FindByUser", ctx, userID).Return(mockBadges, nil)

	res, err := uc.ListUserBadges(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockBadges, res)
}

func TestJourneyUseCase_ListPhases_Stub(t *testing.T) {
	uc := usecase.NewJourneyUseCase(nil, nil, nil, nil, nil)
	res, err := uc.ListPhases(context.Background())
	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestJourneyUseCase_GetHistory_Stub(t *testing.T) {
	uc := usecase.NewJourneyUseCase(nil, nil, nil, nil, nil)
	res, err := uc.GetHistory(context.Background(), "usr_1")
	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestJourneyUseCase_GetCreditScoreHistory_Stub(t *testing.T) {
	uc := usecase.NewJourneyUseCase(nil, nil, nil, nil, nil)
	res, err := uc.GetCreditScoreHistory(context.Background(), "usr_1")
	assert.Nil(t, res)
	assert.NoError(t, err)
}
