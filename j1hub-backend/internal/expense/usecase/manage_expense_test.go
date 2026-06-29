package expenseusecase_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/j1hub/backend/internal/domain"
	expensedomain "github.com/j1hub/backend/internal/expense/domain"
	expenseusecase "github.com/j1hub/backend/internal/expense/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStoragePort struct {
	mock.Mock
}

func (m *MockStoragePort) UploadFile(ctx context.Context, bucket, key string, data io.Reader, contentType string) (string, error) {
	args := m.Called(ctx, bucket, key, data, contentType)
	return args.String(0), args.Error(1)
}

type MockNotifierPort struct {
	mock.Mock
}

func (m *MockNotifierPort) Send(ctx context.Context, userID, title, body string) error {
	args := m.Called(ctx, userID, title, body)
	return args.Error(0)
}

type MockClock struct {
	NowTime time.Time
}

func (m *MockClock) Now() time.Time {
	return m.NowTime
}

// Tests

func TestNewManageExpenseUseCase(t *testing.T) {
	uc := expenseusecase.NewManageExpenseUseCase(nil, nil, nil, nil, nil)
	assert.NotNil(t, uc)
}

func TestManageExpenseUseCase_CreateExpense(t *testing.T) {
	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		splitRepo := new(MockExpenseSplitRepository)
		storage := new(MockStoragePort)
		notifier := new(MockNotifierPort)
		clock := &MockClock{NowTime: nowTime}
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, storage, notifier, clock)

		cmd := expenseusecase.CreateExpenseCmd{
			Title:       "Dinner",
			TotalAmount: 100.0,
			Currency:    "USD",
			Splits: []struct {
				UserID    string
				OweAmount float64
			}{{UserID: "usr_debtor", OweAmount: 50.0}},
		}

		txnRepo.On("Insert", ctx, mock.AnythingOfType("*expensedomain.ExpenseTransaction")).Return(nil)
		splitRepo.On("BulkInsert", ctx, mock.AnythingOfType("[]expensedomain.ExpenseSplit")).Return(nil)
		notifier.On("Send", ctx, "usr_debtor", "New expense", "You owe 50.00 for Dinner").Return(nil)

		err := uc.CreateExpense(ctx, "usr_payer", cmd)
		assert.NoError(t, err)
		txnRepo.AssertExpectations(t)
		splitRepo.AssertExpectations(t)
		notifier.AssertExpectations(t)
	})

	t.Run("Self split error", func(t *testing.T) {
		clock := &MockClock{NowTime: nowTime}
		uc := expenseusecase.NewManageExpenseUseCase(nil, nil, nil, nil, clock)
		cmd := expenseusecase.CreateExpenseCmd{
			Splits: []struct {
				UserID    string
				OweAmount float64
			}{{UserID: "usr_payer", OweAmount: 50.0}},
		}
		err := uc.CreateExpense(ctx, "usr_payer", cmd)
		assert.Equal(t, domain.ErrSelfSplit, err)
	})

	t.Run("Insert txn error", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		clock := &MockClock{NowTime: nowTime}
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, nil, nil, nil, clock)
		cmd := expenseusecase.CreateExpenseCmd{
			Splits: []struct {
				UserID    string
				OweAmount float64
			}{{UserID: "usr_debtor", OweAmount: 50.0}},
		}
		expectedErr := errors.New("db error")
		txnRepo.On("Insert", ctx, mock.AnythingOfType("*expensedomain.ExpenseTransaction")).Return(expectedErr)
		err := uc.CreateExpense(ctx, "usr_payer", cmd)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("BulkInsert split error", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		splitRepo := new(MockExpenseSplitRepository)
		clock := &MockClock{NowTime: nowTime}
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, nil, nil, clock)
		cmd := expenseusecase.CreateExpenseCmd{
			Splits: []struct {
				UserID    string
				OweAmount float64
			}{{UserID: "usr_debtor", OweAmount: 50.0}},
		}
		expectedErr := errors.New("db error")
		txnRepo.On("Insert", ctx, mock.AnythingOfType("*expensedomain.ExpenseTransaction")).Return(nil)
		splitRepo.On("BulkInsert", ctx, mock.AnythingOfType("[]expensedomain.ExpenseSplit")).Return(expectedErr)
		err := uc.CreateExpense(ctx, "usr_payer", cmd)
		assert.Equal(t, expectedErr, err)
	})
}

