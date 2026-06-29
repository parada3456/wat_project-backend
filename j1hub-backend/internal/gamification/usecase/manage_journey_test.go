package gamificationusecase_test

import (
	"context"
	"testing"

	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"

	gamificationusecase "github.com/j1hub/backend/internal/gamification/usecase"
	"github.com/stretchr/testify/assert"
)

func TestJourneyUseCase_ListUserBadges_Success(t *testing.T) {
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	badgeRepo := new(MockBadgeRepository)
	ubRepo := new(MockUserBadgeRepository)
	creditRepo := new(MockCreditScoreRepository)
	ledgerRepo := new(MockPointLedgerRepository)

	uc := gamificationusecase.NewJourneyUseCase(phaseRepo, historyRepo, badgeRepo, ubRepo, creditRepo, ledgerRepo)

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
	uc := gamificationusecase.NewJourneyUseCase(nil, nil, nil, nil, nil, nil)
	res, err := uc.ListPhases(context.Background())
	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestJourneyUseCase_ListPhases_Success(t *testing.T) {
	phaseRepo := new(MockJourneyPhaseRepository)
	uc := gamificationusecase.NewJourneyUseCase(phaseRepo, nil, nil, nil, nil, nil)

	ctx := context.Background()
	mockPhases := []missiondomain.JourneyPhase{
		{PhaseID: "p1", PhaseNumber: 1, Title: "Phase 1"},
		{PhaseID: "p2", PhaseNumber: 2, Title: "Phase 2"},
	}

	phaseRepo.On("ListAll", ctx).Return(mockPhases, nil)

	res, err := uc.ListPhases(ctx)

	assert.NoError(t, err)
	assert.Equal(t, mockPhases, res)
	phaseRepo.AssertExpectations(t)
}

func TestJourneyUseCase_GetHistory_Stub(t *testing.T) {
	uc := gamificationusecase.NewJourneyUseCase(nil, nil, nil, nil, nil, nil)
	res, err := uc.GetHistory(context.Background(), "usr_1")
	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestJourneyUseCase_GetHistory_Success(t *testing.T) {
	historyRepo := new(MockUserPhaseHistoryRepository)
	uc := gamificationusecase.NewJourneyUseCase(nil, historyRepo, nil, nil, nil, nil)

	ctx := context.Background()
	userID := "usr_123"
	mockHistory := []missiondomain.UserPhaseHistory{
		{HistoryID: "h1", UserID: userID, PhaseID: "p1"},
	}

	historyRepo.On("FindByUser", ctx, userID).Return(mockHistory, nil)

	res, err := uc.GetHistory(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockHistory, res)
	historyRepo.AssertExpectations(t)
}

func TestJourneyUseCase_GetCreditScoreHistory_Stub(t *testing.T) {
	uc := gamificationusecase.NewJourneyUseCase(nil, nil, nil, nil, nil, nil)
	res, err := uc.GetCreditScoreHistory(context.Background(), "usr_1")
	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestJourneyUseCase_GetCreditScoreHistory_Success(t *testing.T) {
	ledgerRepo := new(MockPointLedgerRepository)
	uc := gamificationusecase.NewJourneyUseCase(nil, nil, nil, nil, nil, ledgerRepo)

	ctx := context.Background()
	userID := "usr_123"
	mockLedger := []gamificationdomain.PointLedger{
		{LedgerID: "ldg_1", UserID: userID, SourceType: gamificationdomain.SourceExpensePenalty, Delta: -50},
	}

	ledgerRepo.On("FindByUserAndSourceType", ctx, userID, gamificationdomain.SourceExpensePenalty).Return(mockLedger, nil)

	res, err := uc.GetCreditScoreHistory(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockLedger, res)
	ledgerRepo.AssertExpectations(t)
}

func TestJourneyUseCase_GetPointsLedger_Success(t *testing.T) {
	phaseRepo := new(MockJourneyPhaseRepository)
	historyRepo := new(MockUserPhaseHistoryRepository)
	badgeRepo := new(MockBadgeRepository)
	ubRepo := new(MockUserBadgeRepository)
	creditRepo := new(MockCreditScoreRepository)
	ledgerRepo := new(MockPointLedgerRepository)

	uc := gamificationusecase.NewJourneyUseCase(phaseRepo, historyRepo, badgeRepo, ubRepo, creditRepo, ledgerRepo)

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
