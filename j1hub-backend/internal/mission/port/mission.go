package port

import (
	"context"
	"time"
	"io"
	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"
)

type JourneyPhaseRepository interface {
	FindByNumber(ctx context.Context, number int) (*missiondomain.JourneyPhase, error)
	FindByID(ctx context.Context, id string) (*missiondomain.JourneyPhase, error)
	ListAll(ctx context.Context) ([]missiondomain.JourneyPhase, error)
}

type UserPhaseHistoryRepository interface {
	Insert(ctx context.Context, h *missiondomain.UserPhaseHistory) error
	CompleteCurrentPhase(ctx context.Context, userID string, points int, completedAt time.Time) error
	FindByUserAndPhase(ctx context.Context, userID, phaseID string) (*missiondomain.UserPhaseHistory, error)
	FindByUser(ctx context.Context, userID string) ([]missiondomain.UserPhaseHistory, error)
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
	FindByID(ctx context.Context, id string) (*missiondomain.UserTask, error)
}

type UserRepository interface {
	IncrementPoints(ctx context.Context, userID string, lifetimeDelta, phaseDelta int) error
	FindByID(ctx context.Context, id string) (*userdomain.User, error)
	Update(ctx context.Context, u *userdomain.User) error
}

type PointLedgerRepository interface {
	Insert(ctx context.Context, ledger *gamificationdomain.PointLedger) error
	InsertBatch(ctx context.Context, ledgers []gamificationdomain.PointLedger) error
}

type BadgeRepository interface {
	FindEligible(ctx context.Context, userID string, triggerType gamificationdomain.TriggerType) ([]gamificationdomain.Badge, error)
}

type UserBadgeRepository interface {
	Insert(ctx context.Context, ub *gamificationdomain.UserBadge) error
}

type StoragePort interface {
	UploadFile(ctx context.Context, bucket, key string, data io.Reader, contentType string) (url string, err error)
}

type NotifierPort interface {
	Send(ctx context.Context, userID, title, body string) error
}
