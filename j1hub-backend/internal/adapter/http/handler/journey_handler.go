package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"

	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/j1hub/backend/pkg/apperror"
)

type JourneyUC interface {
	ListPhases(ctx context.Context) ([]gamificationdomain.JourneyPhase, error)
	GetHistory(ctx context.Context, userID string) ([]gamificationdomain.UserPhaseHistory, error)
	ListUserBadges(ctx context.Context, userID string) ([]gamificationdomain.UserBadge, error)
	GetCreditScoreHistory(ctx context.Context, userID string) ([]gamificationdomain.PointLedger, error)
	GetPointsLedger(ctx context.Context, userID string) ([]gamificationdomain.PointLedger, error)
}

type AdvancePhaseUC interface {
	TryAdvancePhase(ctx context.Context, userID string) (*usecase.PhaseTransitionResponse, error)
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

func parsePagination(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	return page, pageSize
}

func (h *JourneyHandler) ListPhases(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JourneyHandler).ListPhases")
	phases, err := h.journeyUC.ListPhases(r.Context())
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	page, pageSize := parsePagination(r)
	apperror.RespondList(w, phases, page, pageSize, len(phases))
}

func (h *JourneyHandler) AdvancePhase(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JourneyHandler).AdvancePhase")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	resp, err := h.advanceUC.TryAdvancePhase(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
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
	page, pageSize := parsePagination(r)
	apperror.RespondList(w, history, page, pageSize, len(history))
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
	page, pageSize := parsePagination(r)
	apperror.RespondList(w, entries, page, pageSize, len(entries))
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
	page, pageSize := parsePagination(r)
	apperror.RespondList(w, badges, page, pageSize, len(badges))
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
	page, pageSize := parsePagination(r)
	apperror.RespondList(w, history, page, pageSize, len(history))
}

func (h *JourneyHandler) GetPointsLedger(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JourneyHandler).GetPointsLedger")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	ledger, err := h.journeyUC.GetPointsLedger(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	page, pageSize := parsePagination(r)
	apperror.RespondList(w, ledger, page, pageSize, len(ledger))
}
