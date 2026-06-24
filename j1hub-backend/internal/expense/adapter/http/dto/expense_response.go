package dto

import (
	expensedomain "github.com/j1hub/backend/internal/expense/domain"
)

type ExpenseDetailResponse struct {
	Transaction *expensedomain.ExpenseTransaction `json:"transaction"`
	Splits      []expensedomain.ExpenseSplit      `json:"splits"`
}

func NewExpenseDetailResponse(txn *expensedomain.ExpenseTransaction, splits []expensedomain.ExpenseSplit) *ExpenseDetailResponse {
	return &ExpenseDetailResponse{
		Transaction: txn,
		Splits:      splits,
	}
}