func TestManageExpenseUseCase_SubmitSlip(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		splitRepo := new(MockExpenseSplitRepository)
		storage := new(MockStoragePort)
		uc := expenseusecase.NewManageExpenseUseCase(nil, splitRepo, storage, nil, &MockClock{})

		mockSplit := &expensedomain.ExpenseSplit{SplitID: "spl_123", UserID: "usr_debtor"}
		file := bytes.NewReader([]byte("fake"))

		splitRepo.On("FindByID", ctx, "spl_123").Return(mockSplit, nil)
		storage.On("UploadFile", ctx, "slips", "spl_123", file, "image/jpeg").Return("http://url", nil)
		splitRepo.On("UpdatePaymentStatus", ctx, "spl_123", expensedomain.PaymentSubmitted, "http://url").Return(nil)

		err := uc.SubmitSlip(ctx, "usr_debtor", "spl_123", file, "image/jpeg")
		assert.NoError(t, err)
	})

	t.Run("FindByID error", func(t *testing.T) {
		splitRepo := new(MockExpenseSplitRepository)
		uc := expenseusecase.NewManageExpenseUseCase(nil, splitRepo, nil, nil, &MockClock{})
		expectedErr := errors.New("not found")
		splitRepo.On("FindByID", ctx, "spl_123").Return(nil, expectedErr)
		err := uc.SubmitSlip(ctx, "usr_debtor", "spl_123", nil, "image/jpeg")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Forbidden", func(t *testing.T) {
		splitRepo := new(MockExpenseSplitRepository)
		uc := expenseusecase.NewManageExpenseUseCase(nil, splitRepo, nil, nil, &MockClock{})
		mockSplit := &expensedomain.ExpenseSplit{SplitID: "spl_123", UserID: "other"}
		splitRepo.On("FindByID", ctx, "spl_123").Return(mockSplit, nil)
		err := uc.SubmitSlip(ctx, "usr_debtor", "spl_123", nil, "image/jpeg")
		assert.Equal(t, domain.ErrForbidden, err)
	})

	t.Run("Upload error", func(t *testing.T) {
		splitRepo := new(MockExpenseSplitRepository)
		storage := new(MockStoragePort)
		uc := expenseusecase.NewManageExpenseUseCase(nil, splitRepo, storage, nil, &MockClock{})
		mockSplit := &expensedomain.ExpenseSplit{SplitID: "spl_123", UserID: "usr_debtor"}
		file := bytes.NewReader([]byte("fake"))
		expectedErr := errors.New("upload failed")
		splitRepo.On("FindByID", ctx, "spl_123").Return(mockSplit, nil)
		storage.On("UploadFile", ctx, "slips", "spl_123", file, "image/jpeg").Return("", expectedErr)
		err := uc.SubmitSlip(ctx, "usr_debtor", "spl_123", file, "image/jpeg")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Update status error", func(t *testing.T) {
		splitRepo := new(MockExpenseSplitRepository)
		storage := new(MockStoragePort)
		uc := expenseusecase.NewManageExpenseUseCase(nil, splitRepo, storage, nil, &MockClock{})
		mockSplit := &expensedomain.ExpenseSplit{SplitID: "spl_123", UserID: "usr_debtor"}
		file := bytes.NewReader([]byte("fake"))
		expectedErr := errors.New("update failed")
		splitRepo.On("FindByID", ctx, "spl_123").Return(mockSplit, nil)
		storage.On("UploadFile", ctx, "slips", "spl_123", file, "image/jpeg").Return("http://url", nil)
		splitRepo.On("UpdatePaymentStatus", ctx, "spl_123", expensedomain.PaymentSubmitted, "http://url").Return(expectedErr)
		err := uc.SubmitSlip(ctx, "usr_debtor", "spl_123", file, "image/jpeg")
		assert.Equal(t, expectedErr, err)
	})
}

