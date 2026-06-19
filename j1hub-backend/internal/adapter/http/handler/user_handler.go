package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/j1hub/backend/pkg/apperror"
)

type UserUC interface {
	GetProfile(ctx context.Context, userID string) (*usecase.UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, cmd usecase.UpdateProfileCommand) error
	UpdateLocation(ctx context.Context, userID string, lat, lng float64) error
	UpdateSettings(ctx context.Context, userID string, settings map[string]interface{}) error
	DeleteAccount(ctx context.Context, userID string, password string) error
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
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	resp, err := h.userUC.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

type updateProfileReq struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).UpdateProfile")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var req updateProfileReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Invalid request body", Err: err})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Validation failed", Err: err})
		return
	}

	err := h.userUC.UpdateProfile(r.Context(), claims.UserID, usecase.UpdateProfileCommand{
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

type updateLocationReq struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

func (h *UserHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).UpdateLocation")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var req updateLocationReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Invalid request body", Err: err})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Validation failed", Err: err})
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
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var settings map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Invalid request body", Err: err})
		return
	}

	err := h.userUC.UpdateSettings(r.Context(), claims.UserID, settings)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type deleteAccountReq struct {
	CurrentPassword string `json:"current_password" validate:"required"`
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).DeleteAccount")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var req deleteAccountReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Invalid request body", Err: err})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Validation failed", Err: err})
		return
	}

	err := h.userUC.DeleteAccount(r.Context(), claims.UserID, req.CurrentPassword)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
