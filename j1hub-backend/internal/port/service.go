package port

import (
	"context"
	"io"
	"time"

	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"
)

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(plain, hash string) bool
}

type Claims struct {
	UserID  string
	IsAdmin bool
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type TokenIssuer interface {
	Issue(userID string, isAdmin bool) (*TokenPair, error)
	Verify(token string) (*Claims, error)
	Refresh(refreshToken string) (*TokenPair, error)
}

type StoragePort interface {
	UploadFile(ctx context.Context, bucket, key string, data io.Reader, contentType string) (url string, err error)
}

type NotifierPort interface {
	Send(ctx context.Context, userID, title, body string) error
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

type AdminUseCase interface {
	GetDashboardStats(ctx context.Context) (*AdminStats, error)
	ListPendingVerifications(ctx context.Context) ([]missiondomain.UserMission, error)
	VerifyMission(ctx context.Context, adminID, userMissionID string, approved bool, rejectionReason *string) (*missiondomain.UserMission, error)
	ListUsers(ctx context.Context, search string) ([]userdomain.User, error)
	GetUserDetail(ctx context.Context, id string) (*userdomain.User, error)
	AdjustPoints(ctx context.Context, userID string, delta int, reason string) (*PointsAdjustmentResult, error)
}
