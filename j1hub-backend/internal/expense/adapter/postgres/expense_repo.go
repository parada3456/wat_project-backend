package postgres

import (
	"context"
	"log"
	"time"

	"github.com/j1hub/backend/internal/domain"
	expensedomain "github.com/j1hub/backend/internal/expense/domain"
	port "github.com/j1hub/backend/internal/expense/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type expenseRepo struct {
	pool *pgxpool.Pool
}

func NewExpenseRepository(pool *pgxpool.Pool) port.ExpenseTransactionRepository {
	log.Println("debugprint: entering NewExpenseRepository")
	return &expenseRepo{pool: pool}
}

func (r *expenseRepo) Insert(ctx context.Context, t *expensedomain.ExpenseTransaction) error {
	log.Println("debugprint: entering (*expenseRepo).Insert")
	_, err := r.pool.Exec(ctx, `INSERT INTO expense_transactions (transaction_id, paid_by_user_id, title, total_amount, currency, memo, transaction_date, due_date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		t.TransactionID, t.PaidByUserID, t.Title, t.TotalAmount, t.Currency, t.Memo, t.TransactionDate, t.DueDate, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *expenseRepo) FindByID(ctx context.Context, id string) (*expensedomain.ExpenseTransaction, error) {
	log.Println("debugprint: entering (*expenseRepo).FindByID")
	var t expensedomain.ExpenseTransaction
	err := r.pool.QueryRow(ctx, `SELECT transaction_id, paid_by_user_id, title, total_amount, currency, memo, transaction_date, due_date, created_at, updated_at FROM expense_transactions WHERE transaction_id = $1`, id).Scan(&t.TransactionID, &t.PaidByUserID, &t.Title, &t.TotalAmount, &t.Currency, &t.Memo, &t.TransactionDate, &t.DueDate, &t.CreatedAt, &t.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &t, err
}

func (r *expenseRepo) MarkSettled(ctx context.Context, id string) error {
	log.Println("debugprint: entering (*expenseRepo).MarkSettled")
	_, err := r.pool.Exec(ctx, `UPDATE expense_transactions SET updated_at = NOW() WHERE transaction_id = $1`, id)
	return err
}

func (r *expenseRepo) FindByUser(ctx context.Context, userID string, limit, offset int) ([]expensedomain.ExpenseTransaction, int, error) {
	log.Println("debugprint: entering (*expenseRepo).FindByUser")
	
	var totalCount int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM expense_transactions WHERE paid_by_user_id = $1`, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []expensedomain.ExpenseTransaction{}, 0, nil
	}

	query := `SELECT transaction_id, paid_by_user_id, title, total_amount, currency, memo, transaction_date, due_date, created_at, updated_at FROM expense_transactions WHERE paid_by_user_id = $1 ORDER BY transaction_date DESC LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var txns []expensedomain.ExpenseTransaction
	for rows.Next() {
		var t expensedomain.ExpenseTransaction
		if err := rows.Scan(&t.TransactionID, &t.PaidByUserID, &t.Title, &t.TotalAmount, &t.Currency, &t.Memo, &t.TransactionDate, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		txns = append(txns, t)
	}
	return txns, totalCount, nil
}

func (r *expenseRepo) Delete(ctx context.Context, id string) error {
	log.Println("debugprint: entering (*expenseRepo).Delete")
	_, err := r.pool.Exec(ctx, `DELETE FROM expense_transactions WHERE transaction_id = $1`, id)
	return err
}

type expenseSplitRepo struct {
	pool *pgxpool.Pool
}

func NewExpenseSplitRepository(pool *pgxpool.Pool) port.ExpenseSplitRepository {
	log.Println("debugprint: entering NewExpenseSplitRepository")
	return &expenseSplitRepo{pool: pool}
}

func (r *expenseSplitRepo) BulkInsert(ctx context.Context, splits []expensedomain.ExpenseSplit) error {
	log.Println("debugprint: entering (*expenseSplitRepo).BulkInsert")
	batch := &pgx.Batch{}
	for _, s := range splits {
		batch.Queue(`INSERT INTO expense_splits (split_id, transaction_id, user_id, owe_amount, payment_status, payment_method, payslip_url, approval_status, settled_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			s.SplitID, s.TransactionID, s.UserID, s.OweAmount, s.PaymentStatus, s.PaymentMethod, s.PayslipURL, s.ApprovalStatus, s.SettledAt, s.UpdatedAt)
	}
	return r.pool.SendBatch(ctx, batch).Close()
}

