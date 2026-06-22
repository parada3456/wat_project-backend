package dto

import (
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
)

type GetProfileResponse struct {
	User        *domain.User    `json:"user"`
	Profile     *domain.Profile `json:"profile"`
	UserID      string          `json:"user_id,omitempty"`
	Email       string          `json:"email,omitempty"`
	FirstName   string          `json:"first_name,omitempty"`
	LastName    string          `json:"last_name,omitempty"`
	Points      int             `json:"points"`
	Bio         string          `json:"bio,omitempty"`
	AvatarURL   string          `json:"avatar_url,omitempty"`
	CreditScore int             `json:"credit_score"`
}

type GetPublicProfileResponse struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

func NewGetPublicProfileResponse(user *domain.User, profile *domain.Profile) *GetPublicProfileResponse {
	resp := &GetPublicProfileResponse{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	if profile != nil {
		resp.AvatarURL = profile.AvatarURL
	}
	return resp
}

func NewGetProfileResponse(resp *usecase.UserProfileResponse) GetProfileResponse {
	dto := GetProfileResponse{}

	if resp == nil {
		return dto
	}

	if resp.User != nil {
		dto.UserID = resp.User.UserID
		dto.Email = resp.User.Email
		dto.FirstName = resp.User.FirstName
		dto.LastName = resp.User.LastName
		dto.Points = resp.User.TotalLifetimePoints
	}

	if resp.Profile != nil {
		dto.Bio = resp.Profile.Bio
		dto.AvatarURL = resp.Profile.AvatarURL
	}

	if resp.CreditScore != nil {
		dto.CreditScore = resp.CreditScore.CurrentScore
	} else {
		dto.CreditScore = 0
	}

	return dto
}
