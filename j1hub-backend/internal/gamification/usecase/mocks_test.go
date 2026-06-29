package gamificationusecase_test

import (
	authport "github.com/j1hub/backend/internal/auth/port"
	"context"
	"io"
	"time"

	frienddomain "github.com/j1hub/backend/internal/friend/domain"
	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	notificationdomain "github.com/j1hub/backend/internal/notification/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	expensedomain "github.com/j1hub/backend/internal/expense/domain"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) Create(ctx context.Context, u *userdomain.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*userdomain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userdomain.User), args.Error(1)
}
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userdomain.User), args.Error(1)
}
func (m *MockUserRepository) Update(ctx context.Context, u *userdomain.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *MockUserRepository) IncrementPoints(ctx context.Context, userID string, l, p int) error {
	return m.Called(ctx, userID, l, p).Error(0)
}
func (m *MockUserRepository) ResetStreak(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}
func (m *MockUserRepository) SetPhase(ctx context.Context, userID, phaseID string) error {
	return m.Called(ctx, userID, phaseID).Error(0)
}
func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// MockProfileRepository
type MockProfileRepository struct{ mock.Mock }

func (m *MockProfileRepository) Create(ctx context.Context, p *userdomain.Profile) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockProfileRepository) FindByUserID(ctx context.Context, userID string) (*userdomain.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userdomain.Profile), args.Error(1)
}
func (m *MockProfileRepository) Update(ctx context.Context, p *userdomain.Profile) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockProfileRepository) UpdateLocation(ctx context.Context, userID string, lat, lng float64) error {
	return m.Called(ctx, userID, lat, lng).Error(0)
}
func (m *MockProfileRepository) UpdateVisibility(ctx context.Context, userID string, visibility userdomain.RadarVisibility) error {
	return m.Called(ctx, userID, visibility).Error(0)
}

// MockCreditScoreRepository
type MockCreditScoreRepository struct{ mock.Mock }

func (m *MockCreditScoreRepository) Create(ctx context.Context, c *gamificationdomain.CreditScore) error {
	return m.Called(ctx, c).Error(0)
}
func (m *MockCreditScoreRepository) FindByUserID(ctx context.Context, userID string) (*gamificationdomain.CreditScore, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gamificationdomain.CreditScore), args.Error(1)
}
func (m *MockCreditScoreRepository) Decrement(ctx context.Context, userID string, delta int) error {
	return m.Called(ctx, userID, delta).Error(0)
}

// MockJourneyPhaseRepository
type MockJourneyPhaseRepository struct{ mock.Mock }

func (m *MockJourneyPhaseRepository) FindByNumber(ctx context.Context, number int) (*missiondomain.JourneyPhase, error) {
	args := m.Called(ctx, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*missiondomain.JourneyPhase), args.Error(1)
}
func (m *MockJourneyPhaseRepository) FindByID(ctx context.Context, id string) (*missiondomain.JourneyPhase, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*missiondomain.JourneyPhase), args.Error(1)
}
func (m *MockJourneyPhaseRepository) ListAll(ctx context.Context) ([]missiondomain.JourneyPhase, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.JourneyPhase), args.Error(1)
}

// MockUserPhaseHistoryRepository
type MockUserPhaseHistoryRepository struct{ mock.Mock }

func (m *MockUserPhaseHistoryRepository) Insert(ctx context.Context, h *missiondomain.UserPhaseHistory) error {
	return m.Called(ctx, h).Error(0)
}
func (m *MockUserPhaseHistoryRepository) CompleteCurrentPhase(ctx context.Context, userID string, points int, completedAt time.Time) error {
	return m.Called(ctx, userID, points, completedAt).Error(0)
}
func (m *MockUserPhaseHistoryRepository) FindByUserAndPhase(ctx context.Context, userID, phaseID string) (*missiondomain.UserPhaseHistory, error) {
	args := m.Called(ctx, userID, phaseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*missiondomain.UserPhaseHistory), args.Error(1)
}
func (m *MockUserPhaseHistoryRepository) FindByUser(ctx context.Context, userID string) ([]missiondomain.UserPhaseHistory, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.UserPhaseHistory), args.Error(1)
}

// MockMissionRepository
type MockMissionRepository struct{ mock.Mock }

func (m *MockMissionRepository) FindByPhase(ctx context.Context, phaseID string) ([]missiondomain.Mission, error) {
	args := m.Called(ctx, phaseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.Mission), args.Error(1)
}
func (m *MockMissionRepository) FindByID(ctx context.Context, id string) (*missiondomain.Mission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*missiondomain.Mission), args.Error(1)
}

