package userdomain

import (
	"log"
	"time"
)

type RadarVisibility string

const (
	VisibilityShowAnonymous RadarVisibility = "show_anonymous"
	VisibilityShowFriends   RadarVisibility = "show_friends"
	VisibilityHidden        RadarVisibility = "hidden"
)

func (v RadarVisibility) Valid() bool {
	log.Println("debugprint: entering (RadarVisibility).Valid")
	switch v {
	case VisibilityShowAnonymous, VisibilityShowFriends, VisibilityHidden:
		return true
	}
	return false
}

type User struct {
	UserID              string    `json:"user_id"`
	Email               string    `json:"email"`
	PasswordHash        string    `json:"password_hash,omitempty"`
	CurrentPhaseID      string    `json:"current_phase_id"`
	TotalLifetimePoints int       `json:"total_lifetime_points"`
	CurrentPhasePoints  int       `json:"current_phase_points"`
	MissionStreak       int       `json:"mission_streak"`
	ArrivalDate         time.Time `json:"arrival_date"`
	JobStartDate        time.Time `json:"job_start_date"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type Profile struct {
	ProfileID         string          `json:"profile_id"`
	UserID            string          `json:"user_id"`
	FirstName         string          `json:"first_name"`
	LastName          string          `json:"last_name"`
	PhoneNumber       string          `json:"phone_number"`
	Bio               string          `json:"bio"`
	AvatarURL         string          `json:"avatar_url"`
	RadarVisibility   RadarVisibility `json:"radar_visibility"`
	Lat               float64         `json:"lat"`
	Lng               float64         `json:"lng"`
	LocationUpdatedAt time.Time       `json:"location_updated_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type UserJob struct {
	UserID     string     `json:"user_id"`
	JobID      string     `json:"job_id"`
	AssignedAt time.Time  `json:"assigned_at"`
	IsMain     bool       `json:"is_main"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
}
