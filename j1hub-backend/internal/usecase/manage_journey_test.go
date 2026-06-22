package usecase_test

import (
	"context"
	"testing"

	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"

	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestJourneyUseCase_ListUserBadges_Success(t *testing.T) {
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	badgeRepo := new(MockBadgeRepository)
	ubRepo := new(MockUserBadgeRepository)
	creditRepo := new(MockCreditScoreRepository)
	ledgerRepo := new(MockPointLedgerRepository)

	uc := usecase.NewJourneyUseCase(phaseRepo, historyRepo, badgeRepo, ubRepo, creditRepo, ledgerRepo)

	ctx := context.Background()
	userID := "usr_123"
	mockBadges := []gamificationdomain.UserBadge{
		{UserID: userID, BadgeID: "badge_1"},
		{UserID: userID, BadgeID: "badge_2"},
	}

	ubRepo.On("FindByUser", ctx, userID).Return(mockBadges, nil)

	res, err := uc.ListUserBadges(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockBadges, res)
}

func TestJourneyUseCase_ListPhases_Stub(t *testing.T) {
	uc := usecase.NewJourneyUseCase(nil, nil, nil, nil, nil, nil)
	res, err := uc.ListPhases(context.Background())
	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestJourneyUseCase_GetHistory_Stub(t *testing.T) {
	uc := usecase.NewJourneyUseCase(nil, nil, nil, nil, nil, nil)
	res, err := uc.GetHistory(context.Background(), "usr_1")
	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestJourneyUseCase_GetCreditScoreHistory_Stub(t *testing.T) {
	uc := usecase.NewJourneyUseCase(nil, nil, nil, nil, nil, nil)
	res, err := uc.GetCreditScoreHistory(context.Background(), "usr_1")
	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestJourneyUseCase_GetPointsLedger_Success(t *testing.T) {
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	badgeRepo := new(MockBadgeRepository)
	ubRepo := new(MockUserBadgeRepository)
	creditRepo := new(MockCreditScoreRepository)
	ledgerRepo := new(MockPointLedgerRepository)

	uc := usecase.NewJourneyUseCase(phaseRepo, historyRepo, badgeRepo, ubRepo, creditRepo, ledgerRepo)

	ctx := context.Background()
	userID := "usr_123"
	mockLedger := []gamificationdomain.PointLedger{
		{LedgerID: "ldg_1", UserID: userID, Delta: 100},
	}

	ledgerRepo.On("FindByUser", ctx, userID).Return(mockLedger, nil)

	res, err := uc.GetPointsLedger(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockLedger, res)
}
