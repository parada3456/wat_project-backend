package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/transport/http/middleware"
	userhandlerdto "github.com/j1hub/backend/internal/user/adapter/http/dto"
	userdomain "github.com/j1hub/backend/internal/user/domain"
	userusecase "github.com/j1hub/backend/internal/user/usecase"
	"github.com/j1hub/backend/pkg/apperror"
)

type UserUC interface {
	GetProfile(ctx context.Context, userID string) (*userusecase.UserProfileResponse, error)
	GetPublicProfile(ctx context.Context, currentUserID, targetUserID string) (*userdomain.User, *userdomain.Profile, error)
	UpdateProfile(ctx context.Context, userID string, cmd userusecase.UpdateProfileCommand) error
	UpdateLocation(ctx context.Context, userID string, lat, lng float64) error
	UpdateSettings(ctx context.Context, userID string, settings map[string]interface{}) error
	DeleteAccount(ctx context.Context, userID string, password string) error
	AssignJob(ctx context.Context, userID, jobID string, isMain bool, startDate, endDate *time.Time) error
	UpdatePassword(ctx context.Context, userID string, currentPassword, newPassword string) error
}

type UserHandler struct {
	userUC   UserUC
	validate *validator.Validate
}

func NewUserHandler(userUC UserUC) *UserHandler {
	log.Println("debugprint: entering NewUserHandler")
	return &UserHandler{userUC: userUC, validate: validator.New()}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).GetProfile")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		log.Println("debugprint: claim is nil")
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	resp, err := h.userUC.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		log.Println("debugprint: err get profile usecase")
		apperror.RespondError(w, err)
		return
	}

	respDTO := userhandlerdto.NewGetProfileResponse(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respDTO)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).UpdateProfile")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	var req userhandlerdto.UpdateProfileReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Validation failed: %w", domain.ErrInvalidInput))
		return
	}

	err := h.userUC.UpdateProfile(r.Context(), claims.UserID, userusecase.UpdateProfileCommand{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
	})
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).UpdateLocation")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	var req userhandlerdto.UpdateLocationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Validation failed: %w", domain.ErrInvalidInput))
		return
	}

	err := h.userUC.UpdateLocation(r.Context(), claims.UserID, req.Latitude, req.Longitude)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).UpdateSettings")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	var settings map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	err := h.userUC.UpdateSettings(r.Context(), claims.UserID, settings)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).DeleteAccount")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	var req userhandlerdto.DeleteAccountReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Validation failed: %w", domain.ErrInvalidInput))
		return
	}

	err := h.userUC.DeleteAccount(r.Context(), claims.UserID, req.CurrentPassword)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).GetPublicProfile")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	targetID := chi.URLParam(r, "id")
	user, profile, err := h.userUC.GetPublicProfile(r.Context(), claims.UserID, targetID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	respDTO := userhandlerdto.GetPublicProfileResponse{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	if profile != nil {
		respDTO.AvatarURL = profile.AvatarURL
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respDTO)
}

func (h *UserHandler) AssignJob(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).AssignJob")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	var req struct {
		JobID     string     `json:"job_id" validate:"required"`
		IsMain    bool       `json:"is_main"`
		StartDate *time.Time `json:"start_date"`
		EndDate   *time.Time `json:"end_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	if err := h.validate.Var(req.JobID, "required"); err != nil {
		apperror.RespondError(w, fmt.Errorf("Validation failed: %w", domain.ErrInvalidInput))
		return
	}

	err := h.userUC.AssignJob(r.Context(), claims.UserID, req.JobID, req.IsMain, req.StartDate, req.EndDate)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).UpdatePassword")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	var req userhandlerdto.UpdatePasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Validation failed: %w", domain.ErrInvalidInput))
		return
	}

	err := h.userUC.UpdatePassword(r.Context(), claims.UserID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

