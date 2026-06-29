package notificationusecase

import (
	"context"
	"log"

	notificationdomain "github.com/j1hub/backend/internal/notification/domain"

	port "github.com/j1hub/backend/internal/notification/port"
)

type NotificationUseCase struct {
	notifRepo port.NotificationRepository
}

func NewNotificationUseCase(notifRepo port.NotificationRepository) *NotificationUseCase {
	log.Println("debugprint: entering NewNotificationUseCase")
	return &NotificationUseCase{notifRepo: notifRepo}
}

func (uc *NotificationUseCase) ListNotifications(
	ctx context.Context,
	userID string,
	isRead *bool,
	page,
	pageSize int,
) ([]notificationdomain.Notification, int, error) {
	log.Println("debugprint: entering (*NotificationUseCase).ListNotifications")

	// 1. Calculate database boundary variables
	limit := pageSize
	offset := (page - 1) * pageSize

	// 2. Pass everything down to your updated repository method
	return uc.notifRepo.FindByUser(ctx, userID, isRead, limit, offset)
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
