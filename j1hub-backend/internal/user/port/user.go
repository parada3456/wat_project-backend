package port

import (
	"context"
	"time"
	userdomain "github.com/j1hub/backend/internal/user/domain"
	frienddomain "github.com/j1hub/backend/internal/friend/domain"
	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *userdomain.User) error
	FindByID(ctx context.Context, id string) (*userdomain.User, error)
	FindByEmail(ctx context.Context, email string) (*userdomain.User, error)
	Update(ctx context.Context, u *userdomain.User) error
	IncrementPoints(ctx context.Context, userID string, lifetimeDelta, phaseDelta int) error
	ResetStreak(ctx context.Context, userID string) error
	SetPhase(ctx context.Context, userID, phaseID string) error
	Delete(ctx context.Context, id string) error
	FindUserJob(ctx context.Context, userID string) (*userdomain.UserJob, error)
	FindUserJobs(ctx context.Context, userID string) ([]userdomain.UserJob, error)
	AssignJob(ctx context.Context, userID, jobID string, isMain bool, startDate, endDate *time.Time) error
}

type ProfileRepository interface {
	Create(ctx context.Context, p *userdomain.Profile) error
	FindByUserID(ctx context.Context, userID string) (*userdomain.Profile, error)
	Update(ctx context.Context, p *userdomain.Profile) error
	UpdateLocation(ctx context.Context, userID string, lat, lng float64) error
	UpdateVisibility(ctx context.Context, userID string, visibility userdomain.RadarVisibility) error
}

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(plain, hash string) bool
}

type FriendshipRepository interface {
	FindByCanonicalPair(ctx context.Context, u1, u2 string) (*frienddomain.Friendship, error)
}

type CreditScoreRepository interface {
	FindByUserID(ctx context.Context, userID string) (*gamificationdomain.CreditScore, error)
}