func TestManageExpenseUseCase_ApproveSplit(t *testing.T) {
	ctx := context.Background()
	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)

	t.Run("Success", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		splitRepo := new(MockExpenseSplitRepository)
		notifier := new(MockNotifierPort)
		clock := &MockClock{NowTime: nowTime}
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, nil, notifier, clock)

		mockSplit := &expensedomain.ExpenseSplit{SplitID: "spl_123", TransactionID: "txn_1", UserID: "usr_debtor"}
		mockTxn := &expensedomain.ExpenseTransaction{TransactionID: "txn_1", PaidByUserID: "usr_payer"}

		splitRepo.On("FindByID", ctx, "spl_123").Return(mockSplit, nil)
		txnRepo.On("FindByID", ctx, "txn_1").Return(mockTxn, nil)
		splitRepo.On("UpdateApproval", ctx, "spl_123", expensedomain.ApprovalApproved, &nowTime).Return(nil)
		notifier.On("Send", ctx, "usr_debtor", "Payment approved", "Your payment has been approved!").Return(nil)

		err := uc.ApproveSplit(ctx, "usr_payer", "spl_123")
		assert.NoError(t, err)
	})

	t.Run("FindByID split error", func(t *testing.T) {
		splitRepo := new(MockExpenseSplitRepository)
		uc := expenseusecase.NewManageExpenseUseCase(nil, splitRepo, nil, nil, &MockClock{})
		expectedErr := errors.New("not found")
		splitRepo.On("FindByID", ctx, "spl_123").Return(nil, expectedErr)
		err := uc.ApproveSplit(ctx, "usr_payer", "spl_123")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("FindByID txn error", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		splitRepo := new(MockExpenseSplitRepository)
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, nil, nil, &MockClock{})
		mockSplit := &expensedomain.ExpenseSplit{SplitID: "spl_123", TransactionID: "txn_1"}
		expectedErr := errors.New("not found")
		splitRepo.On("FindByID", ctx, "spl_123").Return(mockSplit, nil)
		txnRepo.On("FindByID", ctx, "txn_1").Return(nil, expectedErr)
		err := uc.ApproveSplit(ctx, "usr_payer", "spl_123")
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Forbidden", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		splitRepo := new(MockExpenseSplitRepository)
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, nil, nil, &MockClock{})
		mockSplit := &expensedomain.ExpenseSplit{SplitID: "spl_123", TransactionID: "txn_1"}
		mockTxn := &expensedomain.ExpenseTransaction{TransactionID: "txn_1", PaidByUserID: "other"}
		splitRepo.On("FindByID", ctx, "spl_123").Return(mockSplit, nil)
		txnRepo.On("FindByID", ctx, "txn_1").Return(mockTxn, nil)
		err := uc.ApproveSplit(ctx, "usr_payer", "spl_123")
		assert.Equal(t, domain.ErrForbidden, err)
	})

	t.Run("UpdateApproval error", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		splitRepo := new(MockExpenseSplitRepository)
		clock := &MockClock{NowTime: nowTime}
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, nil, nil, clock)

		mockSplit := &expensedomain.ExpenseSplit{SplitID: "spl_123", TransactionID: "txn_1"}
		mockTxn := &expensedomain.ExpenseTransaction{TransactionID: "txn_1", PaidByUserID: "usr_payer"}
		expectedErr := errors.New("update error")

		splitRepo.On("FindByID", ctx, "spl_123").Return(mockSplit, nil)
		txnRepo.On("FindByID", ctx, "txn_1").Return(mockTxn, nil)
		splitRepo.On("UpdateApproval", ctx, "spl_123", expensedomain.ApprovalApproved, &nowTime).Return(expectedErr)

		err := uc.ApproveSplit(ctx, "usr_payer", "spl_123")
		assert.Equal(t, expectedErr, err)
	})
}

