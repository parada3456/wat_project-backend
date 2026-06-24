package frienddomain

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
	FriendshipID string           `json:"friendship_id"`
	UserID1      string           `json:"user_id1"`
	UserID2      string           `json:"user_id2"`
	Status       FriendshipStatus `json:"status"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

func CanonicalOrder(a, b string) (string, string) {
	log.Println("debugprint: entering CanonicalOrder")
	if a < b {
		return a, b
	}
	return b, a
}
