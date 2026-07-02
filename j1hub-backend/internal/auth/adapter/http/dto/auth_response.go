package dto

import (
	"fmt"
	"time"

	userdto "github.com/j1hub/backend/internal/user/adapter/http/dto"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	port "github.com/j1hub/backend/internal/auth/port"
)

type AuthData struct {
	*port.TokenPair
	TokenType string `json:"token_type"`
}

type LoginResponse struct {
	User *userdto.UserAccountDTO `json:"user"`
	Auth AuthData                `json:"auth"`
}

func NewLoginResponse(user *userdomain.User, profile *userdomain.Profile, tokens *port.TokenPair) *LoginResponse {
	var arrival *time.Time
	if user != nil && !user.ArrivalDate.IsZero() {
		arrival = &user.ArrivalDate
	}
	var jobStart *time.Time
	if user != nil && !user.JobStartDate.IsZero() {
		jobStart = &user.JobStartDate
	}

	var userDTO *userdto.UserAccountDTO
	if user != nil {
		userDTO = &userdto.UserAccountDTO{
			ID:                  user.UserID,
			Email:               user.Email,
			CurrentPhaseID:      user.CurrentPhaseID,
			TotalLifetimePoints: user.TotalLifetimePoints,
			CurrentPhasePoints:  user.CurrentPhasePoints,
			MissionStreak:       user.MissionStreak,
			ArrivalDate:         arrival,
			JobStartDate:        jobStart,
			CreatedAt:           user.CreatedAt,
			UpdatedAt:           user.UpdatedAt,
		}

		if profile != nil {
			var locUpdated *time.Time
			if !profile.LocationUpdatedAt.IsZero() {
				locUpdated = &profile.LocationUpdatedAt
			}
			var coords string
			if profile.Lat != 0 || profile.Lng != 0 {
				coords = fmt.Sprintf("%f,%f", profile.Lat, profile.Lng)
			}
			userDTO.ProfileID = profile.ProfileID
			userDTO.Username = profile.Username
			userDTO.FirstName = profile.FirstName
			userDTO.LastName = profile.LastName
			userDTO.PhoneNumber = profile.PhoneNumber
			userDTO.Bio = profile.Bio
			userDTO.AvatarURL = profile.AvatarURL
			userDTO.RadarVisibility = string(profile.RadarVisibility)
			userDTO.CurrentCoordinates = coords
			userDTO.LocationUpdatedAt = locUpdated
		}
	}

	return &LoginResponse{
		User: userDTO,
		Auth: AuthData{
			TokenPair: tokens,
			TokenType: "Bearer",
		},
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
