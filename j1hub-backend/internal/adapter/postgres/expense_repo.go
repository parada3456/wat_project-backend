package postgres

import (
	"context"
	"time"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type expenseRepo struct {
	pool *pgxpool.Pool
}

func NewExpenseRepository(pool *pgxpool.Pool) port.ExpenseTransactionRepository {
	return &expenseRepo{pool: pool}
}

func (r *expenseRepo) Insert(ctx context.Context, t *domain.ExpenseTransaction) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO expense_transactions (transaction_id, paid_by_user_id, title, total_amount, currency, memo, transaction_date, due_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		t.TransactionID, t.PaidByUserID, t.Title, t.TotalAmount, t.Currency, t.Memo, t.TransactionDate, t.DueDate, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *expenseRepo) FindByID(ctx context.Context, id string) (*domain.ExpenseTransaction, error) {
	var t domain.ExpenseTransaction
	err := r.pool.QueryRow(ctx, `SELECT transaction_id, paid_by_user_id, title, total_amount, currency, memo, transaction_date, due_date, created_at, updated_at FROM expense_transactions WHERE transaction_id = $1`, id).Scan(&t.TransactionID, &t.PaidByUserID, &t.Title, &t.TotalAmount, &t.Currency, &t.Memo, &t.TransactionDate, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &t, err
}

func (r *expenseRepo) MarkSettled(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE expense_transactions SET updated_at = NOW() WHERE transaction_id = $1`, id)
	return err
}

type expenseSplitRepo struct {
	pool *pgxpool.Pool
}

func NewExpenseSplitRepository(pool *pgxpool.Pool) port.ExpenseSplitRepository {
	return &expenseSplitRepo{pool: pool}
}

func (r *expenseSplitRepo) BulkInsert(ctx context.Context, splits []domain.ExpenseSplit) error {
	batch := &pgx.Batch{}
	for _, s := range splits {
		batch.Queue(`INSERT INTO expense_splits (split_id, transaction_id, user_id, owe_amount, payment_status, payment_method, payslip_url, approval_status, settled_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			s.SplitID, s.TransactionID, s.UserID, s.OweAmount, s.PaymentStatus, s.PaymentMethod, s.PayslipURL, s.ApprovalStatus, s.SettledAt, s.UpdatedAt)
	}
	return r.pool.SendBatch(ctx, batch).Close()
}

func (r *expenseSplitRepo) FindByID(ctx context.Context, id string) (*domain.ExpenseSplit, error) {
	var s domain.ExpenseSplit
	err := r.pool.QueryRow(ctx, `SELECT split_id, transaction_id, user_id, owe_amount, payment_status, payment_method, payslip_url, approval_status, settled_at, updated_at FROM expense_splits WHERE split_id = $1`, id).Scan(&s.SplitID, &s.TransactionID, &s.UserID, &s.OweAmount, &s.PaymentStatus, &s.PaymentMethod, &s.PayslipURL, &s.ApprovalStatus, &s.SettledAt, &s.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &s, err
}

func (r *expenseSplitRepo) UpdatePaymentStatus(ctx context.Context, id string, status domain.PaymentStatus, slipURL string) error {
	_, err := r.pool.Exec(ctx, `UPDATE expense_splits SET payment_status = $1, payslip_url = $2, updated_at = NOW() WHERE split_id = $3`, status, slipURL, id)
	return err
}

func (r *expenseSplitRepo) UpdateApproval(ctx context.Context, id string, status domain.ApprovalStatus, settledAt *time.Time) error {
	_, err := r.pool.Exec(ctx, `UPDATE expense_splits SET approval_status = $1, payment_status = CASE WHEN $1 = 'Approved' THEN 'Approved' ELSE payment_status END, settled_at = $2, updated_at = NOW() WHERE split_id = $3`, status, settledAt, id)
	return err
}

func (r *expenseSplitRepo) FindOverdue(ctx context.Context) ([]domain.ExpenseSplit, error) {
	query := `
		SELECT s.split_id, s.transaction_id, s.user_id, s.owe_amount, s.payment_status, s.payment_method, s.payslip_url, s.approval_status, s.settled_at, s.updated_at 
		FROM expense_splits s 
		JOIN expense_transactions t ON s.transaction_id = t.transaction_id 
		WHERE s.payment_status IN ('Pending', 'Submitted') AND t.due_date < NOW()`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var splits []domain.ExpenseSplit
	for rows.Next() {
		var s domain.ExpenseSplit
		if err := rows.Scan(&s.SplitID, &s.TransactionID, &s.UserID, &s.OweAmount, &s.PaymentStatus, &s.PaymentMethod, &s.PayslipURL, &s.ApprovalStatus, &s.SettledAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		splits = append(splits, s)
	}
	return splits, nil
}

func (r *expenseSplitRepo) CountUnsettled(ctx context.Context, transactionID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM expense_splits WHERE transaction_id = $1 AND payment_status != 'Approved'`, transactionID).Scan(&count)
	return count, err
}
