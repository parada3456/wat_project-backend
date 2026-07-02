package port

import (
	"context"
	"time"
	"github.com/jackc/pgx/v5"
	userdomain "github.com/j1hub/backend/internal/user/domain"
	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
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

type UserRepository interface {
	Create(ctx context.Context, u *userdomain.User) error
	FindByID(ctx context.Context, id string) (*userdomain.User, error)
	FindByEmail(ctx context.Context, email string) (*userdomain.User, error)
	Update(ctx context.Context, u *userdomain.User) error
}

type ProfileRepository interface {
	Create(ctx context.Context, p *userdomain.Profile) error
	FindByUserID(ctx context.Context, userID string) (*userdomain.Profile, error)
}

type CreditScoreRepository interface {
	Create(ctx context.Context, c *gamificationdomain.CreditScore) error
}

type JourneyPhaseRepository interface {
	FindByNumber(ctx context.Context, number int) (*missiondomain.JourneyPhase, error)
}

type UserPhaseHistoryRepository interface {
	Insert(ctx context.Context, h *missiondomain.UserPhaseHistory) error
}

type MissionRepository interface {
	FindByPhase(ctx context.Context, phaseID string) ([]missiondomain.Mission, error)
}

type UserMissionRepository interface {
	BulkInsert(ctx context.Context, ums []missiondomain.UserMission) error
}

type TxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}
