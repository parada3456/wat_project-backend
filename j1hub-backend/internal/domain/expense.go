package domain

import (
	"log"
	"time"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "Pending"
	PaymentSubmitted PaymentStatus = "Submitted"
	PaymentApproved  PaymentStatus = "Approved"
	PaymentOverdue   PaymentStatus = "Overdue"
)

func (s PaymentStatus) Valid() bool {
	log.Println("debugprint: entering (PaymentStatus).Valid")
	switch s {
	case PaymentPending, PaymentSubmitted, PaymentApproved, PaymentOverdue:
		return true
	}
	return false
}

type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "Pending_Approval"
	ApprovalApproved ApprovalStatus = "Approved"
	ApprovalRejected ApprovalStatus = "Rejected"
)

func (s ApprovalStatus) Valid() bool {
	log.Println("debugprint: entering (ApprovalStatus).Valid")
	switch s {
	case ApprovalPending, ApprovalApproved, ApprovalRejected:
		return true
	}
	return false
}

type ExpenseTransaction struct {
	TransactionID   string    `json:"transaction_id"`
	PaidByUserID    string    `json:"paid_by_user_id"`
	Title           string    `json:"title"`
	TotalAmount     float64   `json:"total_amount"`
	Currency        string    `json:"currency"`
	Memo            string    `json:"memo"`
	TransactionDate time.Time `json:"transaction_date"`
	DueDate         time.Time `json:"due_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ExpenseSplit struct {
	SplitID        string         `json:"split_id"`
	TransactionID  string         `json:"transaction_id"`
	UserID         string         `json:"user_id"`
	OweAmount      float64        `json:"owe_amount"`
	PaymentStatus  PaymentStatus  `json:"payment_status"`
	PaymentMethod  string         `json:"payment_method"`
	PayslipURL     string         `json:"payslip_url"`
	ApprovalStatus ApprovalStatus `json:"approval_status"`
	SettledAt      *time.Time     `json:"settled_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (s *ExpenseSplit) IsSettled() bool {
	log.Println("debugprint: entering (*ExpenseSplit).IsSettled")
	return s.PaymentStatus == PaymentApproved
}
