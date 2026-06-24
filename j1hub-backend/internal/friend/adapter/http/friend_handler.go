package http

import (
	"github.com/j1hub/backend/internal/friend/adapter/http/dto"
	frienddomain "github.com/j1hub/backend/internal/friend/domain"

	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/transport/http/middleware"
	gamificationusecase "github.com/j1hub/backend/internal/gamification/usecase"
	"github.com/j1hub/backend/pkg/apperror"
)

type FriendshipUC interface {
	SendRequest(ctx context.Context, senderID, targetID string) error
	ListPendingRequests(ctx context.Context, userID string) ([]frienddomain.Friendship, error)
	RespondToRequest(ctx context.Context, userID, friendshipID string, accept bool) error
	ListFriends(ctx context.Context, userID string) ([]frienddomain.Friendship, error)
	RemoveFriend(ctx context.Context, userID, friendID string) error
}

type RadarUC interface {
	GetRadar(ctx context.Context, requesterID string) ([]gamificationusecase.RadarEntry, error)
}

type FriendHandler struct {
	friendshipUC FriendshipUC
	radarUC      RadarUC
}

func NewFriendHandler(friendshipUC FriendshipUC, radarUC RadarUC) *FriendHandler {
	log.Println("debugprint: entering NewFriendHandler")
	return &FriendHandler{friendshipUC: friendshipUC, radarUC: radarUC}
}

func (h *FriendHandler) SendRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*FriendHandler).SendRequest")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	var req dto.FriendRequestReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
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
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	requests, err := h.friendshipUC.ListPendingRequests(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	page, pageSize := apperror.ParsePagination(r)
	apperror.RespondList(w, requests, page, pageSize, len(requests))
}

func (h *FriendHandler) RespondToRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*FriendHandler).RespondToRequest")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	var req dto.RespondFriendReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	accept := false
	if req.Accept != nil {
		accept = *req.Accept
	} else if req.Status == "Accepted" {
		accept = true
	}

	err := h.friendshipUC.RespondToRequest(r.Context(), claims.UserID, id, accept)
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
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	friends, err := h.friendshipUC.ListFriends(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	page, pageSize := apperror.ParsePagination(r)
	apperror.RespondList(w, friends, page, pageSize, len(friends))
}

func (h *FriendHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*FriendHandler).RemoveFriend")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "friendshipId")
	if id == "" {
		id = chi.URLParam(r, "id")
	}

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
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	radar, err := h.radarUC.GetRadar(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	page, pageSize := apperror.ParsePagination(r)
	apperror.RespondList(w, radar, page, pageSize, len(radar))
}
