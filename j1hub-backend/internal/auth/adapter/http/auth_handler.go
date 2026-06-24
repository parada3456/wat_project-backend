package http

import (
	"github.com/j1hub/backend/internal/auth/adapter/http/dto"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/j1hub/backend/internal/domain"
	authusecase "github.com/j1hub/backend/internal/auth/usecase"
	port "github.com/j1hub/backend/internal/auth/port"
	"github.com/j1hub/backend/pkg/apperror"
)

type RegisterUserUC interface {
	Register(ctx context.Context, cmd authusecase.RegisterCommand) (*userdomain.User, *port.TokenPair, error)
}

type LoginUC interface {
	Login(ctx context.Context, cmd authusecase.LoginCommand) (*userdomain.User, *port.TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*port.TokenPair, error)
}

type AuthHandler struct {
	registerUC RegisterUserUC
	loginUC    LoginUC
	validate   *validator.Validate
}

func NewAuthHandler(registerUC RegisterUserUC, loginUC LoginUC) *AuthHandler {
	log.Println("debugprint: entering NewAuthHandler")
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
		validate:   validator.New(),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*AuthHandler).Register")
	var req dto.RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Validation failed: %w", domain.ErrInvalidInput))
		return
	}

	user, tokens, err := h.registerUC.Register(r.Context(), authusecase.RegisterCommand{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	respDTO := dto.NewRegisterResponse(user, tokens)
	json.NewEncoder(w).Encode(respDTO)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*AuthHandler).Login")

	if r.Body == nil || r.ContentLength == 0 {
		apperror.RespondError(w, fmt.Errorf("Request body cannot be empty: %w", domain.ErrInvalidInput))
		return
	}

	var req dto.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Validation failed: %w", domain.ErrInvalidInput))
		return
	}

	user, tokens, err := h.loginUC.Login(r.Context(), authusecase.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("login failed for email %s: %v", req.Email, err)
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respDTO := dto.NewLoginResponse(user, tokens)
	if err := json.NewEncoder(w).Encode(respDTO); err != nil {
		log.Printf("failed to encode login response: %v", err)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*AuthHandler).Logout")
	var req dto.LogoutReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*AuthHandler).Refresh")
	var req dto.RefreshReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	tokens, err := h.loginUC.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid refresh token: %w", domain.ErrUnauthorized))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respDTO := dto.NewRefreshResponse(tokens)
	json.NewEncoder(w).Encode(respDTO)
}
