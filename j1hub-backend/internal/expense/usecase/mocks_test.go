package expenseusecase_test

import (
	"context"
	"time"

	expensedomain "github.com/parada3456/wat_project-backend/internal/expense/domain"
	"github.com/stretchr/testify/mock"
)

type MockExpenseTransactionRepository struct {
	mock.Mock
}

func (m *MockExpenseTransactionRepository) Insert(ctx context.Context, t *expensedomain.ExpenseTransaction) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockExpenseTransactionRepository) FindByID(ctx context.Context, id string) (*expensedomain.ExpenseTransaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*expensedomain.ExpenseTransaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExpenseTransactionRepository) MarkSettled(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockExpenseTransactionRepository) FindByUser(ctx context.Context, userID string, limit, offset int) ([]expensedomain.ExpenseTransaction, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]expensedomain.ExpenseTransaction), args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func (m *MockExpenseTransactionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockExpenseSplitRepository struct {
	mock.Mock
}

func (m *MockExpenseSplitRepository) BulkInsert(ctx context.Context, splits []expensedomain.ExpenseSplit) error {
	args := m.Called(ctx, splits)
	return args.Error(0)
}

func (m *MockExpenseSplitRepository) FindByID(ctx context.Context, id string) (*expensedomain.ExpenseSplit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*expensedomain.ExpenseSplit), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExpenseSplitRepository) UpdatePaymentStatus(ctx context.Context, id string, status expensedomain.PaymentStatus, slipURL string) error {
	args := m.Called(ctx, id, status, slipURL)
	return args.Error(0)
}

func (m *MockExpenseSplitRepository) UpdateApproval(ctx context.Context, id string, status expensedomain.ApprovalStatus, settledAt *time.Time) error {
	args := m.Called(ctx, id, status, settledAt)
	return args.Error(0)
}

func (m *MockExpenseSplitRepository) FindOverdue(ctx context.Context) ([]expensedomain.ExpenseSplit, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]expensedomain.ExpenseSplit), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExpenseSplitRepository) CountUnsettled(ctx context.Context, transactionID string) (int, error) {
	args := m.Called(ctx, transactionID)
	return args.Int(0), args.Error(1)
}

func (m *MockExpenseSplitRepository) FindByUser(ctx context.Context, userID string, limit, offset int) ([]expensedomain.ExpenseSplit, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]expensedomain.ExpenseSplit), args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func (m *MockExpenseSplitRepository) FindPendingByUser(ctx context.Context, userID string, limit, offset int) ([]expensedomain.ExpenseSplit, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]expensedomain.ExpenseSplit), args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func (m *MockExpenseSplitRepository) FindByTransactionID(ctx context.Context, transactionID string) ([]expensedomain.ExpenseSplit, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) != nil {
		return args.Get(0).([]expensedomain.ExpenseSplit), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockExpenseSplitRepository) DeleteByTransactionID(ctx context.Context, transactionID string) error {
	args := m.Called(ctx, transactionID)
	return args.Error(0)
}
