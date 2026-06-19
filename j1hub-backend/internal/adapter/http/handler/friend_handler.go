package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/j1hub/backend/pkg/apperror"
)

type FriendshipUC interface {
	SendRequest(ctx context.Context, senderID, targetID string) error
	ListPendingRequests(ctx context.Context, userID string) ([]domain.Friendship, error)
	RespondToRequest(ctx context.Context, userID, friendshipID string, accept bool) error
	ListFriends(ctx context.Context, userID string) ([]domain.Friendship, error)
	RemoveFriend(ctx context.Context, userID, friendID string) error
}

type RadarUC interface {
	GetRadar(ctx context.Context, requesterID string) ([]usecase.RadarEntry, error)
}

type FriendHandler struct {
	friendshipUC FriendshipUC
	radarUC      RadarUC
}

func NewFriendHandler(friendshipUC FriendshipUC, radarUC RadarUC) *FriendHandler {
	log.Println("debugprint: entering NewFriendHandler")
	return &FriendHandler{friendshipUC: friendshipUC, radarUC: radarUC}
}

type friendRequestReq struct {
	TargetUserID string `json:"target_user_id" validate:"required"`
}

func (h *FriendHandler) SendRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*FriendHandler).SendRequest")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var req friendRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Invalid request body", Err: err})
		return
	}

	err := h.friendshipUC.SendRequest(r.Context(), claims.UserID, req.TargetUserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *FriendHandler) ListPendingRequests(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*FriendHandler).ListPendingRequests")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	requests, err := h.friendshipUC.ListPendingRequests(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	json.NewEncoder(w).Encode(requests)
}

type respondFriendReq struct {
	FriendshipID string `json:"friendship_id" validate:"required"`
	Accept       bool   `json:"accept"`
}

func (h *FriendHandler) RespondToRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*FriendHandler).RespondToRequest")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var req respondFriendReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Invalid request body", Err: err})
		return
	}

	err := h.friendshipUC.RespondToRequest(r.Context(), claims.UserID, req.FriendshipID, req.Accept)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FriendHandler) ListFriends(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*FriendHandler).ListFriends")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	friends, err := h.friendshipUC.ListFriends(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	json.NewEncoder(w).Encode(friends)
}

func (h *FriendHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*FriendHandler).RemoveFriend")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	err := h.friendshipUC.RemoveFriend(r.Context(), claims.UserID, id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FriendHandler) GetRadar(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*FriendHandler).GetRadar")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	radar, err := h.radarUC.GetRadar(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	json.NewEncoder(w).Encode(radar)
}
