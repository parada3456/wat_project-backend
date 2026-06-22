package usecase

import (
	"context"
	"log"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
)

type JourneyUseCase struct {
	phaseRepo   port.JourneyPhaseRepository
	historyRepo port.UserPhaseHistoryRepository
	badgeRepo   port.BadgeRepository
	ubRepo      port.UserBadgeRepository
	creditRepo  port.CreditScoreRepository
	ledgerRepo  port.PointLedgerRepository
}

func NewJourneyUseCase(
	phaseRepo port.JourneyPhaseRepository,
	historyRepo port.UserPhaseHistoryRepository,
	badgeRepo port.BadgeRepository,
	ubRepo port.UserBadgeRepository,
	creditRepo port.CreditScoreRepository,
	ledgerRepo port.PointLedgerRepository,
) *JourneyUseCase {
	log.Println("debugprint: entering NewJourneyUseCase")
	return &JourneyUseCase{
		phaseRepo:   phaseRepo,
		historyRepo: historyRepo,
		badgeRepo:   badgeRepo,
		ubRepo:      ubRepo,
		creditRepo:  creditRepo,
		ledgerRepo:  ledgerRepo,
	}
}

func (uc *JourneyUseCase) ListPhases(ctx context.Context) ([]domain.JourneyPhase, error) {
	log.
		// Need ListAll in repo
		Println("debugprint: entering (*JourneyUseCase).ListPhases")

	return nil, nil
}

func (uc *JourneyUseCase) GetHistory(ctx context.Context, userID string) ([]domain.UserPhaseHistory, error) {
	log.
		// Need FindByUser in repo
		Println("debugprint: entering (*JourneyUseCase).GetHistory")

	return nil, nil
}

func (uc *JourneyUseCase) ListUserBadges(ctx context.Context, userID string) ([]domain.UserBadge, error) {
	log.Println("debugprint: entering (*JourneyUseCase).ListUserBadges")
	return uc.ubRepo.FindByUser(ctx, userID)
}

func (uc *JourneyUseCase) GetCreditScoreHistory(ctx context.Context, userID string) ([]domain.PointLedger, error) {
	log.Println("debugprint: entering (*JourneyUseCase).GetCreditScoreHistory")
	return nil, nil
}

func (uc *JourneyUseCase) GetPointsLedger(ctx context.Context, userID string) ([]domain.PointLedger, error) {
	log.Println("debugprint: entering (*JourneyUseCase).GetPointsLedger")
	return uc.ledgerRepo.FindByUser(ctx, userID)
}
