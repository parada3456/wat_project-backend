package port

import (
	"context"
	"time"
	"github.com/jackc/pgx/v5"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"
	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
)

type TxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type AdminStats struct {
	TotalUsers           int `json:"totalUsers"`
	ActiveUsers          int `json:"activeUsers"`
	PendingVerifications int `json:"pendingVerifications"`
	ActiveJobs           int `json:"activeJobs"`
	AverageCreditScore   int `json:"averageCreditScore"`
	TotalPointsAwarded   int `json:"totalPointsAwarded"`
}

type PointsAdjustmentResult struct {
	UserID               string `json:"userId"`
	LifetimeBalanceAfter int    `json:"lifetimeBalanceAfter"`
	PhaseBalanceAfter    int    `json:"phaseBalanceAfter"`
	LedgerID             string `json:"ledgerId"`
}

type UserWithProfile struct {
	User    userdomain.User    `json:"user"`
	Profile userdomain.Profile `json:"profile"`
}

type AdminRepository interface {
	GetStats(ctx context.Context) (*AdminStats, error)
	ListPendingVerifications(ctx context.Context, limit, offset int) ([]missiondomain.UserMission, int, error)
	SearchUsers(ctx context.Context, query string, limit, offset int) ([]UserWithProfile, int, error)
}

type UserRepository interface {
	IncrementPoints(ctx context.Context, userID string, lifetimeDelta, phaseDelta int) error
	FindByID(ctx context.Context, id string) (*userdomain.User, error)
}

type ProfileRepository interface {
	FindByUserID(ctx context.Context, userID string) (*userdomain.Profile, error)
}

type UserMissionRepository interface {
	FindByID(ctx context.Context, id string) (*missiondomain.UserMission, error)
	UpdateVerification(ctx context.Context, id string, verifiedAt time.Time, verifiedBy string) error
	UpdateReward(ctx context.Context, id string, reward *gamificationdomain.PointReward, rewardedAt time.Time) error
}

type MissionRepository interface {
	FindByID(ctx context.Context, id string) (*missiondomain.Mission, error)
}

type PointLedgerRepository interface {
	Insert(ctx context.Context, l *gamificationdomain.PointLedger) error
}

type NotifierPort interface {
	Send(ctx context.Context, userID, title, body string) error
}

type AdminUseCase interface {
	GetDashboardStats(ctx context.Context) (*AdminStats, error)
	ListPendingVerifications(ctx context.Context, page, pageSize int) ([]missiondomain.UserMission, int, error)
	VerifyMission(ctx context.Context, adminID, userMissionID string, approved bool, rejectionReason *string) (*missiondomain.UserMission, error)
	ListUsers(ctx context.Context, search string, page, pageSize int) ([]UserWithProfile, int, error)
	GetUserDetail(ctx context.Context, id string) (*UserWithProfile, error)
	AdjustPoints(ctx context.Context, userID string, delta int, reason string) (*PointsAdjustmentResult, error)
}
