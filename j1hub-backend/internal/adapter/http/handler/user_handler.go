package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/j1hub/backend/pkg/apperror"
)

type UserUC interface {
	GetProfile(ctx context.Context, userID string) (*usecase.UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, cmd usecase.UpdateProfileCommand) error
}

type UserHandler struct {
	userUC   UserUC
	validate *validator.Validate
}

func NewUserHandler(userUC UserUC) *UserHandler {
	return &UserHandler{userUC: userUC, validate: validator.New()}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
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