// MockUserMissionRepository
type MockUserMissionRepository struct{ mock.Mock }

func (m *MockUserMissionRepository) BulkInsert(ctx context.Context, ums []missiondomain.UserMission) error {
	return m.Called(ctx, ums).Error(0)
}
func (m *MockUserMissionRepository) FindByUserAndPhase(ctx context.Context, userID, phaseID string) ([]missiondomain.UserMission, error) {
	args := m.Called(ctx, userID, phaseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.UserMission), args.Error(1)
}
func (m *MockUserMissionRepository) FindByID(ctx context.Context, id string) (*missiondomain.UserMission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*missiondomain.UserMission), args.Error(1)
}
func (m *MockUserMissionRepository) UpdateStatus(ctx context.Context, id string, status missiondomain.UserMissionStatus) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *MockUserMissionRepository) UpdateVerification(ctx context.Context, id string, verifiedAt time.Time, verifiedBy string) error {
	return m.Called(ctx, id, verifiedAt, verifiedBy).Error(0)
}
func (m *MockUserMissionRepository) UpdateReward(ctx context.Context, id string, reward *gamificationdomain.PointReward, rewardedAt time.Time) error {
	return m.Called(ctx, id, reward, rewardedAt).Error(0)
}
func (m *MockUserMissionRepository) FindOverdue(ctx context.Context) ([]missiondomain.UserMission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.UserMission), args.Error(1)
}

// MockTaskRepository
type MockTaskRepository struct{ mock.Mock }

func (m *MockTaskRepository) FindByMission(ctx context.Context, missionID string) ([]missiondomain.Task, error) {
	args := m.Called(ctx, missionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.Task), args.Error(1)
}

// MockUserTaskRepository
type MockUserTaskRepository struct{ mock.Mock }

func (m *MockUserTaskRepository) Upsert(ctx context.Context, ut *missiondomain.UserTask) error {
	return m.Called(ctx, ut).Error(0)
}
func (m *MockUserTaskRepository) FindByUserMission(ctx context.Context, userMissionID string) ([]missiondomain.UserTask, error) {
	args := m.Called(ctx, userMissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.UserTask), args.Error(1)
}

// MockPointLedgerRepository
type MockPointLedgerRepository struct{ mock.Mock }

func (m *MockPointLedgerRepository) Insert(ctx context.Context, l *gamificationdomain.PointLedger) error {
	return m.Called(ctx, l).Error(0)
}
func (m *MockPointLedgerRepository) InsertBatch(ctx context.Context, ledgers []gamificationdomain.PointLedger) error {
	return m.Called(ctx, ledgers).Error(0)
}
func (m *MockPointLedgerRepository) FindByUser(ctx context.Context, userID string) ([]gamificationdomain.PointLedger, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]gamificationdomain.PointLedger), args.Error(1)
}
func (m *MockPointLedgerRepository) FindByUserAndSourceType(ctx context.Context, userID string, sourceType gamificationdomain.SourceType) ([]gamificationdomain.PointLedger, error) {
	args := m.Called(ctx, userID, sourceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]gamificationdomain.PointLedger), args.Error(1)
}

// MockBadgeRepository
type MockBadgeRepository struct{ mock.Mock }

func (m *MockBadgeRepository) FindByTriggerType(ctx context.Context, triggerType gamificationdomain.TriggerType) ([]gamificationdomain.Badge, error) {
	args := m.Called(ctx, triggerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]gamificationdomain.Badge), args.Error(1)
}
func (m *MockBadgeRepository) FindEligible(ctx context.Context, userID string, triggerType gamificationdomain.TriggerType) ([]gamificationdomain.Badge, error) {
	args := m.Called(ctx, userID, triggerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]gamificationdomain.Badge), args.Error(1)
}

// MockUserBadgeRepository
type MockUserBadgeRepository struct{ mock.Mock }

func (m *MockUserBadgeRepository) Insert(ctx context.Context, ub *gamificationdomain.UserBadge) error {
	return m.Called(ctx, ub).Error(0)
}
func (m *MockUserBadgeRepository) FindByUser(ctx context.Context, userID string) ([]gamificationdomain.UserBadge, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]gamificationdomain.UserBadge), args.Error(1)
}

// MockFriendshipRepository
type MockFriendshipRepository struct{ mock.Mock }

