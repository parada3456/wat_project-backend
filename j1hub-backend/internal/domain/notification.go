package domain

import "time"

type Notification struct {
	NotificationID string
	UserID         string
	Title          string
	Body           string
	IsRead         bool
	CreatedAt      time.Time
}