func TestManageExpenseUseCase_ListExpenses(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, nil, nil, nil, &MockClock{})

		mockTxns := []expensedomain.ExpenseTransaction{
			{TransactionID: "txn_1", PaidByUserID: "usr_1", Title: "Rent"},
			{TransactionID: "txn_2", PaidByUserID: "usr_1", Title: "Food"},
		}
		txnRepo.On("FindByUser", ctx, "usr_1", 10, 0).Return(mockTxns, 2, nil)

		res, totalCount, err := uc.ListExpenses(ctx, "usr_1", 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, mockTxns, res)
		assert.Equal(t, 2, totalCount)
		txnRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, nil, nil, nil, &MockClock{})

		expectedErr := errors.New("db error")
		txnRepo.On("FindByUser", ctx, "usr_1", 10, 0).Return(nil, 0, expectedErr)

		res, totalCount, err := uc.ListExpenses(ctx, "usr_1", 1, 10)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, res)
		assert.Equal(t, 0, totalCount)
		txnRepo.AssertExpectations(t)
	})
}

func TestManageExpenseUseCase_GetExpenseDetail(t *testing.T) {
	ctx := context.Background()

	t.Run("Success as Payer", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		splitRepo := new(MockExpenseSplitRepository)
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, nil, nil, &MockClock{})

		mockTxn := &expensedomain.ExpenseTransaction{TransactionID: "txn_1", PaidByUserID: "usr_1"}
		mockSplits := []expensedomain.ExpenseSplit{
			{SplitID: "spl_1", TransactionID: "txn_1", UserID: "usr_2"},
		}

		txnRepo.On("FindByID", ctx, "txn_1").Return(mockTxn, nil)
		splitRepo.On("FindByTransactionID", ctx, "txn_1").Return(mockSplits, nil)

		txn, splits, err := uc.GetExpenseDetail(ctx, "usr_1", "txn_1")
		assert.NoError(t, err)
		assert.Equal(t, mockTxn, txn)
		assert.Equal(t, mockSplits, splits)
		txnRepo.AssertExpectations(t)
		splitRepo.AssertExpectations(t)
	})

	t.Run("Success as Debtor", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		splitRepo := new(MockExpenseSplitRepository)
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, nil, nil, &MockClock{})

		mockTxn := &expensedomain.ExpenseTransaction{TransactionID: "txn_1", PaidByUserID: "other"}
		mockSplits := []expensedomain.ExpenseSplit{
			{SplitID: "spl_1", TransactionID: "txn_1", UserID: "usr_1"},
		}

		txnRepo.On("FindByID", ctx, "txn_1").Return(mockTxn, nil)
		splitRepo.On("FindByTransactionID", ctx, "txn_1").Return(mockSplits, nil)

		txn, splits, err := uc.GetExpenseDetail(ctx, "usr_1", "txn_1")
		assert.NoError(t, err)
		assert.Equal(t, mockTxn, txn)
		assert.Equal(t, mockSplits, splits)
		txnRepo.AssertExpectations(t)
		splitRepo.AssertExpectations(t)
	})

	t.Run("Error Forbidden", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		splitRepo := new(MockExpenseSplitRepository)
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, nil, nil, &MockClock{})

		mockTxn := &expensedomain.ExpenseTransaction{TransactionID: "txn_1", PaidByUserID: "other"}
		mockSplits := []expensedomain.ExpenseSplit{
			{SplitID: "spl_1", TransactionID: "txn_1", UserID: "other_debtor"},
		}

		txnRepo.On("FindByID", ctx, "txn_1").Return(mockTxn, nil)
		splitRepo.On("FindByTransactionID", ctx, "txn_1").Return(mockSplits, nil)

		txn, splits, err := uc.GetExpenseDetail(ctx, "usr_1", "txn_1")
		assert.Equal(t, domain.ErrForbidden, err)
		assert.Nil(t, txn)
		assert.Nil(t, splits)
		txnRepo.AssertExpectations(t)
		splitRepo.AssertExpectations(t)
	})

	t.Run("Error Transaction Not Found", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, nil, nil, nil, &MockClock{})

		txnRepo.On("FindByID", ctx, "txn_1").Return(nil, domain.ErrNotFound)

		txn, splits, err := uc.GetExpenseDetail(ctx, "usr_1", "txn_1")
		assert.Equal(t, domain.ErrNotFound, err)
		assert.Nil(t, txn)
		assert.Nil(t, splits)
		txnRepo.AssertExpectations(t)
	})
}

