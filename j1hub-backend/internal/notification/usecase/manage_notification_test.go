package notificationusecase_test

import (
	"context"
	"testing"

	notificationusecase "github.com/j1hub/backend/internal/notification/usecase"

	notificationdomain "github.com/j1hub/backend/internal/notification/domain"

	"github.com/stretchr/testify/assert"
)

func TestNotificationUseCase_ListNotifications_Success(t *testing.T) {
	notifRepo := new(MockNotificationRepository)
	uc := notificationusecase.NewNotificationUseCase(notifRepo)

	ctx := context.Background()
	userID := "usr_123"
	mockNotifs := []notificationdomain.Notification{
		{NotificationID: "notif_1", UserID: userID, Title: "Test Title", Body: "Test Body"},
	}

	notifRepo.On("FindByUser", ctx, userID).Return(mockNotifs, nil)

	res, err := uc.ListNotifications(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockNotifs, res)
}

func TestNotificationUseCase_MarkRead_Success(t *testing.T) {
	notifRepo := new(MockNotificationRepository)
	uc := notificationusecase.NewNotificationUseCase(notifRepo)

	ctx := context.Background()
	id := "notif_1"

	notifRepo.On("MarkAsRead", ctx, id).Return(nil)

	err := uc.MarkRead(ctx, id)

	assert.NoError(t, err)
}

func TestNotificationUseCase_MarkAllRead_Success(t *testing.T) {
	notifRepo := new(MockNotificationRepository)
	uc := notificationusecase.NewNotificationUseCase(notifRepo)

	ctx := context.Background()
	userID := "usr_123"

	notifRepo.On("MarkAllAsRead", ctx, userID).Return(nil)

	err := uc.MarkAllRead(ctx, userID)

	assert.NoError(t, err)
}

func TestNotificationUseCase_Delete_Success(t *testing.T) {
	notifRepo := new(MockNotificationRepository)
	uc := notificationusecase.NewNotificationUseCase(notifRepo)

	ctx := context.Background()
	id := "notif_1"

	notifRepo.On("Delete", ctx, id).Return(nil)

	err := uc.Delete(ctx, id)

	assert.NoError(t, err)
}
