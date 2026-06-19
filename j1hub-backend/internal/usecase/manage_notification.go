package usecase

import (
	"context"
	"log"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
)

type NotificationUseCase struct {
	notifRepo port.NotificationRepository
}

func NewNotificationUseCase(notifRepo port.NotificationRepository) *NotificationUseCase {
	log.Println("debugprint: entering NewNotificationUseCase")
	return &NotificationUseCase{notifRepo: notifRepo}
}

func (uc *NotificationUseCase) ListNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	log.Println("debugprint: entering (*NotificationUseCase).ListNotifications")
	return uc.notifRepo.FindByUser(ctx, userID)
}

func (uc *NotificationUseCase) MarkRead(ctx context.Context, id string) error {
	log.Println("debugprint: entering (*NotificationUseCase).MarkRead")
	return uc.notifRepo.MarkAsRead(ctx, id)
}

func (uc *NotificationUseCase) MarkAllRead(ctx context.Context, userID string) error {
	log.Println("debugprint: entering (*NotificationUseCase).MarkAllRead")
	return uc.notifRepo.MarkAllAsRead(ctx, userID)
}

func (uc *NotificationUseCase) Delete(ctx context.Context, id string) error {
	log.Println("debugprint: entering (*NotificationUseCase).Delete")
	return uc.notifRepo.Delete(ctx, id)
}
