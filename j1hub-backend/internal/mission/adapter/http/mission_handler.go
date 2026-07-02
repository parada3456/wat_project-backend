package http

import (
	"github.com/parada3456/wat_project-backend/internal/mission/adapter/http/dto"
	missiondomain "github.com/parada3456/wat_project-backend/internal/mission/domain"

	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/parada3456/wat_project-backend/internal/domain"
	missionusecase "github.com/parada3456/wat_project-backend/internal/mission/usecase"
	"github.com/parada3456/wat_project-backend/internal/transport/http/middleware"
	"github.com/parada3456/wat_project-backend/pkg/apperror"
)

type MissionUC interface {
	ListAvailableMissions(ctx context.Context, userID string, ids []string) ([]missiondomain.UserMission, error)
	ListStaticMissions(ctx context.Context, userID string, ids []string) ([]missiondomain.Mission, error)
	GetMissionDetail(ctx context.Context, userID, userMissionID string) (*missionusecase.MissionDetailResponse, error)
	ToggleTask(ctx context.Context, userID, userTaskID string, completed bool) error
	ListTasks(ctx context.Context, ids []string) ([]missiondomain.Task, error)
	ListUserTasks(ctx context.Context, ids []string) ([]missiondomain.UserTask, error)
}

func (h *MissionHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*MissionHandler).ListTasks")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}
	ids := r.URL.Query()["ids"]
	tasks, err := h.missionUC.ListTasks(r.Context(), ids)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	pago := apperror.ParsePagination(r)
	apperror.RespondList(w, tasks, pago.Page, pago.PageSize, len(tasks))
}

func (h *MissionHandler) ListUserTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*MissionHandler).ListUserTasks")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}
	ids := r.URL.Query()["ids"]
	userTasks, err := h.missionUC.ListUserTasks(r.Context(), ids)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	pago := apperror.ParsePagination(r)
	apperror.RespondList(w, userTasks, pago.Page, pago.PageSize, len(userTasks))
}

type CompleteMissionUC interface {
	SubmitProof(ctx context.Context, userID, userMissionID string, file io.Reader, contentType string) error
}

type MissionHandler struct {
	missionUC  MissionUC
	completeUC CompleteMissionUC
}

func NewMissionHandler(missionUC MissionUC, completeUC CompleteMissionUC) *MissionHandler {
	log.Println("debugprint: entering NewMissionHandler")
	return &MissionHandler{
		missionUC:  missionUC,
		completeUC: completeUC,
	}
}

func (h *MissionHandler) ListMissions(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*MissionHandler).ListMissions")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	ids := r.URL.Query()["ids"]
	missions, err := h.missionUC.ListStaticMissions(r.Context(), claims.UserID, ids)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	pago := apperror.ParsePagination(r)
	apperror.RespondList(w, missions, pago.Page, pago.PageSize, len(missions))
}

func (h *MissionHandler) ListUserMissions(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*MissionHandler).ListUserMissions")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	ids := r.URL.Query()["ids"]
	missions, err := h.missionUC.ListAvailableMissions(r.Context(), claims.UserID, ids)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	pago := apperror.ParsePagination(r)
	apperror.RespondList(w, missions, pago.Page, pago.PageSize, len(missions))
}

func (h *MissionHandler) GetMissionDetail(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*MissionHandler).GetMissionDetail")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	detail, err := h.missionUC.GetMissionDetail(r.Context(), claims.UserID, id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respDTO := dto.NewMissionDetailResponse(detail)
	json.NewEncoder(w).Encode(respDTO)
}

func (h *MissionHandler) SubmitProof(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*MissionHandler).SubmitProof")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	r.ParseMultipartForm(10 << 20) // 10MB max
	file, header, err := r.FormFile("proof")
	if err != nil {
		apperror.RespondError(w, fmt.Errorf("Proof file required: %w", domain.ErrInvalidInput))
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

func (h *MissionHandler) ToggleTask(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*MissionHandler).ToggleTask")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	taskID := chi.URLParam(r, "taskId")
	if taskID == "" {
		taskID = id
	}

	var req dto.ToggleTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	completed := false
	if req.Completed != nil {
		completed = *req.Completed
	} else if req.IsCompleted != nil {
		completed = *req.IsCompleted
	}

	err := h.missionUC.ToggleTask(r.Context(), claims.UserID, taskID, completed)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
