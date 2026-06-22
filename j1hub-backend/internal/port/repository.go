package port

import (
	"context"
	"time"

	frienddomain "github.com/j1hub/backend/internal/friend/domain"
	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	notificationdomain "github.com/j1hub/backend/internal/notification/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/j1hub/backend/internal/domain"
	expensedomain "github.com/j1hub/backend/internal/expense/domain"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
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
}

type ProfileRepository interface {
	Create(ctx context.Context, p *userdomain.Profile) error
	FindByUserID(ctx context.Context, userID string) (*userdomain.Profile, error)
	Update(ctx context.Context, p *userdomain.Profile) error
	UpdateLocation(ctx context.Context, userID string, lat, lng float64) error
	UpdateVisibility(ctx context.Context, userID string, visibility userdomain.RadarVisibility) error
}

type CreditScoreRepository interface {
	Create(ctx context.Context, c *gamificationdomain.CreditScore) error
	FindByUserID(ctx context.Context, userID string) (*gamificationdomain.CreditScore, error)
	Decrement(ctx context.Context, userID string, delta int) error
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
	FindByID(ctx context.Context, id string) (*missiondomain.Mission, error)
}

type UserMissionRepository interface {
	BulkInsert(ctx context.Context, ums []missiondomain.UserMission) error
	FindByUserAndPhase(ctx context.Context, userID, phaseID string) ([]missiondomain.UserMission, error)
	FindByID(ctx context.Context, id string) (*missiondomain.UserMission, error)
	UpdateStatus(ctx context.Context, id string, status missiondomain.UserMissionStatus) error
	UpdateVerification(ctx context.Context, id string, verifiedAt time.Time, verifiedBy string) error
	UpdateReward(ctx context.Context, id string, reward *gamificationdomain.PointReward, rewardedAt time.Time) error
	FindOverdue(ctx context.Context) ([]missiondomain.UserMission, error)
}

type TaskRepository interface {
	FindByMission(ctx context.Context, missionID string) ([]missiondomain.Task, error)
}

type UserTaskRepository interface {
	Upsert(ctx context.Context, ut *missiondomain.UserTask) error
	FindByUserMission(ctx context.Context, userMissionID string) ([]missiondomain.UserTask, error)
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

type FriendshipRepository interface {
	Insert(ctx context.Context, f *frienddomain.Friendship) error
	FindByCanonicalPair(ctx context.Context, u1, u2 string) (*frienddomain.Friendship, error)
	FindByID(ctx context.Context, id string) (*frienddomain.Friendship, error)
	UpdateStatus(ctx context.Context, id string, status frienddomain.FriendshipStatus) error
	FindFriendsOf(ctx context.Context, userID string) ([]frienddomain.Friendship, error)
}

type ExpenseTransactionRepository interface {
	Insert(ctx context.Context, t *expensedomain.ExpenseTransaction) error
	FindByID(ctx context.Context, id string) (*expensedomain.ExpenseTransaction, error)
	MarkSettled(ctx context.Context, id string) error
}

type ExpenseSplitRepository interface {
	BulkInsert(ctx context.Context, splits []expensedomain.ExpenseSplit) error
	FindByID(ctx context.Context, id string) (*expensedomain.ExpenseSplit, error)
	UpdatePaymentStatus(ctx context.Context, id string, status expensedomain.PaymentStatus, slipURL string) error
	UpdateApproval(ctx context.Context, id string, status expensedomain.ApprovalStatus, settledAt *time.Time) error
	FindOverdue(ctx context.Context) ([]expensedomain.ExpenseSplit, error)
	CountUnsettled(ctx context.Context, transactionID string) (int, error)
}

type JobPostingRepository interface {
	FindWithFilters(ctx context.Context, filters map[string]interface{}) ([]jobdomain.JobPosting, error)
	FindByID(ctx context.Context, id string) (*jobdomain.JobPosting, error)
	Upsert(ctx context.Context, job *jobdomain.JobPosting) error
}

type JobHousingRepository interface {
	FindByJobID(ctx context.Context, jobID string) ([]jobdomain.JobHousing, error)
	Upsert(ctx context.Context, housing *jobdomain.JobHousing) error
}

type JobOverallRatingRepository interface {
	FindByJobID(ctx context.Context, jobID string) (*jobdomain.JobOverallRating, error)
	Recalculate(ctx context.Context, jobID string) error
}

type JobReviewRepository interface {
	Insert(ctx context.Context, r *jobdomain.JobReview) error
	FindByJobID(ctx context.Context, jobID string) ([]jobdomain.JobReview, error)
}

type UserCartRepository interface {
	Insert(ctx context.Context, c *jobdomain.UserCart) error
	FindByUserAndJob(ctx context.Context, userID, jobID string) (*jobdomain.UserCart, error)
	FindByID(ctx context.Context, id string) (*jobdomain.UserCart, error)
	UpdateStatus(ctx context.Context, id string, status jobdomain.CartStatus) error
	FindByUser(ctx context.Context, userID string) ([]jobdomain.UserCart, error)
	Delete(ctx context.Context, id string) error
}

type RadarRepository interface {
	FindNearby(ctx context.Context, lat, lng, radius float64, staleMinutes int) ([]userdomain.Profile, error)
}

type NotificationRepository interface {
	Insert(ctx context.Context, n *notificationdomain.Notification) error
	FindByUser(ctx context.Context, userID string) ([]notificationdomain.Notification, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
}

type LeaderboardRepository interface {
	FindByScope(ctx context.Context, scope, jobID string) ([]userdomain.User, error)
}

type AdminRepository interface {
	GetStats(ctx context.Context) (*AdminStats, error)
	ListPendingVerifications(ctx context.Context) ([]missiondomain.UserMission, error)
	SearchUsers(ctx context.Context, query string) ([]userdomain.User, error)
}
