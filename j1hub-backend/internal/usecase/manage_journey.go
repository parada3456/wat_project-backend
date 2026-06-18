package usecase

import (
	"context"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
)

type JourneyUseCase struct {
	phaseRepo   port.JourneyPhaseRepository
	historyRepo port.UserPhaseHistoryRepository
	badgeRepo   port.BadgeRepository
	ubRepo      port.UserBadgeRepository
	creditRepo  port.CreditScoreRepository
}

func NewJourneyUseCase(
	phaseRepo port.JourneyPhaseRepository,
	historyRepo port.UserPhaseHistoryRepository,
	badgeRepo port.BadgeRepository,
	ubRepo port.UserBadgeRepository,
	creditRepo port.CreditScoreRepository,
) *JourneyUseCase {
	return &JourneyUseCase{
		phaseRepo:   phaseRepo,
		historyRepo: historyRepo,
		badgeRepo:   badgeRepo,
		ubRepo:      ubRepo,
		creditRepo:  creditRepo,
	}
}

func (uc *JourneyUseCase) ListPhases(ctx context.Context) ([]domain.JourneyPhase, error) {
	// Need ListAll in repo
	return nil, nil
}

func (uc *JourneyUseCase) GetHistory(ctx context.Context, userID string) ([]domain.UserPhaseHistory, error) {
	// Need FindByUser in repo
	return nil, nil
}

func (uc *JourneyUseCase) ListUserBadges(ctx context.Context, userID string) ([]domain.UserBadge, error) {
	return uc.ubRepo.FindByUser(ctx, userID)
}

func (uc *JourneyUseCase) GetCreditScoreHistory(ctx context.Context, userID string) ([]domain.PointLedger, error) {
	// Use point ledger repo
	return nil, nil
}