func TestManageExpenseUseCase_DeleteExpense(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		splitRepo := new(MockExpenseSplitRepository)
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, splitRepo, nil, nil, &MockClock{})

		mockTxn := &expensedomain.ExpenseTransaction{TransactionID: "txn_1", PaidByUserID: "usr_1"}

		txnRepo.On("FindByID", ctx, "txn_1").Return(mockTxn, nil)
		splitRepo.On("DeleteByTransactionID", ctx, "txn_1").Return(nil)
		txnRepo.On("Delete", ctx, "txn_1").Return(nil)

		err := uc.DeleteExpense(ctx, "usr_1", "txn_1")
		assert.NoError(t, err)
		txnRepo.AssertExpectations(t)
		splitRepo.AssertExpectations(t)
	})

	t.Run("Forbidden", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, nil, nil, nil, &MockClock{})

		mockTxn := &expensedomain.ExpenseTransaction{TransactionID: "txn_1", PaidByUserID: "other"}

		txnRepo.On("FindByID", ctx, "txn_1").Return(mockTxn, nil)

		err := uc.DeleteExpense(ctx, "usr_1", "txn_1")
		assert.Equal(t, domain.ErrForbidden, err)
		txnRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		txnRepo := new(MockExpenseTransactionRepository)
		uc := expenseusecase.NewManageExpenseUseCase(txnRepo, nil, nil, nil, &MockClock{})

		txnRepo.On("FindByID", ctx, "txn_1").Return(nil, domain.ErrNotFound)

		err := uc.DeleteExpense(ctx, "usr_1", "txn_1")
		assert.Equal(t, domain.ErrNotFound, err)
		txnRepo.AssertExpectations(t)
	})
}

func TestManageExpenseUseCase_ListPendingExpenses(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		splitRepo := new(MockExpenseSplitRepository)
		uc := expenseusecase.NewManageExpenseUseCase(nil, splitRepo, nil, nil, &MockClock{})

		mockPendingSplits := []expensedomain.ExpenseSplit{
			{SplitID: "spl_1", PaymentStatus: expensedomain.PaymentPending},
			{SplitID: "spl_3", PaymentStatus: expensedomain.PaymentSubmitted},
		}

		splitRepo.On("FindPendingByUser", ctx, "usr_1", 10, 0).Return(mockPendingSplits, 2, nil)

		res, totalCount, err := uc.ListPendingExpenses(ctx, "usr_1", 1, 10)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, "spl_1", res[0].SplitID)
		assert.Equal(t, "spl_3", res[1].SplitID)
		assert.Equal(t, 2, totalCount)
		splitRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		splitRepo := new(MockExpenseSplitRepository)
		uc := expenseusecase.NewManageExpenseUseCase(nil, splitRepo, nil, nil, &MockClock{})

		expectedErr := errors.New("db error")
		splitRepo.On("FindPendingByUser", ctx, "usr_1", 10, 0).Return(nil, 0, expectedErr)

		res, totalCount, err := uc.ListPendingExpenses(ctx, "usr_1", 1, 10)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, res)
		assert.Equal(t, 0, totalCount)
		splitRepo.AssertExpectations(t)
	})
}
