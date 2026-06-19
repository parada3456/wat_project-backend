package domain

import (
	"log"
	"time"
)

type FriendshipStatus string

const (
	FriendshipPending  FriendshipStatus = "Pending"
	FriendshipAccepted FriendshipStatus = "Accepted"
	FriendshipBlocked  FriendshipStatus = "Blocked"
)

func (s FriendshipStatus) Valid() bool {
	log.Println("debugprint: entering (FriendshipStatus).Valid")
	switch s {
	case FriendshipPending, FriendshipAccepted, FriendshipBlocked:
		return true
	}
	return false
}

type Friendship struct {
	FriendshipID string
	UserID1      string
	UserID2      string
	Status       FriendshipStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func CanonicalOrder(a, b string) (string, string) {
	log.Println("debugprint: entering CanonicalOrder")
	if a < b {
		return a, b
	}
	return b, a
}