func (r *expenseSplitRepo) FindByID(ctx context.Context, id string) (*expensedomain.ExpenseSplit, error) {
	log.Println("debugprint: entering (*expenseSplitRepo).FindByID")
	var s expensedomain.ExpenseSplit
	err := r.pool.QueryRow(ctx, `SELECT split_id, transaction_id, user_id, owe_amount, payment_status, payment_method, payslip_url, approval_status, settled_at, updated_at FROM expense_splits WHERE split_id = $1`, id).Scan(&s.SplitID, &s.TransactionID, &s.UserID, &s.OweAmount, &s.PaymentStatus, &s.PaymentMethod, &s.PayslipURL, &s.ApprovalStatus, &s.SettledAt, &s.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return &s, err
}

func (r *expenseSplitRepo) UpdatePaymentStatus(ctx context.Context, id string, status expensedomain.PaymentStatus, slipURL string) error {
	log.Println("debugprint: entering (*expenseSplitRepo).UpdatePaymentStatus")
	_, err := r.pool.Exec(ctx, `UPDATE expense_splits SET payment_status = $1, payslip_url = $2, updated_at = NOW() WHERE split_id = $3`, status, slipURL, id)
	return err
}

func (r *expenseSplitRepo) UpdateApproval(ctx context.Context, id string, status expensedomain.ApprovalStatus, settledAt *time.Time) error {
	log.Println("debugprint: entering (*expenseSplitRepo).UpdateApproval")
	_, err := r.pool.Exec(ctx, `UPDATE expense_splits SET approval_status = $1, payment_status = CASE WHEN $1 = 'approved' THEN 'approved' ELSE payment_status END, settled_at = $2, updated_at = NOW() WHERE split_id = $3`, status, settledAt, id)
	return err
}

func (r *expenseSplitRepo) FindOverdue(ctx context.Context) ([]expensedomain.ExpenseSplit, error) {
	log.Println("debugprint: entering (*expenseSplitRepo).FindOverdue")
	query := `
		SELECT s.split_id, s.transaction_id, s.user_id, s.owe_amount, s.payment_status, s.payment_method, s.payslip_url, s.approval_status, s.settled_at, s.updated_at 
		FROM expense_splits s 
		JOIN expense_transactions t ON s.transaction_id = t.transaction_id 
		WHERE s.payment_status IN ('pending', 'submitted') AND t.due_date < NOW()`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var splits []expensedomain.ExpenseSplit
	for rows.Next() {
		var s expensedomain.ExpenseSplit
		if err := rows.Scan(&s.SplitID, &s.TransactionID, &s.UserID, &s.OweAmount, &s.PaymentStatus, &s.PaymentMethod, &s.PayslipURL, &s.ApprovalStatus, &s.SettledAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		splits = append(splits, s)
	}
	return splits, nil
}

func (r *expenseSplitRepo) CountUnsettled(ctx context.Context, transactionID string) (int, error) {
	log.Println("debugprint: entering (*expenseSplitRepo).CountUnsettled")
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM expense_splits WHERE transaction_id = $1 AND payment_status != 'approved'`, transactionID).Scan(&count)
	return count, err
}

func (r *expenseSplitRepo) FindByUser(ctx context.Context, userID string, limit, offset int) ([]expensedomain.ExpenseSplit, int, error) {
	log.Println("debugprint: entering (*expenseSplitRepo).FindByUser")
	
	var totalCount int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM expense_splits WHERE user_id = $1`, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []expensedomain.ExpenseSplit{}, 0, nil
	}

	query := `SELECT split_id, transaction_id, user_id, owe_amount, payment_status, payment_method, payslip_url, approval_status, settled_at, updated_at FROM expense_splits WHERE user_id = $1 LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var splits []expensedomain.ExpenseSplit
	for rows.Next() {
		var s expensedomain.ExpenseSplit
		if err := rows.Scan(&s.SplitID, &s.TransactionID, &s.UserID, &s.OweAmount, &s.PaymentStatus, &s.PaymentMethod, &s.PayslipURL, &s.ApprovalStatus, &s.SettledAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		splits = append(splits, s)
	}
	return splits, totalCount, nil
}

func (r *expenseSplitRepo) FindPendingByUser(ctx context.Context, userID string, limit, offset int) ([]expensedomain.ExpenseSplit, int, error) {
	log.Println("debugprint: entering (*expenseSplitRepo).FindPendingByUser")
	
	var totalCount int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM expense_splits WHERE user_id = $1 AND payment_status != 'approved'`, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []expensedomain.ExpenseSplit{}, 0, nil
	}

	query := `SELECT split_id, transaction_id, user_id, owe_amount, payment_status, payment_method, payslip_url, approval_status, settled_at, updated_at FROM expense_splits WHERE user_id = $1 AND payment_status != 'approved' LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var splits []expensedomain.ExpenseSplit
	for rows.Next() {
		var s expensedomain.ExpenseSplit
		if err := rows.Scan(&s.SplitID, &s.TransactionID, &s.UserID, &s.OweAmount, &s.PaymentStatus, &s.PaymentMethod, &s.PayslipURL, &s.ApprovalStatus, &s.SettledAt, &s.UpdatedAt); err != nil {
			return nil, 0, err
		}
		splits = append(splits, s)
	}
	return splits, totalCount, nil
}

func (r *expenseSplitRepo) FindByTransactionID(ctx context.Context, transactionID string) ([]expensedomain.ExpenseSplit, error) {
	log.Println("debugprint: entering (*expenseSplitRepo).FindByTransactionID")
	query := `SELECT split_id, transaction_id, user_id, owe_amount, payment_status, payment_method, payslip_url, approval_status, settled_at, updated_at FROM expense_splits WHERE transaction_id = $1`
	rows, err := r.pool.Query(ctx, query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var splits []expensedomain.ExpenseSplit
	for rows.Next() {
		var s expensedomain.ExpenseSplit
		if err := rows.Scan(&s.SplitID, &s.TransactionID, &s.UserID, &s.OweAmount, &s.PaymentStatus, &s.PaymentMethod, &s.PayslipURL, &s.ApprovalStatus, &s.SettledAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		splits = append(splits, s)
	}
	return splits, nil
}

func (r *expenseSplitRepo) DeleteByTransactionID(ctx context.Context, transactionID string) error {
	log.Println("debugprint: entering (*expenseSplitRepo).DeleteByTransactionID")
	_, err := r.pool.Exec(ctx, `DELETE FROM expense_splits WHERE transaction_id = $1`, transactionID)
	return err
}
