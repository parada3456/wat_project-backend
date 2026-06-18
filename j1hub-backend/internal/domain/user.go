package domain

import (
	"time"
)

type RadarVisibility string

const (
	VisibilityShowAnonymous RadarVisibility = "Show_Anonymous"
	VisibilityShowFriends   RadarVisibility = "Show_Friends"
	VisibilityHidden        RadarVisibility = "Hidden"
)

func (v RadarVisibility) Valid() bool {
	switch v {
	case VisibilityShowAnonymous, VisibilityShowFriends, VisibilityHidden:
		return true
	}
	return false
}

type User struct {
	UserID              string
	Email               string
	PasswordHash        string
	FirstName           string
	LastName            string
	CurrentPhaseID      string
	TotalLifetimePoints int
	CurrentPhasePoints  int
	MissionStreak       int
	ArrivalDate         time.Time
	JobStartDate        time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Profile struct {
	ProfileID         string
	UserID            string
	PhoneNumber       string
	Bio               string
	AvatarURL         string
	RadarVisibility   RadarVisibility
	Lat               float64
	Lng               float64
	LocationUpdatedAt time.Time
	UpdatedAt         time.Time
}