func (m *MockFriendshipRepository) Insert(ctx context.Context, f *frienddomain.Friendship) error {
	return m.Called(ctx, f).Error(0)
}
func (m *MockFriendshipRepository) FindByCanonicalPair(ctx context.Context, u1, u2 string) (*frienddomain.Friendship, error) {
	args := m.Called(ctx, u1, u2)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*frienddomain.Friendship), args.Error(1)
}
func (m *MockFriendshipRepository) FindByID(ctx context.Context, id string) (*frienddomain.Friendship, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*frienddomain.Friendship), args.Error(1)
}
func (m *MockFriendshipRepository) UpdateStatus(ctx context.Context, id string, status frienddomain.FriendshipStatus) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *MockFriendshipRepository) FindFriendsOf(ctx context.Context, userID string, limit, offset int) ([]frienddomain.Friendship, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]frienddomain.Friendship), args.Int(1), args.Error(2)
}
func (m *MockFriendshipRepository) FindPendingFor(ctx context.Context, userID string, limit, offset int) ([]frienddomain.Friendship, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]frienddomain.Friendship), args.Int(1), args.Error(2)
}
func (m *MockFriendshipRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// MockExpenseTransactionRepository
type MockExpenseTransactionRepository struct{ mock.Mock }

func (m *MockExpenseTransactionRepository) Insert(ctx context.Context, t *expensedomain.ExpenseTransaction) error {
	return m.Called(ctx, t).Error(0)
}
func (m *MockExpenseTransactionRepository) FindByID(ctx context.Context, id string) (*expensedomain.ExpenseTransaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expensedomain.ExpenseTransaction), args.Error(1)
}
func (m *MockExpenseTransactionRepository) MarkSettled(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// MockExpenseSplitRepository
type MockExpenseSplitRepository struct{ mock.Mock }

func (m *MockExpenseSplitRepository) BulkInsert(ctx context.Context, splits []expensedomain.ExpenseSplit) error {
	return m.Called(ctx, splits).Error(0)
}
func (m *MockExpenseSplitRepository) FindByID(ctx context.Context, id string) (*expensedomain.ExpenseSplit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*expensedomain.ExpenseSplit), args.Error(1)
}
func (m *MockExpenseSplitRepository) UpdatePaymentStatus(ctx context.Context, id string, status expensedomain.PaymentStatus, slipURL string) error {
	return m.Called(ctx, id, status, slipURL).Error(0)
}
func (m *MockExpenseSplitRepository) UpdateApproval(ctx context.Context, id string, status expensedomain.ApprovalStatus, settledAt *time.Time) error {
	return m.Called(ctx, id, status, settledAt).Error(0)
}
func (m *MockExpenseSplitRepository) FindOverdue(ctx context.Context) ([]expensedomain.ExpenseSplit, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]expensedomain.ExpenseSplit), args.Error(1)
}
func (m *MockExpenseSplitRepository) CountUnsettled(ctx context.Context, transactionID string) (int, error) {
	args := m.Called(ctx, transactionID)
	return args.Int(0), args.Error(1)
}

// MockJobPostingRepository
type MockJobPostingRepository struct{ mock.Mock }

func (m *MockJobPostingRepository) FindWithFilters(ctx context.Context, filters map[string]interface{}) ([]jobdomain.JobPosting, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jobdomain.JobPosting), args.Error(1)
}
func (m *MockJobPostingRepository) FindByID(ctx context.Context, id string) (*jobdomain.JobPosting, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jobdomain.JobPosting), args.Error(1)
}
func (m *MockJobPostingRepository) Upsert(ctx context.Context, job *jobdomain.JobPosting) error {
	return m.Called(ctx, job).Error(0)
}

// MockJobHousingRepository
type MockJobHousingRepository struct{ mock.Mock }

func (m *MockJobHousingRepository) FindByJobID(ctx context.Context, jobID string) ([]jobdomain.JobHousing, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jobdomain.JobHousing), args.Error(1)
}
func (m *MockJobHousingRepository) Upsert(ctx context.Context, housing *jobdomain.JobHousing) error {
	return m.Called(ctx, housing).Error(0)
}

// MockJobOverallRatingRepository
type MockJobOverallRatingRepository struct{ mock.Mock }

