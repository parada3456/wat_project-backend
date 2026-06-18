package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/j1hub/backend/pkg/apperror"
)

type MissionUC interface {
	ListAvailableMissions(ctx context.Context, userID string) ([]domain.UserMission, error)
	GetMissionDetail(ctx context.Context, userID, userMissionID string) (*usecase.MissionDetailResponse, error)
	ToggleTask(ctx context.Context, userID, userTaskID string, completed bool) error
}

type CompleteMissionUC interface {
	SubmitProof(ctx context.Context, userID, userMissionID string, file io.Reader, contentType string) error
}

type MissionHandler struct {
	missionUC  MissionUC
	completeUC CompleteMissionUC
}

func NewMissionHandler(missionUC MissionUC, completeUC CompleteMissionUC) *MissionHandler {
	return &MissionHandler{
		missionUC:  missionUC,
		completeUC: completeUC,
	}
}

func (h *MissionHandler) ListMissions(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	missions, err := h.missionUC.ListAvailableMissions(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	json.NewEncoder(w).Encode(missions)
}

func (h *MissionHandler) GetMissionDetail(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	detail, err := h.missionUC.GetMissionDetail(r.Context(), claims.UserID, id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	json.NewEncoder(w).Encode(detail)
}

func (h *MissionHandler) SubmitProof(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	
	// Need to handle multipart form for file upload
	// For simplicity, let's assume JSON for now or handle multipart
	r.ParseMultipartForm(10 << 20) // 10MB max
	file, header, err := r.FormFile("proof")
	if err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Proof file required", Err: err})
		return
	}
	defer file.Close()

	err = h.completeUC.SubmitProof(r.Context(), claims.UserID, id, file, header.Header.Get("Content-Type"))
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type toggleTaskReq struct {
	Completed bool `json:"completed"`
}

func (h *MissionHandler) ToggleTask(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	var req toggleTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Invalid request body", Err: err})
		return
	}

	err := h.missionUC.ToggleTask(r.Context(), claims.UserID, id, req.Completed)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
