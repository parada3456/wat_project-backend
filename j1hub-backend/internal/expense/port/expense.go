package port

import (
	"context"
	"time"
	"io"
	expensedomain "github.com/j1hub/backend/internal/expense/domain"
)

type ExpenseTransactionRepository interface {
	Insert(ctx context.Context, t *expensedomain.ExpenseTransaction) error
	FindByID(ctx context.Context, id string) (*expensedomain.ExpenseTransaction, error)
	MarkSettled(ctx context.Context, id string) error
	FindByUser(ctx context.Context, userID string, limit, offset int) ([]expensedomain.ExpenseTransaction, int, error)
	Delete(ctx context.Context, id string) error
}

type ExpenseSplitRepository interface {
	BulkInsert(ctx context.Context, splits []expensedomain.ExpenseSplit) error
	FindByID(ctx context.Context, id string) (*expensedomain.ExpenseSplit, error)
	UpdatePaymentStatus(ctx context.Context, id string, status expensedomain.PaymentStatus, slipURL string) error
	UpdateApproval(ctx context.Context, id string, status expensedomain.ApprovalStatus, settledAt *time.Time) error
	FindOverdue(ctx context.Context) ([]expensedomain.ExpenseSplit, error)
	CountUnsettled(ctx context.Context, transactionID string) (int, error)
	FindByUser(ctx context.Context, userID string, limit, offset int) ([]expensedomain.ExpenseSplit, int, error)
	FindPendingByUser(ctx context.Context, userID string, limit, offset int) ([]expensedomain.ExpenseSplit, int, error)
	FindByTransactionID(ctx context.Context, transactionID string) ([]expensedomain.ExpenseSplit, error)
	DeleteByTransactionID(ctx context.Context, transactionID string) error
}

type StoragePort interface {
	UploadFile(ctx context.Context, bucket, key string, data io.Reader, contentType string) (url string, err error)
}

type NotifierPort interface {
	Send(ctx context.Context, userID, title, body string) error
}
