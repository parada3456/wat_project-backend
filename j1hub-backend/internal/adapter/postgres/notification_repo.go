package postgres

import (
	"context"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/jackc/pgx/v5/pgxpool"
)

type notificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) port.NotificationRepository {
	return &notificationRepository{pool: pool}
}

func (r *notificationRepository) Insert(ctx context.Context, n *domain.Notification) error {
	// Implementation
	return nil
}

func (r *notificationRepository) FindByUser(ctx context.Context, userID string) ([]domain.Notification, error) {
	// Implementation
	return nil, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id string) error {
	// Implementation
	return nil
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	// Implementation
	return nil
}

func (r *notificationRepository) Delete(ctx context.Context, id string) error {
	// Implementation
	return nil
}
