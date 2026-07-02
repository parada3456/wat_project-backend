package dto

import (
	"fmt"
	"time"

	userdomain "github.com/parada3456/wat_project-backend/internal/user/domain"
	userusecase "github.com/parada3456/wat_project-backend/internal/user/usecase"
)

type UserAccountDTO struct {
	ID                  string     `json:"user_id"`
	Email               string     `json:"email"`
	Username            string     `json:"username"`
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	ProfileID           string     `json:"profile_id,omitempty"`
	PhoneNumber         string     `json:"phone_number,omitempty"`
	Bio                 string     `json:"bio,omitempty"`
	AvatarURL           string     `json:"avatar_url,omitempty"`
	RadarVisibility     string     `json:"radar_visibility,omitempty"`
	CurrentCoordinates  string     `json:"current_coordinates,omitempty"`
	LocationUpdatedAt   *time.Time `json:"location_updated_at,omitempty"`
	CurrentPhaseID      string     `json:"current_phase_id"`
	TotalLifetimePoints int        `json:"total_lifetime_points"`
	CurrentPhasePoints  int        `json:"current_phase_points"`
	MissionStreak       int        `json:"mission_streak"`
	ArrivalDate         *time.Time `json:"arrival_date,omitempty"`
	JobStartDate        *time.Time `json:"job_start_date,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type ProfileCreditScoreDTO struct {
	CreditID     string    `json:"credit_id"`
	UserID       string    `json:"user_id"`
	CurrentScore int       `json:"current_score"`
	LastUpdated  time.Time `json:"last_updated"`
}

type UserJobDTO struct {
	UserID     string     `json:"user_id"`
	JobID      string     `json:"job_id"`
	AssignedAt time.Time  `json:"assigned_at"`
	IsMain     bool       `json:"is_main"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
}

type GetProfileResponse struct {
	User        *UserAccountDTO        `json:"user"`
	CreditScore *ProfileCreditScoreDTO `json:"credit_score"`
	UserJobs    []UserJobDTO           `json:"user_jobs"`
}

// type GetPublicProfileResponse struct {
// 	UserID    string   `json:"user_id"`
// 	FirstName string   `json:"first_name"`
// 	LastName  string   `json:"last_name"`
// 	AvatarURL string   `json:"avatar_url,omitempty"`
// 	UserJobs  []string `json:"user_jobs,omitempty"`
// }

// func NewGetPublicProfileResponse(user *userdomain.User, profile *userdomain.Profile) *GetPublicProfileResponse {
// 	if user == nil {
// 		return nil
// 	}

//		resp := &GetPublicProfileResponse{
//			UserID:    user.UserID,
//			FirstName: user.FirstName,
//			LastName:  user.LastName,
//		}
//		if profile != nil {
//			resp.AvatarURL = profile.AvatarURL
//		}
//		return resp
//	}
type ProfileResponse struct {
	UserID             string     `json:"user_id"`
	ProfileID          string     `json:"profile_id,omitempty"`
	FirstName          string     `json:"first_name"`
	LastName           string     `json:"last_name"`
	PhoneNumber        string     `json:"phone_number,omitempty"`
	Bio                string     `json:"bio,omitempty"`
	AvatarURL          string     `json:"avatar_url,omitempty"`
	RadarVisibility    string     `json:"radar_visibility,omitempty"`
	CurrentCoordinates string     `json:"current_coordinates,omitempty"`
	LocationUpdatedAt  *time.Time `json:"location_updated_at,omitempty"`
	UpdatedAt          time.Time  `json:"updated_at,omitempty"`
}

type GetPublicProfileResponse struct {
	UserID  string           `json:"user_id"`
	Profile *ProfileResponse `json:"profile,omitempty"` // FIX: Changed type to *ProfileResponse
}

