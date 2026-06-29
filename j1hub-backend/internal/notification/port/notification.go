package port

import (
	"context"
	notificationdomain "github.com/j1hub/backend/internal/notification/domain"
)

type NotificationRepository interface {
	Insert(ctx context.Context, n *notificationdomain.Notification) error
	FindByUser(ctx context.Context, userID string, isRead *bool, limit, offset int) ([]notificationdomain.Notification, int, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
}

type NotifierPort interface {
	Send(ctx context.Context, userID, title, body string) error
}
