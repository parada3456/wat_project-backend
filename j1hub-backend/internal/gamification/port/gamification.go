package port

import (
	"context"
	"time"
	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"
	frienddomain "github.com/j1hub/backend/internal/friend/domain"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
)

type CreditScoreRepository interface {
	Create(ctx context.Context, c *gamificationdomain.CreditScore) error
	FindByUserID(ctx context.Context, userID string) (*gamificationdomain.CreditScore, error)
	Decrement(ctx context.Context, userID string, delta int) error
}

type PointLedgerRepository interface {
	Insert(ctx context.Context, l *gamificationdomain.PointLedger) error
	InsertBatch(ctx context.Context, ledgers []gamificationdomain.PointLedger) error
	FindByUser(ctx context.Context, userID string) ([]gamificationdomain.PointLedger, error)
}

type BadgeRepository interface {
	FindByTriggerType(ctx context.Context, triggerType gamificationdomain.TriggerType) ([]gamificationdomain.Badge, error)
	FindEligible(ctx context.Context, userID string, triggerType gamificationdomain.TriggerType) ([]gamificationdomain.Badge, error)
}

type UserBadgeRepository interface {
	Insert(ctx context.Context, ub *gamificationdomain.UserBadge) error
	FindByUser(ctx context.Context, userID string) ([]gamificationdomain.UserBadge, error)
}

type RadarRepository interface {
	FindNearby(ctx context.Context, lat, lng, radius float64, staleMinutes int) ([]userdomain.Profile, error)
}

type LeaderboardRepository interface {
	FindByScope(ctx context.Context, scope, jobID string) ([]userdomain.User, error)
}

type UserRepository interface {
	IncrementPoints(ctx context.Context, userID string, lifetimeDelta, phaseDelta int) error
	SetPhase(ctx context.Context, userID, phaseID string) error
	FindByID(ctx context.Context, id string) (*userdomain.User, error)
	Update(ctx context.Context, u *userdomain.User) error
}

type ProfileRepository interface {
	FindByUserID(ctx context.Context, userID string) (*userdomain.Profile, error)
}

type FriendshipRepository interface {
	FindFriendsOf(ctx context.Context, userID string) ([]frienddomain.Friendship, error)
	FindByCanonicalPair(ctx context.Context, u1, u2 string) (*frienddomain.Friendship, error)
}

type UserMissionRepository interface {
	FindByUserAndPhase(ctx context.Context, userID, phaseID string) ([]missiondomain.UserMission, error)
	FindOverdue(ctx context.Context) ([]missiondomain.UserMission, error)
	BulkInsert(ctx context.Context, ums []missiondomain.UserMission) error
}

type JourneyPhaseRepository interface {
	FindByNumber(ctx context.Context, number int) (*missiondomain.JourneyPhase, error)
	FindByID(ctx context.Context, id string) (*missiondomain.JourneyPhase, error)
}

type UserPhaseHistoryRepository interface {
	Insert(ctx context.Context, h *missiondomain.UserPhaseHistory) error
	CompleteCurrentPhase(ctx context.Context, userID string, points int, completedAt time.Time) error
	FindByUserAndPhase(ctx context.Context, userID, phaseID string) (*missiondomain.UserPhaseHistory, error)
}

type MissionRepository interface {
	FindByPhase(ctx context.Context, phaseID string) ([]missiondomain.Mission, error)
}

type NotifierPort interface {
	Send(ctx context.Context, userID, title, body string) error
}
