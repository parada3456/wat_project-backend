package expenseusecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/parada3456/wat_project-backend/internal/domain"
	expensedomain "github.com/parada3456/wat_project-backend/internal/expense/domain"
	port "github.com/parada3456/wat_project-backend/internal/expense/port"
	"github.com/parada3456/wat_project-backend/pkg/timeutil"
	"github.com/parada3456/wat_project-backend/pkg/uid"
)

type ManageExpenseUseCase struct {
	txnRepo   port.ExpenseTransactionRepository
	splitRepo port.ExpenseSplitRepository
	storage   port.StoragePort
	notifier  port.NotifierPort
	clock     timeutil.Clock
}

func NewManageExpenseUseCase(
	txnRepo port.ExpenseTransactionRepository,
	splitRepo port.ExpenseSplitRepository,
	storage port.StoragePort,
	notifier port.NotifierPort,
	clock timeutil.Clock,
) *ManageExpenseUseCase {
	log.Println("debugprint: entering NewManageExpenseUseCase")
	return &ManageExpenseUseCase{
		txnRepo:   txnRepo,
		splitRepo: splitRepo,
		storage:   storage,
		notifier:  notifier,
		clock:     clock,
	}
}

type CreateExpenseCmd struct {
	Title       string
	TotalAmount float64
	Currency    string
	Memo        string
	DueDate     time.Time
	Splits      []struct {
		UserID    string
		OweAmount float64
	}
}

func (uc *ManageExpenseUseCase) CreateExpense(ctx context.Context, payerID string, cmd CreateExpenseCmd) error {
	log.Println("debugprint: entering (*ManageExpenseUseCase).CreateExpense")
	txn := &expensedomain.ExpenseTransaction{
		TransactionID:   uid.New("txn_"),
		PaidByUserID:    payerID,
		Title:           cmd.Title,
		TotalAmount:     cmd.TotalAmount,
		Currency:        cmd.Currency,
		Memo:            cmd.Memo,
		TransactionDate: uc.clock.Now(),
		DueDate:         cmd.DueDate,
		CreatedAt:       uc.clock.Now(),
		UpdatedAt:       uc.clock.Now(),
	}

	var splits []expensedomain.ExpenseSplit
	for _, s := range cmd.Splits {
		if s.UserID == payerID {
			return domain.ErrSelfSplit
		}
		splits = append(splits, expensedomain.ExpenseSplit{
			SplitID:        uid.New("spl_"),
			TransactionID:  txn.TransactionID,
			UserID:         s.UserID,
			OweAmount:      s.OweAmount,
			PaymentStatus:  expensedomain.PaymentPending,
			ApprovalStatus: expensedomain.ApprovalPending,
			UpdatedAt:      uc.clock.Now(),
		})
	}

	if err := uc.txnRepo.Insert(ctx, txn); err != nil {
		return err
	}
	if err := uc.splitRepo.BulkInsert(ctx, splits); err != nil {
		return err
	}

	for _, s := range splits {
		uc.notifier.Send(ctx, s.UserID, "New expense", fmt.Sprintf("You owe %.2f for %s", s.OweAmount, txn.Title))
	}

	return nil
}

func (uc *ManageExpenseUseCase) SubmitSlip(ctx context.Context, debtorID, splitID string, file io.Reader, contentType string) error {
	log.Println("debugprint: entering (*ManageExpenseUseCase).SubmitSlip")
	split, err := uc.splitRepo.FindByID(ctx, splitID)
	if err != nil {
		return err
	}
	if split.UserID != debtorID {
		return domain.ErrForbidden
	}

	url, err := uc.storage.UploadFile(ctx, "slips", splitID, file, contentType)
	if err != nil {
		return err
	}

	return uc.splitRepo.UpdatePaymentStatus(ctx, splitID, expensedomain.PaymentSubmitted, url)
}

func (uc *ManageExpenseUseCase) ApproveSplit(ctx context.Context, payerID, splitID string) error {
	log.Println("debugprint: entering (*ManageExpenseUseCase).ApproveSplit")
	split, err := uc.splitRepo.FindByID(ctx, splitID)
	if err != nil {
		return err
	}

	txn, err := uc.txnRepo.FindByID(ctx, split.TransactionID)
	if err != nil {
		return err
	}
	if txn.PaidByUserID != payerID {
		return domain.ErrForbidden
	}

	now := uc.clock.Now()
	if err := uc.splitRepo.UpdateApproval(ctx, splitID, expensedomain.ApprovalApproved, &now); err != nil {
		return err
	}

	uc.notifier.Send(ctx, split.UserID, "Payment approved", "Your payment has been approved!")

	return nil
}

func (uc *ManageExpenseUseCase) ListExpenses(ctx context.Context, userID string, page, pageSize int) ([]expensedomain.ExpenseTransaction, int, error) {
	log.Println("debugprint: entering (*ManageExpenseUseCase).ListExpenses")
	limit := pageSize
	offset := (page - 1) * pageSize
	return uc.txnRepo.FindByUser(ctx, userID, limit, offset)
}

func (uc *ManageExpenseUseCase) GetExpenseDetail(ctx context.Context, userID, transactionID string) (*expensedomain.ExpenseTransaction, []expensedomain.ExpenseSplit, error) {
	log.Println("debugprint: entering (*ManageExpenseUseCase).GetExpenseDetail")
	txn, err := uc.txnRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, nil, err
	}
	splits, err := uc.splitRepo.FindByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, nil, err
	}

	isAuthorized := txn.PaidByUserID == userID
	if !isAuthorized {
		for _, s := range splits {
			if s.UserID == userID {
				isAuthorized = true
				break
			}
		}
	}
	if !isAuthorized {
		return nil, nil, domain.ErrForbidden
	}

	return txn, splits, nil
}

func (uc *ManageExpenseUseCase) DeleteExpense(ctx context.Context, userID, transactionID string) error {
	log.Println("debugprint: entering (*ManageExpenseUseCase).DeleteExpense")
	txn, err := uc.txnRepo.FindByID(ctx, transactionID)
	if err != nil {
		return err
	}
	if txn.PaidByUserID != userID {
		return domain.ErrForbidden
	}
	if err := uc.splitRepo.DeleteByTransactionID(ctx, transactionID); err != nil {
		return err
	}
	return uc.txnRepo.Delete(ctx, transactionID)
}

func (uc *ManageExpenseUseCase) ListPendingExpenses(ctx context.Context, userID string, page, pageSize int) ([]expensedomain.ExpenseSplit, int, error) {
	log.Println("debugprint: entering (*ManageExpenseUseCase).ListPendingExpenses")
	limit := pageSize
	offset := (page - 1) * pageSize
	return uc.splitRepo.FindPendingByUser(ctx, userID, limit, offset)
}
