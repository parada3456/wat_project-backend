package dto

import (
	userdomain "github.com/j1hub/backend/internal/user/domain"
	userusecase "github.com/j1hub/backend/internal/user/usecase"
)

type GetProfileResponse struct {
	User        *userdomain.User     `json:"user"`
	Profile     *userdomain.Profile  `json:"profile"`
	UserJobs    []userdomain.UserJob `json:"user_jobs,omitempty"`
	UserID      string               `json:"user_id,omitempty"`
	Email       string               `json:"email,omitempty"`
	FirstName   string               `json:"first_name,omitempty"`
	LastName    string               `json:"last_name,omitempty"`
	Points      int                  `json:"points"`
	Bio         string               `json:"bio,omitempty"`
	AvatarURL   string               `json:"avatar_url,omitempty"`
	CreditScore int                  `json:"credit_score"`
}

type GetPublicProfileResponse struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

func NewGetPublicProfileResponse(user *userdomain.User, profile *userdomain.Profile) *GetPublicProfileResponse {
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

func NewGetProfileResponse(resp *userusecase.UserProfileResponse) GetProfileResponse {
	dto := GetProfileResponse{}

	if resp == nil {
		return dto
	}

	if resp.User != nil {
		dto.User = resp.User
		dto.UserID = resp.User.UserID
		dto.Email = resp.User.Email
		dto.FirstName = resp.User.FirstName
		dto.LastName = resp.User.LastName
		dto.Points = resp.User.TotalLifetimePoints
	}

	if resp.Profile != nil {
		dto.Profile = resp.Profile
		dto.Bio = resp.Profile.Bio
		dto.AvatarURL = resp.Profile.AvatarURL
	}

	if resp.UserJobs != nil {
		dto.UserJobs = resp.UserJobs
	}

	if resp.CreditScore != nil {
		dto.CreditScore = resp.CreditScore.CurrentScore
	} else {
		dto.CreditScore = 0
	}

	return dto
}
