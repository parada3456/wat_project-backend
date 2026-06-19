package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/j1hub/backend/pkg/apperror"
)

type JourneyUC interface {
	ListPhases(ctx context.Context) ([]domain.JourneyPhase, error)
	GetHistory(ctx context.Context, userID string) ([]domain.UserPhaseHistory, error)
	ListUserBadges(ctx context.Context, userID string) ([]domain.UserBadge, error)
	GetCreditScoreHistory(ctx context.Context, userID string) ([]domain.PointLedger, error)
}

type AdvancePhaseUC interface {
	TryAdvancePhase(ctx context.Context, userID string) (bool, error)
}

type LeaderboardUC interface {
	GetLeaderboard(ctx context.Context, scope, jobID string) ([]usecase.LeaderboardEntry, error)
}

type JourneyHandler struct {
	journeyUC     JourneyUC
	advanceUC     AdvancePhaseUC
	leaderboardUC LeaderboardUC
}

func NewJourneyHandler(
	journeyUC JourneyUC,
	advanceUC AdvancePhaseUC,
	leaderboardUC LeaderboardUC,
) *JourneyHandler {
	log.Println("debugprint: entering NewJourneyHandler")
	return &JourneyHandler{
		journeyUC:     journeyUC,
		advanceUC:     advanceUC,
		leaderboardUC: leaderboardUC,
	}
}

func (h *JourneyHandler) ListPhases(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JourneyHandler).ListPhases")
	phases, err := h.journeyUC.ListPhases(r.Context())
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	json.NewEncoder(w).Encode(phases)
}

func (h *JourneyHandler) AdvancePhase(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JourneyHandler).AdvancePhase")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	success, err := h.advanceUC.TryAdvancePhase(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": success})
}

func (h *JourneyHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JourneyHandler).GetHistory")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	history, err := h.journeyUC.GetHistory(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	json.NewEncoder(w).Encode(history)
}

func (h *JourneyHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JourneyHandler).GetLeaderboard")
	scope := r.URL.Query().Get("scope")
	jobID := r.URL.Query().Get("job_id")

	entries, err := h.leaderboardUC.GetLeaderboard(r.Context(), scope, jobID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	json.NewEncoder(w).Encode(entries)
}

func (h *JourneyHandler) ListBadges(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JourneyHandler).ListBadges")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	badges, err := h.journeyUC.ListUserBadges(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	json.NewEncoder(w).Encode(badges)
}

func (h *JourneyHandler) GetCreditHistory(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JourneyHandler).GetCreditHistory")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	history, err := h.journeyUC.GetCreditScoreHistory(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	json.NewEncoder(w).Encode(history)
}
