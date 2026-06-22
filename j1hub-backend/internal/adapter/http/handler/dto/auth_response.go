package dto

import (
	"time"

	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/j1hub/backend/internal/port"
)

type RegisterResponse struct {
	UserID         string    `json:"user_id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	CurrentPhaseID string    `json:"current_phase_id"`
	CreatedAt      time.Time `json:"created_at"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	ExpiresAt      time.Time `json:"expires_at"`
	TokenType      string    `json:"token_type"`
}

func NewRegisterResponse(user *userdomain.User, tokens *port.TokenPair) *RegisterResponse {
	return &RegisterResponse{
		UserID:         user.UserID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		CurrentPhaseID: user.CurrentPhaseID,
		CreatedAt:      user.CreatedAt,
		AccessToken:    tokens.AccessToken,
		RefreshToken:   tokens.RefreshToken,
		ExpiresAt:      tokens.ExpiresAt,
		TokenType:      "Bearer",
	}
}

type LoginResponse struct {
	UserID       string    `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

func NewLoginResponse(user *userdomain.User, tokens *port.TokenPair) *LoginResponse {
	return &LoginResponse{
		UserID:       user.UserID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
		TokenType:    "Bearer",
	}
}

type RefreshResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	TokenType   string    `json:"token_type"`
}

func NewRefreshResponse(tokens *port.TokenPair) *RefreshResponse {
	return &RefreshResponse{
		AccessToken: tokens.AccessToken,
		ExpiresAt:   tokens.ExpiresAt,
		TokenType:   "Bearer",
	}
}
