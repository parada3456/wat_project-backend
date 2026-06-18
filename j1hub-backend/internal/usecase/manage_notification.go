package usecase

import (
	"context"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
)

type NotificationUseCase struct {
	notifRepo port.NotificationRepository
}

func NewNotificationUseCase(notifRepo port.NotificationRepository) *NotificationUseCase {
	return &NotificationUseCase{notifRepo: notifRepo}
}

func (uc *NotificationUseCase) ListNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	return uc.notifRepo.FindByUser(ctx, userID)
}

func (uc *NotificationUseCase) MarkRead(ctx context.Context, id string) error {
	return uc.notifRepo.MarkAsRead(ctx, id)
}

func (uc *NotificationUseCase) MarkAllRead(ctx context.Context, userID string) error {
	return uc.notifRepo.MarkAllAsRead(ctx, userID)
}

func (uc *NotificationUseCase) Delete(ctx context.Context, id string) error {
	return uc.notifRepo.Delete(ctx, id)
}
