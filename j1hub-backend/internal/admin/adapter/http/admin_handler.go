package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/parada3456/wat_project-backend/internal/admin/adapter/http/dto"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	port "github.com/parada3456/wat_project-backend/internal/admin/port"
	"github.com/parada3456/wat_project-backend/internal/domain"
	missiondomain "github.com/parada3456/wat_project-backend/internal/mission/domain"
	"github.com/parada3456/wat_project-backend/internal/transport/http/middleware"
	userdto "github.com/parada3456/wat_project-backend/internal/user/adapter/http/dto"
	"github.com/parada3456/wat_project-backend/pkg/apperror"
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

func mapUserWithProfileToDTO(up port.UserWithProfile) *userdto.UserAccountDTO {
	var arrival *time.Time
	if !up.User.ArrivalDate.IsZero() {
		arrival = &up.User.ArrivalDate
	}
	var jobStart *time.Time
	if !up.User.JobStartDate.IsZero() {
		jobStart = &up.User.JobStartDate
	}

	var locUpdated *time.Time
	if !up.Profile.LocationUpdatedAt.IsZero() {
		locUpdated = &up.Profile.LocationUpdatedAt
	}

	var coords string
	if up.Profile.Lat != 0 || up.Profile.Lng != 0 {
		coords = fmt.Sprintf("%f,%f", up.Profile.Lat, up.Profile.Lng)
	}

	return &userdto.UserAccountDTO{
		ID:                  up.User.UserID,
		Email:               up.User.Email,
		Username:            up.Profile.Username,
		FirstName:           up.Profile.FirstName,
		LastName:            up.Profile.LastName,
		ProfileID:           up.Profile.ProfileID,
		PhoneNumber:         up.Profile.PhoneNumber,
		Bio:                 up.Profile.Bio,
		AvatarURL:           up.Profile.AvatarURL,
		RadarVisibility:     string(up.Profile.RadarVisibility),
		CurrentCoordinates:  coords,
		LocationUpdatedAt:   locUpdated,
		CurrentPhaseID:      up.User.CurrentPhaseID,
		TotalLifetimePoints: up.User.TotalLifetimePoints,
		CurrentPhasePoints:  up.User.CurrentPhasePoints,
		MissionStreak:       up.User.MissionStreak,
		ArrivalDate:         arrival,
		JobStartDate:        jobStart,
		CreatedAt:           up.User.CreatedAt,
		UpdatedAt:           up.User.UpdatedAt,
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
	pago := apperror.ParsePagination(r)
	ums, totalCount, err := h.adminUseCase.ListPendingVerifications(r.Context(), pago.Page, pago.PageSize)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	apperror.RespondList(w, ums, pago.Page, pago.PageSize, totalCount)
}

func (h *AdminHandler) VerifyMission(w http.ResponseWriter, r *http.Request) {
	log.Printf("debugprint: entering (*AdminHandler).VerifyMission")
	userMissionID := chi.URLParam(r, "id")
	claims := middleware.GetClaims(r.Context())
	if claims == nil || !claims.IsAdmin {
		apperror.RespondError(w, domain.ErrForbidden)
		return
	}

	var req dto.VerifyMissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Malformed request body: %w", domain.ErrInvalidInput))
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
	pago := apperror.ParsePagination(r)
	users, totalCount, err := h.adminUseCase.ListUsers(r.Context(), searchQuery, pago.Page, pago.PageSize)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	dtos := make([]*userdto.UserAccountDTO, len(users))
	for i, u := range users {
		dtos[i] = mapUserWithProfileToDTO(u)
	}

	apperror.RespondList(w, dtos, pago.Page, pago.PageSize, totalCount)
}

func (h *AdminHandler) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	up, err := h.adminUseCase.GetUserDetail(r.Context(), userID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	dto := mapUserWithProfileToDTO(*up)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto)
}

func (h *AdminHandler) AdjustPoints(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var req dto.AdjustPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Malformed request body: %w", domain.ErrInvalidInput))
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

func (h *AdminHandler) CreateMission(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*AdminHandler).CreateMission")
	var req dto.CreateMissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Malformed request body: %w", domain.ErrInvalidInput))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, fmt.Errorf("%s: %w", err.Error(), domain.ErrInvalidInput))
		return
	}

	taskCmds := make([]port.CreateTaskCmd, len(req.Tasks))
	for i, t := range req.Tasks {
		taskCmds[i] = port.CreateTaskCmd{Title: t.Title, Description: t.Description}
	}

	cmd := port.CreateMissionCmd{
		PhaseID:              req.PhaseID,
		Title:                req.Title,
		Description:          req.Description,
		Location:             req.Location,
		BasePoints:           req.BasePoints,
		IsMandatory:          req.IsMandatory,
		VerificationType:     missiondomain.VerificationType(req.VerificationType),
		DueDateType:          req.DueDateType,
		FixedDueDate:         req.FixedDueDate,
		RelativeTriggerEvent: req.RelativeTriggerEvent,
		RelativeDaysOffset:   req.RelativeDaysOffset,
		Tasks:                taskCmds,
	}

	result, err := h.adminUseCase.CreateMission(r.Context(), cmd)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
