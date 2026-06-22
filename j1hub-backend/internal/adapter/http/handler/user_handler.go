package handler

import (
	"github.com/j1hub/backend/internal/adapter/http/handler/dto"

	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/j1hub/backend/pkg/apperror"
)

type UserUC interface {
	GetProfile(ctx context.Context, userID string) (*usecase.UserProfileResponse, error)
	GetPublicProfile(ctx context.Context, currentUserID, targetUserID string) (*domain.User, *domain.Profile, error)
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

	// respDTO := dto.GetProfileResponse{
	// 	User:    resp.User,
	// 	Profile: resp.Profile,
	// }
	// if resp.User != nil {
	// 	respDTO.UserID = resp.User.UserID
	// 	respDTO.Email = resp.User.Email
	// 	respDTO.FirstName = resp.User.FirstName
	// 	respDTO.LastName = resp.User.LastName
	// 	respDTO.Points = resp.User.TotalLifetimePoints
	// }
	// if resp.Profile != nil {
	// 	respDTO.Bio = resp.Profile.Bio
	// 	respDTO.AvatarURL = resp.Profile.AvatarURL
	// }
	// if resp.CreditScore != nil {
	// 	respDTO.CreditScore = resp.CreditScore.CurrentScore
	// }
	respDTO := dto.NewGetProfileResponse(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respDTO)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).UpdateProfile")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var req dto.UpdateProfileReq
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

func (h *UserHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).UpdateLocation")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var req dto.UpdateLocationReq
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

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).DeleteAccount")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var req dto.DeleteAccountReq
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

func (h *UserHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*UserHandler).GetPublicProfile")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	targetID := chi.URLParam(r, "id")
	user, profile, err := h.userUC.GetPublicProfile(r.Context(), claims.UserID, targetID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	respDTO := dto.GetPublicProfileResponse{
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
