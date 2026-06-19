package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/pkg/apperror"
)

type NotificationUC interface {
	ListNotifications(ctx context.Context, userID string) ([]domain.Notification, error)
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
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	notifs, err := h.notifUC.ListNotifications(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	json.NewEncoder(w).Encode(notifs)
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
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
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