func (m *MockJobOverallRatingRepository) FindByJobID(ctx context.Context, jobID string) (*jobdomain.JobOverallRating, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jobdomain.JobOverallRating), args.Error(1)
}
func (m *MockJobOverallRatingRepository) Recalculate(ctx context.Context, jobID string) error {
	return m.Called(ctx, jobID).Error(0)
}

// MockJobReviewRepository
type MockJobReviewRepository struct{ mock.Mock }

func (m *MockJobReviewRepository) Insert(ctx context.Context, r *jobdomain.JobReview) error {
	return m.Called(ctx, r).Error(0)
}
func (m *MockJobReviewRepository) FindByJobID(ctx context.Context, jobID string) ([]jobdomain.JobReview, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jobdomain.JobReview), args.Error(1)
}

// MockUserCartRepository
type MockUserCartRepository struct{ mock.Mock }

func (m *MockUserCartRepository) Insert(ctx context.Context, c *jobdomain.UserCart) error {
	return m.Called(ctx, c).Error(0)
}
func (m *MockUserCartRepository) FindByUserAndJob(ctx context.Context, userID, jobID string) (*jobdomain.UserCart, error) {
	args := m.Called(ctx, userID, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jobdomain.UserCart), args.Error(1)
}
func (m *MockUserCartRepository) FindByID(ctx context.Context, id string) (*jobdomain.UserCart, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jobdomain.UserCart), args.Error(1)
}
func (m *MockUserCartRepository) UpdateStatus(ctx context.Context, id string, status jobdomain.CartStatus) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *MockUserCartRepository) FindByUser(ctx context.Context, userID string) ([]jobdomain.UserCart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jobdomain.UserCart), args.Error(1)
}
func (m *MockUserCartRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// MockRadarRepository
type MockRadarRepository struct{ mock.Mock }

func (m *MockRadarRepository) FindNearby(ctx context.Context, lat, lng, radius float64, staleMinutes int) ([]userdomain.Profile, error) {
	args := m.Called(ctx, lat, lng, radius, staleMinutes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdomain.Profile), args.Error(1)
}

// MockNotificationRepository
type MockNotificationRepository struct{ mock.Mock }

func (m *MockNotificationRepository) Insert(ctx context.Context, n *notificationdomain.Notification) error {
	return m.Called(ctx, n).Error(0)
}
func (m *MockNotificationRepository) FindByUser(ctx context.Context, userID string) ([]notificationdomain.Notification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]notificationdomain.Notification), args.Error(1)
}
func (m *MockNotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockNotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}
func (m *MockNotificationRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// MockLeaderboardRepository
type MockLeaderboardRepository struct{ mock.Mock }

func (m *MockLeaderboardRepository) FindByScope(ctx context.Context, scope, jobID string) ([]userdomain.User, error) {
	args := m.Called(ctx, scope, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdomain.User), args.Error(1)
}

// MockPasswordHasher
type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) Hash(plain string) (string, error) {
	args := m.Called(plain)
	return args.String(0), args.Error(1)
}
func (m *MockHasher) Verify(plain, hash string) bool {
	args := m.Called(plain, hash)
	return args.Bool(0)
}

// MockTokenIssuer
type MockIssuer struct {
	mock.Mock
}

func (m *MockIssuer) Issue(userID string, isAdmin bool) (*authport.TokenPair, error) {
	args := m.Called(userID, isAdmin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authport.TokenPair), args.Error(1)
}
func (m *MockIssuer) Verify(token string) (*authport.Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authport.Claims), args.Error(1)
}
func (m *MockIssuer) Refresh(refreshToken string) (*authport.TokenPair, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authport.TokenPair), args.Error(1)
}

// MockStoragePort
type MockStoragePort struct{ mock.Mock }

func (m *MockStoragePort) UploadFile(ctx context.Context, bucket, key string, data io.Reader, contentType string) (string, error) {
	args := m.Called(ctx, bucket, key, data, contentType)
	return args.String(0), args.Error(1)
}

// MockNotifierPort
type MockNotifierPort struct{ mock.Mock }

func (m *MockNotifierPort) Send(ctx context.Context, userID, title, body string) error {
	return m.Called(ctx, userID, title, body).Error(0)
}

// MockClock
type MockClock struct {
	NowTime time.Time
}

func (m *MockClock) Now() time.Time {
	return m.NowTime
}

// MockTx
type MockTx struct {
	pgx.Tx
	mock.Mock
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockTxBeginner
type MockTxBeginner struct {
	mock.Mock
}

func (m *MockTxBeginner) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(pgx.Tx), args.Error(1)
}
