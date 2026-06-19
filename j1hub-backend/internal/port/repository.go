package port

import (
	"context"
	"time"

	"github.com/j1hub/backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	IncrementPoints(ctx context.Context, userID string, lifetimeDelta, phaseDelta int) error
	ResetStreak(ctx context.Context, userID string) error
	SetPhase(ctx context.Context, userID, phaseID string) error
	Delete(ctx context.Context, id string) error
}

type ProfileRepository interface {
	Create(ctx context.Context, p *domain.Profile) error
	FindByUserID(ctx context.Context, userID string) (*domain.Profile, error)
	Update(ctx context.Context, p *domain.Profile) error
	UpdateLocation(ctx context.Context, userID string, lat, lng float64) error
	UpdateVisibility(ctx context.Context, userID string, visibility domain.RadarVisibility) error
}

type CreditScoreRepository interface {
	Create(ctx context.Context, c *domain.CreditScore) error
	FindByUserID(ctx context.Context, userID string) (*domain.CreditScore, error)
	Decrement(ctx context.Context, userID string, delta int) error
}

type JourneyPhaseRepository interface {
	FindByNumber(ctx context.Context, number int) (*domain.JourneyPhase, error)
	FindByID(ctx context.Context, id string) (*domain.JourneyPhase, error)
}

type UserPhaseHistoryRepository interface {
	Insert(ctx context.Context, h *domain.UserPhaseHistory) error
	CompleteCurrentPhase(ctx context.Context, userID string, points int, completedAt time.Time) error
	FindByUserAndPhase(ctx context.Context, userID, phaseID string) (*domain.UserPhaseHistory, error)
}

type MissionRepository interface {
	FindByPhase(ctx context.Context, phaseID string) ([]domain.Mission, error)
	FindByID(ctx context.Context, id string) (*domain.Mission, error)
}

type UserMissionRepository interface {
	BulkInsert(ctx context.Context, ums []domain.UserMission) error
	FindByUserAndPhase(ctx context.Context, userID, phaseID string) ([]domain.UserMission, error)
	FindByID(ctx context.Context, id string) (*domain.UserMission, error)
	UpdateStatus(ctx context.Context, id string, status domain.UserMissionStatus) error
	UpdateVerification(ctx context.Context, id string, verifiedAt time.Time, verifiedBy string) error
	UpdateReward(ctx context.Context, id string, reward *domain.PointReward, rewardedAt time.Time) error
	FindOverdue(ctx context.Context) ([]domain.UserMission, error)
}

type TaskRepository interface {
	FindByMission(ctx context.Context, missionID string) ([]domain.Task, error)
}

type UserTaskRepository interface {
	Upsert(ctx context.Context, ut *domain.UserTask) error
	FindByUserMission(ctx context.Context, userMissionID string) ([]domain.UserTask, error)
}

type PointLedgerRepository interface {
	Insert(ctx context.Context, l *domain.PointLedger) error
	InsertBatch(ctx context.Context, ledgers []domain.PointLedger) error
}

type BadgeRepository interface {
	FindByTriggerType(ctx context.Context, triggerType domain.TriggerType) ([]domain.Badge, error)
	FindEligible(ctx context.Context, userID string, triggerType domain.TriggerType) ([]domain.Badge, error)
}

type UserBadgeRepository interface {
	Insert(ctx context.Context, ub *domain.UserBadge) error
	FindByUser(ctx context.Context, userID string) ([]domain.UserBadge, error)
}

type FriendshipRepository interface {
	Insert(ctx context.Context, f *domain.Friendship) error
	FindByCanonicalPair(ctx context.Context, u1, u2 string) (*domain.Friendship, error)
	FindByID(ctx context.Context, id string) (*domain.Friendship, error)
	UpdateStatus(ctx context.Context, id string, status domain.FriendshipStatus) error
	FindFriendsOf(ctx context.Context, userID string) ([]domain.Friendship, error)
}

type ExpenseTransactionRepository interface {
	Insert(ctx context.Context, t *domain.ExpenseTransaction) error
	FindByID(ctx context.Context, id string) (*domain.ExpenseTransaction, error)
	MarkSettled(ctx context.Context, id string) error
}

type ExpenseSplitRepository interface {
	BulkInsert(ctx context.Context, splits []domain.ExpenseSplit) error
	FindByID(ctx context.Context, id string) (*domain.ExpenseSplit, error)
	UpdatePaymentStatus(ctx context.Context, id string, status domain.PaymentStatus, slipURL string) error
	UpdateApproval(ctx context.Context, id string, status domain.ApprovalStatus, settledAt *time.Time) error
	FindOverdue(ctx context.Context) ([]domain.ExpenseSplit, error)
	CountUnsettled(ctx context.Context, transactionID string) (int, error)
}

type JobPostingRepository interface {
	FindWithFilters(ctx context.Context, filters map[string]interface{}) ([]domain.JobPosting, error)
	FindByID(ctx context.Context, id string) (*domain.JobPosting, error)
	Upsert(ctx context.Context, job *domain.JobPosting) error
}

type JobHousingRepository interface {
	FindByJobID(ctx context.Context, jobID string) ([]domain.JobHousing, error)
	Upsert(ctx context.Context, housing *domain.JobHousing) error
}

type JobOverallRatingRepository interface {
	FindByJobID(ctx context.Context, jobID string) (*domain.JobOverallRating, error)
	Recalculate(ctx context.Context, jobID string) error
}

type JobReviewRepository interface {
	Insert(ctx context.Context, r *domain.JobReview) error
	FindByJobID(ctx context.Context, jobID string) ([]domain.JobReview, error)
}

type UserCartRepository interface {
	Insert(ctx context.Context, c *domain.UserCart) error
	FindByUserAndJob(ctx context.Context, userID, jobID string) (*domain.UserCart, error)
	FindByID(ctx context.Context, id string) (*domain.UserCart, error)
	UpdateStatus(ctx context.Context, id string, status domain.CartStatus) error
}

type RadarRepository interface {
	FindNearby(ctx context.Context, lat, lng, radius float64, staleMinutes int) ([]domain.Profile, error)
}

type NotificationRepository interface {
	Insert(ctx context.Context, n *domain.Notification) error
	FindByUser(ctx context.Context, userID string) ([]domain.Notification, error)
	MarkAsRead(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	Delete(ctx context.Context, id string) error
}

type LeaderboardRepository interface {
	FindByScope(ctx context.Context, scope, jobID string) ([]domain.User, error)
}

type AdminRepository interface {
	GetStats(ctx context.Context) (*AdminStats, error)
	ListPendingVerifications(ctx context.Context) ([]domain.UserMission, error)
	SearchUsers(ctx context.Context, query string) ([]domain.User, error)
}

