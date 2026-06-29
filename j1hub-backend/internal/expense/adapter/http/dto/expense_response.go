package dto

import (
	expensedomain "github.com/j1hub/backend/internal/expense/domain"
)

type ExpenseDetailResponse struct {
	Transaction *expensedomain.ExpenseTransaction `json:"transaction"`
	Splits      []string                          `json:"splits"`
}

func NewExpenseDetailResponse(txn *expensedomain.ExpenseTransaction, splits []expensedomain.ExpenseSplit) *ExpenseDetailResponse {
	splitIDs := make([]string, len(splits))
	for i, s := range splits {
		splitIDs[i] = s.SplitID
	}
	return &ExpenseDetailResponse{
		Transaction: txn,
		Splits:      splitIDs,
	}
}
