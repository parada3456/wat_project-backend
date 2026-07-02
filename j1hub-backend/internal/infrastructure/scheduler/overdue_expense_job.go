package scheduler

import (
	"context"
	"log"

	expensedomain "github.com/parada3456/wat_project-backend/internal/expense/domain"
	expenseport "github.com/parada3456/wat_project-backend/internal/expense/port"
	gamificationdomain "github.com/parada3456/wat_project-backend/internal/gamification/domain"
	gamificationport "github.com/parada3456/wat_project-backend/internal/gamification/port"
	notificationport "github.com/parada3456/wat_project-backend/internal/notification/port"
	"github.com/parada3456/wat_project-backend/pkg/uid"
)

type OverdueExpenseJob struct {
	splitRepo  expenseport.ExpenseSplitRepository
	creditRepo gamificationport.CreditScoreRepository
	ledgerRepo gamificationport.PointLedgerRepository
	notifier   notificationport.NotifierPort
}

func NewOverdueExpenseJob(splitRepo expenseport.ExpenseSplitRepository, creditRepo gamificationport.CreditScoreRepository, ledgerRepo gamificationport.PointLedgerRepository, notifier notificationport.NotifierPort) *OverdueExpenseJob {
	log.Println("debugprint: entering NewOverdueExpenseJob")
	return &OverdueExpenseJob{splitRepo: splitRepo, creditRepo: creditRepo, ledgerRepo: ledgerRepo, notifier: notifier}
}

func (j *OverdueExpenseJob) Run(ctx context.Context) error {
	log.Println("debugprint: entering (*OverdueExpenseJob).Run")
	splits, err := j.splitRepo.FindOverdue(ctx)
	if err != nil {
		return err
	}

	count := 0
	for _, s := range splits {
		if err := j.splitRepo.UpdatePaymentStatus(ctx, s.SplitID, expensedomain.PaymentOverdue, s.PayslipURL); err != nil {
			log.Printf("failed to update split %s: %v", s.SplitID, err)
			continue
		}

		if err := j.creditRepo.Decrement(ctx, s.UserID, 10); err != nil {
			log.Printf("failed to decrement credit for user %s: %v", s.UserID, err)
		}

		// Point ledger for audit
		ledger := gamificationdomain.PointLedger{
			LedgerID:   uid.New("ldg_"),
			UserID:     s.UserID,
			SourceType: gamificationdomain.SourceExpensePenalty,
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
