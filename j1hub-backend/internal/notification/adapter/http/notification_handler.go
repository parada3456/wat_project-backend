package http

import (
	"context"
	"log"
	"net/http"
	"strconv"

	notificationdomain "github.com/parada3456/wat_project-backend/internal/notification/domain"

	"github.com/go-chi/chi/v5"
	"github.com/parada3456/wat_project-backend/internal/domain"
	"github.com/parada3456/wat_project-backend/internal/transport/http/middleware"
	"github.com/parada3456/wat_project-backend/pkg/apperror"
)

type NotificationUC interface {
	ListNotifications(ctx context.Context, userID string, isRead *bool, page, pageSize int) ([]notificationdomain.Notification, int, error)
	MarkRead(ctx context.Context, id string) error
	MarkAllRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
}

type NotificationHandler struct {
	notifUC NotificationUC
}

func NewNotificationHandler(notifUC NotificationUC) *NotificationHandler {
	log.Println("debugprint: entering NewNotificationHandler")
	return &NotificationHandler{notifUC: notifUC}
}

func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*NotificationHandler).ListNotifications")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	pago := apperror.ParsePagination(r)

	var isReadFilter *bool
	isReadStr := r.URL.Query().Get("isRead")
	if isReadStr == "" {
		isReadStr = r.URL.Query().Get("is_read")
	}
	if isReadStr != "" {
		if isRead, err := strconv.ParseBool(isReadStr); err == nil {
			isReadFilter = &isRead
		}
	}

	notifs, totalCount, err := h.notifUC.ListNotifications(r.Context(), claims.UserID, isReadFilter, pago.Page, pago.PageSize)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	apperror.RespondList(w, notifs, pago.Page, pago.PageSize, totalCount)
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*NotificationHandler).MarkRead")
	id := chi.URLParam(r, "id")
	err := h.notifUC.MarkRead(r.Context(), id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*NotificationHandler).MarkAllRead")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	err := h.notifUC.MarkAllRead(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NotificationHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*NotificationHandler).DeleteNotification")
	id := chi.URLParam(r, "id")
	err := h.notifUC.Delete(r.Context(), id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
