package usecase_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestManageExpenseUseCase_CreateExpense_Success(t *testing.T) {
	txnRepo := new(MockExpenseTransactionRepository)
	splitRepo := new(MockExpenseSplitRepository)
	storage := new(MockStoragePort)
	notifier := new(MockNotifierPort)
	
	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}

	uc := usecase.NewManageExpenseUseCase(txnRepo, splitRepo, storage, notifier, clock)

	ctx := context.Background()
	payerID := "usr_payer"
	dueDate := nowTime.Add(24 * time.Hour)

	cmd := usecase.CreateExpenseCmd{
		Title:       "Dinner",
		TotalAmount: 100.0,
		Currency:    "USD",
		Memo:        "Dinner split",
		DueDate:     dueDate,
		Splits: []struct {
			UserID    string
			OweAmount float64
		}{
			{UserID: "usr_debtor", OweAmount: 50.0},
		},
	}

	txnRepo.On("Insert", ctx, mock.AnythingOfType("*domain.ExpenseTransaction")).Return(nil).Run(func(args mock.Arguments) {
		txn := args.Get(1).(*domain.ExpenseTransaction)
		assert.Equal(t, payerID, txn.PaidByUserID)
		assert.Equal(t, "Dinner", txn.Title)
		assert.Equal(t, 100.0, txn.TotalAmount)
	})

	splitRepo.On("BulkInsert", ctx, mock.AnythingOfType("[]domain.ExpenseSplit")).Return(nil).Run(func(args mock.Arguments) {
		splits := args.Get(1).([]domain.ExpenseSplit)
		assert.Len(t, splits, 1)
		assert.Equal(t, "usr_debtor", splits[0].UserID)
		assert.Equal(t, 50.0, splits[0].OweAmount)
	})

	notifier.On("Send", ctx, "usr_debtor", "New expense", "You owe 50.00 for Dinner").Return(nil)

	err := uc.CreateExpense(ctx, payerID, cmd)

	assert.NoError(t, err)
	txnRepo.AssertExpectations(t)
	splitRepo.AssertExpectations(t)
	notifier.AssertExpectations(t)
}

func TestManageExpenseUseCase_CreateExpense_SelfSplit(t *testing.T) {
	txnRepo := new(MockExpenseTransactionRepository)
	splitRepo := new(MockExpenseSplitRepository)
	clock := &MockClock{}

	uc := usecase.NewManageExpenseUseCase(txnRepo, splitRepo, nil, nil, clock)

	ctx := context.Background()
	payerID := "usr_payer"

	cmd := usecase.CreateExpenseCmd{
		Title:       "Dinner",
		TotalAmount: 100.0,
		Currency:    "USD",
		Splits: []struct {
			UserID    string
			OweAmount float64
		}{
			{UserID: "usr_payer", OweAmount: 50.0}, // Self split
		},
	}

	err := uc.CreateExpense(ctx, payerID, cmd)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrSelfSplit, err)
}

func TestManageExpenseUseCase_SubmitSlip_Success(t *testing.T) {
	txnRepo := new(MockExpenseTransactionRepository)
	splitRepo := new(MockExpenseSplitRepository)
	storage := new(MockStoragePort)
	clock := &MockClock{}

	uc := usecase.NewManageExpenseUseCase(txnRepo, splitRepo, storage, nil, clock)

	ctx := context.Background()
	debtorID := "usr_debtor"
	splitID := "spl_123"
	file := bytes.NewReader([]byte("fake_image"))

	mockSplit := &domain.ExpenseSplit{
		SplitID: splitID,
		UserID:  debtorID,
	}

	splitRepo.On("FindByID", ctx, splitID).Return(mockSplit, nil)
	storage.On("UploadFile", ctx, "slips", splitID, file, "image/jpeg").Return("http://storage.com/slip.jpg", nil)
	splitRepo.On("UpdatePaymentStatus", ctx, splitID, domain.PaymentSubmitted, "http://storage.com/slip.jpg").Return(nil)

	err := uc.SubmitSlip(ctx, debtorID, splitID, file, "image/jpeg")

	assert.NoError(t, err)
	splitRepo.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestManageExpenseUseCase_ApproveSplit_Success(t *testing.T) {
	txnRepo := new(MockExpenseTransactionRepository)
	splitRepo := new(MockExpenseSplitRepository)
	notifier := new(MockNotifierPort)
	
	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}

	uc := usecase.NewManageExpenseUseCase(txnRepo, splitRepo, nil, notifier, clock)

	ctx := context.Background()
	payerID := "usr_payer"
	splitID := "spl_123"
	txnID := "txn_456"
	debtorID := "usr_debtor"

	mockSplit := &domain.ExpenseSplit{
		SplitID:       splitID,
		UserID:        debtorID,
		TransactionID: txnID,
	}

	mockTxn := &domain.ExpenseTransaction{
		TransactionID: txnID,
		PaidByUserID:  payerID,
	}

	splitRepo.On("FindByID", ctx, splitID).Return(mockSplit, nil)
	txnRepo.On("FindByID", ctx, txnID).Return(mockTxn, nil)
	splitRepo.On("UpdateApproval", ctx, splitID, domain.ApprovalApproved, &nowTime).Return(nil)
	notifier.On("Send", ctx, debtorID, "Payment approved", "Your payment has been approved!").Return(nil)

	err := uc.ApproveSplit(ctx, payerID, splitID)

	assert.NoError(t, err)
	splitRepo.AssertExpectations(t)
	txnRepo.AssertExpectations(t)
	notifier.AssertExpectations(t)
}

func TestManageExpenseUseCase_ListExpenses_Stub(t *testing.T) {
	uc := usecase.NewManageExpenseUseCase(nil, nil, nil, nil, &MockClock{})
	res, err := uc.ListExpenses(context.Background(), "usr_1")
	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestManageExpenseUseCase_GetExpenseDetail_Stub(t *testing.T) {
	txnRepo := new(MockExpenseTransactionRepository)
	uc := usecase.NewManageExpenseUseCase(txnRepo, nil, nil, nil, &MockClock{})
	ctx := context.Background()
	txnID := "txn_1"
	mockTxn := &domain.ExpenseTransaction{TransactionID: txnID}

	txnRepo.On("FindByID", ctx, txnID).Return(mockTxn, nil)

	txn, splits, err := uc.GetExpenseDetail(ctx, "usr_1", txnID)
	assert.NoError(t, err)
	assert.Equal(t, mockTxn, txn)
	assert.Nil(t, splits)
}

func TestManageExpenseUseCase_DeleteExpense_Stub(t *testing.T) {
	uc := usecase.NewManageExpenseUseCase(nil, nil, nil, nil, &MockClock{})
	err := uc.DeleteExpense(context.Background(), "usr_1", "txn_1")
	assert.NoError(t, err)
}

func TestManageExpenseUseCase_ListPendingExpenses_Stub(t *testing.T) {
	uc := usecase.NewManageExpenseUseCase(nil, nil, nil, nil, &MockClock{})
	res, err := uc.ListPendingExpenses(context.Background(), "usr_1")
	assert.Nil(t, res)
	assert.NoError(t, err)
}
