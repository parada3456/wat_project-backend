package usecase

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/pkg/timeutil"
	"github.com/j1hub/backend/pkg/uid"
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
	txn := &domain.ExpenseTransaction{
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

	var splits []domain.ExpenseSplit
	for _, s := range cmd.Splits {
		if s.UserID == payerID {
			return domain.ErrSelfSplit
		}
		splits = append(splits, domain.ExpenseSplit{
			SplitID:        uid.New("spl_"),
			TransactionID:  txn.TransactionID,
			UserID:         s.UserID,
			OweAmount:      s.OweAmount,
			PaymentStatus:  domain.PaymentPending,
			ApprovalStatus: domain.ApprovalPending,
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

	return uc.splitRepo.UpdatePaymentStatus(ctx, splitID, domain.PaymentSubmitted, url)
}

func (uc *ManageExpenseUseCase) ApproveSplit(ctx context.Context, payerID, splitID string) error {
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
	if err := uc.splitRepo.UpdateApproval(ctx, splitID, domain.ApprovalApproved, &now); err != nil {
		return err
	}

	uc.notifier.Send(ctx, split.UserID, "Payment approved", "Your payment has been approved!")

	return nil
}

func (uc *ManageExpenseUseCase) ListExpenses(ctx context.Context, userID string) ([]domain.ExpenseTransaction, error) {
	// Need FindByUser in repo
	return nil, nil
}

func (uc *ManageExpenseUseCase) GetExpenseDetail(ctx context.Context, userID, transactionID string) (*domain.ExpenseTransaction, []domain.ExpenseSplit, error) {
	txn, err := uc.txnRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, nil, err
	}
	// Fetch splits for txn
	return txn, nil, nil
}

func (uc *ManageExpenseUseCase) DeleteExpense(ctx context.Context, userID, transactionID string) error {
	// Check ownership and delete
	return nil
}

func (uc *ManageExpenseUseCase) ListPendingExpenses(ctx context.Context, userID string) ([]domain.ExpenseSplit, error) {
	// Need FindPendingByUser in repo
	return nil, nil
}
