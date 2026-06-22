package dto

import "github.com/j1hub/backend/internal/domain"

type ExpenseDetailResponse struct {
	Transaction *domain.ExpenseTransaction `json:"transaction"`
	Splits      []domain.ExpenseSplit      `json:"splits"`
}

func NewExpenseDetailResponse(txn *domain.ExpenseTransaction, splits []domain.ExpenseSplit) *ExpenseDetailResponse {
	return &ExpenseDetailResponse{
		Transaction: txn,
		Splits:      splits,
	}
}
