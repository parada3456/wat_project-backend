package domain

import (
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
	switch s {
	case ApprovalPending, ApprovalApproved, ApprovalRejected:
		return true
	}
	return false
}

type ExpenseTransaction struct {
	TransactionID   string
	PaidByUserID    string
	Title           string
	TotalAmount     float64
	Currency        string
	Memo            string
	TransactionDate time.Time
	DueDate         time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ExpenseSplit struct {
	SplitID        string
	TransactionID  string
	UserID         string
	OweAmount      float64
	PaymentStatus  PaymentStatus
	PaymentMethod  string
	PayslipURL     string
	ApprovalStatus ApprovalStatus
	SettledAt      *time.Time
	UpdatedAt      time.Time
}

func (s *ExpenseSplit) IsSettled() bool {
	return s.PaymentStatus == PaymentApproved
}
