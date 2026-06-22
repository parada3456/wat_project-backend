package handler

import (
	"github.com/j1hub/backend/internal/adapter/http/handler/dto"

	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/pkg/apperror"
)

type AdminHandler struct {
	adminUseCase port.AdminUseCase
	validate     *validator.Validate
}

func NewAdminHandler(uc port.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUseCase: uc,
		validate:     validator.New(),
	}
}

func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.adminUseCase.GetDashboardStats(r.Context())
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func (h *AdminHandler) ListPendingVerifications(w http.ResponseWriter, r *http.Request) {
	ums, err := h.adminUseCase.ListPendingVerifications(r.Context())
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	page, pageSize := parsePagination(r)
	apperror.RespondList(w, ums, page, pageSize, len(ums))
}

func (h *AdminHandler) VerifyMission(w http.ResponseWriter, r *http.Request) {
	userMissionID := chi.URLParam(r, "id")
	claims := middleware.GetClaims(r.Context())
	if claims == nil || !claims.IsAdmin {
		apperror.RespondError(w, domain.ErrForbidden)
		return
	}

	var req dto.VerifyMissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Malformed request body", Err: err})
		return
	}

	um, err := h.adminUseCase.VerifyMission(r.Context(), claims.UserID, userMissionID, req.Approved, req.RejectionReason)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respDTO := dto.NewVerifyMissionResponse(um)
	json.NewEncoder(w).Encode(respDTO)
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")
	users, err := h.adminUseCase.ListUsers(r.Context(), searchQuery)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	page, pageSize := parsePagination(r)
	apperror.RespondList(w, users, page, pageSize, len(users))
}

func (h *AdminHandler) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	user, err := h.adminUseCase.GetUserDetail(r.Context(), userID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *AdminHandler) AdjustPoints(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var req dto.AdjustPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Malformed request body", Err: err})
		return
	}

	res, err := h.adminUseCase.AdjustPoints(r.Context(), userID, req.PointsDelta, req.Reason)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
