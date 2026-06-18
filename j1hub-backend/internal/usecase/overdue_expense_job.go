package usecase

import (
	"context"
	"log"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/pkg/uid"
)

type OverdueExpenseJob struct {
	splitRepo  port.ExpenseSplitRepository
	creditRepo port.CreditScoreRepository
	ledgerRepo port.PointLedgerRepository
	notifier   port.NotifierPort
}

func NewOverdueExpenseJob(splitRepo port.ExpenseSplitRepository, creditRepo port.CreditScoreRepository, ledgerRepo port.PointLedgerRepository, notifier port.NotifierPort) *OverdueExpenseJob {
	return &OverdueExpenseJob{splitRepo: splitRepo, creditRepo: creditRepo, ledgerRepo: ledgerRepo, notifier: notifier}
}

func (j *OverdueExpenseJob) Run(ctx context.Context) error {
	splits, err := j.splitRepo.FindOverdue(ctx)
	if err != nil {
		return err
	}

	count := 0
	for _, s := range splits {
		if err := j.splitRepo.UpdatePaymentStatus(ctx, s.SplitID, domain.PaymentOverdue, s.PayslipURL); err != nil {
			log.Printf("failed to update split %s: %v", s.SplitID, err)
			continue
		}

		if err := j.creditRepo.Decrement(ctx, s.UserID, 10); err != nil {
			log.Printf("failed to decrement credit for user %s: %v", s.UserID, err)
		}

		// Point ledger for audit
		ledger := domain.PointLedger{
			LedgerID:   uid.New("ldg_"),
			UserID:     s.UserID,
			SourceType: domain.SourceExpensePenalty,
			SourceID:   s.SplitID,
			Delta:      -10,
			// ... balance after needs to be calculated or skipped if not critical for job
			Note: "Overdue expense penalty",
		}
		j.ledgerRepo.Insert(ctx, &ledger)

		j.notifier.Send(ctx, s.UserID, "Overdue payment", "Your payment is overdue. Credit score -10.")
		count++
	}

	log.Printf("Processed %d overdue expenses", count)
	return nil
}
