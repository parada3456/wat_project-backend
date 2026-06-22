package dto

import "time"

type CreateExpenseReq struct {
	Title       string    `json:"title" validate:"required"`
	TotalAmount float64   `json:"total_amount" validate:"required,gt=0"`
	Currency    string    `json:"currency" validate:"required"`
	Memo        string    `json:"memo"`
	DueDate     time.Time `json:"due_date" validate:"required"`
	Splits      []struct {
		UserID    string  `json:"user_id" validate:"required"`
		OweAmount float64 `json:"owe_amount" validate:"required,gt=0"`
	} `json:"splits" validate:"required,dive"`
}