func NewGetPublicProfileResponse(user *userdomain.User, profile *userdomain.Profile) *ProfileResponse {
	if user == nil {
		return nil
	}

	resp := &ProfileResponse{}

	if profile != nil {
		var locUpdated *time.Time
		// FIX: Read from the incoming parameter 'profile', not 'resp.Profile'
		if !profile.LocationUpdatedAt.IsZero() {
			locUpdated = &profile.LocationUpdatedAt
		}

		var coords string
		// FIX: Read from the incoming parameter 'profile'
		if profile.Lat != 0 || profile.Lng != 0 {
			coords = fmt.Sprintf("%f,%f", profile.Lat, profile.Lng)
		}

		// FIX: Assigning a pointer to ProfileResponse to match the updated struct type
		resp = &ProfileResponse{
			UserID:             user.UserID,
			ProfileID:          profile.ProfileID,
			FirstName:          profile.FirstName,
			LastName:           profile.LastName,
			AvatarURL:          profile.AvatarURL,
			PhoneNumber:        profile.PhoneNumber,
			Bio:                profile.Bio,
			CurrentCoordinates: coords,
			LocationUpdatedAt:  locUpdated,
			RadarVisibility:    string(profile.RadarVisibility),
			UpdatedAt:          profile.UpdatedAt,
		}
	}

	return resp
}
func NewGetProfileResponse(resp *userusecase.UserProfileResponse) GetProfileResponse {
	dto := GetProfileResponse{
		UserJobs: []UserJobDTO{},
	}

	if resp == nil {
		return dto
	}

	if resp.User != nil {
		var arrival *time.Time
		if !resp.User.ArrivalDate.IsZero() {
			arrival = &resp.User.ArrivalDate
		}
		var jobStart *time.Time
		if !resp.User.JobStartDate.IsZero() {
			jobStart = &resp.User.JobStartDate
		}
		dto.User = &UserAccountDTO{
			ID:                  resp.User.UserID,
			Email:               resp.User.Email,
			CurrentPhaseID:      resp.User.CurrentPhaseID,
			TotalLifetimePoints: resp.User.TotalLifetimePoints,
			CurrentPhasePoints:  resp.User.CurrentPhasePoints,
			MissionStreak:       resp.User.MissionStreak,
			ArrivalDate:         arrival,
			JobStartDate:        jobStart,
			CreatedAt:           resp.User.CreatedAt,
			UpdatedAt:           resp.User.UpdatedAt,
		}

		if resp.Profile != nil {
			var locUpdated *time.Time
			if !resp.Profile.LocationUpdatedAt.IsZero() {
				locUpdated = &resp.Profile.LocationUpdatedAt
			}
			var coords string
			if resp.Profile.Lat != 0 || resp.Profile.Lng != 0 {
				coords = fmt.Sprintf("%f,%f", resp.Profile.Lat, resp.Profile.Lng)
			}
			dto.User.ProfileID = resp.Profile.ProfileID
			dto.User.Username = resp.Profile.Username
			dto.User.FirstName = resp.Profile.FirstName
			dto.User.LastName = resp.Profile.LastName
			dto.User.PhoneNumber = resp.Profile.PhoneNumber
			dto.User.Bio = resp.Profile.Bio
			dto.User.AvatarURL = resp.Profile.AvatarURL
			dto.User.RadarVisibility = string(resp.Profile.RadarVisibility)
			dto.User.CurrentCoordinates = coords
			dto.User.LocationUpdatedAt = locUpdated
		}
	}

	if resp.CreditScore != nil {
		dto.CreditScore = &ProfileCreditScoreDTO{
			CreditID:     resp.CreditScore.CreditID,
			UserID:       resp.CreditScore.UserID,
			CurrentScore: resp.CreditScore.CurrentScore,
			LastUpdated:  resp.CreditScore.LastUpdated,
		}
	}

	if resp.UserJobs != nil {
		for _, uj := range resp.UserJobs {
			dto.UserJobs = append(dto.UserJobs, UserJobDTO{
				UserID:     uj.UserID,
				JobID:      uj.JobID,
				AssignedAt: uj.AssignedAt,
				IsMain:     uj.IsMain,
				StartDate:  uj.StartDate,
				EndDate:    uj.EndDate,
			})
		}
	}

	return dto
}
