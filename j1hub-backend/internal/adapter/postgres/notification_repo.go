package postgres

import (
	"context"
	"log"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/jackc/pgx/v5/pgxpool"
)

type notificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) port.NotificationRepository {
	log.Println("debugprint: entering NewNotificationRepository")
	return &notificationRepository{pool: pool}
}

func (r *notificationRepository) Insert(ctx context.Context, n *domain.Notification) error {
	log.
		// Implementation
		Println("debugprint: entering (*notificationRepository).Insert")

	return nil
}

func (r *notificationRepository) FindByUser(ctx context.Context, userID string) ([]domain.Notification, error) {
	log.
		// Implementation
		Println("debugprint: entering (*notificationRepository).FindByUser")

	return nil, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id string) error {
	log.
		// Implementation
		Println("debugprint: entering (*notificationRepository).MarkAsRead")

	return nil
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	log.
		// Implementation
		Println("debugprint: entering (*notificationRepository).MarkAllAsRead")

	return nil
}

func (r *notificationRepository) Delete(ctx context.Context, id string) error {
	log.
		// Implementation
		Println("debugprint: entering (*notificationRepository).Delete")

	return nil
}
