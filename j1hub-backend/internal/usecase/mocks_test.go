package usecase_test

import (
	"context"
	"io"
	"time"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository
type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) Create(ctx context.Context, u *domain.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *MockUserRepository) Update(ctx context.Context, u *domain.User) error {
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

func (m *MockProfileRepository) Create(ctx context.Context, p *domain.Profile) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockProfileRepository) FindByUserID(ctx context.Context, userID string) (*domain.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Profile), args.Error(1)
}
func (m *MockProfileRepository) Update(ctx context.Context, p *domain.Profile) error {
	return m.Called(ctx, p).Error(0)
}
func (m *MockProfileRepository) UpdateLocation(ctx context.Context, userID string, lat, lng float64) error {
	return m.Called(ctx, userID, lat, lng).Error(0)
}
func (m *MockProfileRepository) UpdateVisibility(ctx context.Context, userID string, visibility domain.RadarVisibility) error {
	return m.Called(ctx, userID, visibility).Error(0)
}

// MockCreditScoreRepository
type MockCreditScoreRepository struct{ mock.Mock }

func (m *MockCreditScoreRepository) Create(ctx context.Context, c *domain.CreditScore) error {
	return m.Called(ctx, c).Error(0)
}
func (m *MockCreditScoreRepository) FindByUserID(ctx context.Context, userID string) (*domain.CreditScore, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CreditScore), args.Error(1)
}
func (m *MockCreditScoreRepository) Decrement(ctx context.Context, userID string, delta int) error {
	return m.Called(ctx, userID, delta).Error(0)
}

// MockJourneyPhaseRepository
type MockJourneyPhaseRepository struct{ mock.Mock }

func (m *MockJourneyPhaseRepository) FindByNumber(ctx context.Context, number int) (*domain.JourneyPhase, error) {
	args := m.Called(ctx, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JourneyPhase), args.Error(1)
}
func (m *MockJourneyPhaseRepository) FindByID(ctx context.Context, id string) (*domain.JourneyPhase, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JourneyPhase), args.Error(1)
}

// MockUserPhaseHistoryRepository
type MockUserPhaseHistoryRepository struct{ mock.Mock }

func (m *MockUserPhaseHistoryRepository) Insert(ctx context.Context, h *domain.UserPhaseHistory) error {
	return m.Called(ctx, h).Error(0)
}
func (m *MockUserPhaseHistoryRepository) CompleteCurrentPhase(ctx context.Context, userID string, points int, completedAt time.Time) error {
	return m.Called(ctx, userID, points, completedAt).Error(0)
}
func (m *MockUserPhaseHistoryRepository) FindByUserAndPhase(ctx context.Context, userID, phaseID string) (*domain.UserPhaseHistory, error) {
	args := m.Called(ctx, userID, phaseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserPhaseHistory), args.Error(1)
}

// MockMissionRepository
type MockMissionRepository struct{ mock.Mock }

func (m *MockMissionRepository) FindByPhase(ctx context.Context, phaseID string) ([]domain.Mission, error) {
	args := m.Called(ctx, phaseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Mission), args.Error(1)
}
func (m *MockMissionRepository) FindByID(ctx context.Context, id string) (*domain.Mission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Mission), args.Error(1)
}

// MockUserMissionRepository
type MockUserMissionRepository struct{ mock.Mock }

func (m *MockUserMissionRepository) BulkInsert(ctx context.Context, ums []domain.UserMission) error {
	return m.Called(ctx, ums).Error(0)
}
func (m *MockUserMissionRepository) FindByUserAndPhase(ctx context.Context, userID, phaseID string) ([]domain.UserMission, error) {
	args := m.Called(ctx, userID, phaseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserMission), args.Error(1)
}
func (m *MockUserMissionRepository) FindByID(ctx context.Context, id string) (*domain.UserMission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserMission), args.Error(1)
}
func (m *MockUserMissionRepository) UpdateStatus(ctx context.Context, id string, status domain.UserMissionStatus) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *MockUserMissionRepository) UpdateVerification(ctx context.Context, id string, verifiedAt time.Time, verifiedBy string) error {
	return m.Called(ctx, id, verifiedAt, verifiedBy).Error(0)
}
func (m *MockUserMissionRepository) UpdateReward(ctx context.Context, id string, reward *domain.PointReward, rewardedAt time.Time) error {
	return m.Called(ctx, id, reward, rewardedAt).Error(0)
}
func (m *MockUserMissionRepository) FindOverdue(ctx context.Context) ([]domain.UserMission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserMission), args.Error(1)
}

// MockTaskRepository
type MockTaskRepository struct{ mock.Mock }

func (m *MockTaskRepository) FindByMission(ctx context.Context, missionID string) ([]domain.Task, error) {
	args := m.Called(ctx, missionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Task), args.Error(1)
}

// MockUserTaskRepository
type MockUserTaskRepository struct{ mock.Mock }

func (m *MockUserTaskRepository) Upsert(ctx context.Context, ut *domain.UserTask) error {
	return m.Called(ctx, ut).Error(0)
}
func (m *MockUserTaskRepository) FindByUserMission(ctx context.Context, userMissionID string) ([]domain.UserTask, error) {
	args := m.Called(ctx, userMissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserTask), args.Error(1)
}

// MockPointLedgerRepository
type MockPointLedgerRepository struct{ mock.Mock }

func (m *MockPointLedgerRepository) Insert(ctx context.Context, l *domain.PointLedger) error {
	return m.Called(ctx, l).Error(0)
}
func (m *MockPointLedgerRepository) InsertBatch(ctx context.Context, ledgers []domain.PointLedger) error {
	return m.Called(ctx, ledgers).Error(0)
}

// MockBadgeRepository
type MockBadgeRepository struct{ mock.Mock }

func (m *MockBadgeRepository) FindByTriggerType(ctx context.Context, triggerType domain.TriggerType) ([]domain.Badge, error) {
	args := m.Called(ctx, triggerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Badge), args.Error(1)
}
func (m *MockBadgeRepository) FindEligible(ctx context.Context, userID string, triggerType domain.TriggerType) ([]domain.Badge, error) {
	args := m.Called(ctx, userID, triggerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Badge), args.Error(1)
}

// MockUserBadgeRepository
type MockUserBadgeRepository struct{ mock.Mock }

func (m *MockUserBadgeRepository) Insert(ctx context.Context, ub *domain.UserBadge) error {
	return m.Called(ctx, ub).Error(0)
}
func (m *MockUserBadgeRepository) FindByUser(ctx context.Context, userID string) ([]domain.UserBadge, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserBadge), args.Error(1)
}

// MockFriendshipRepository
type MockFriendshipRepository struct{ mock.Mock }

func (m *MockFriendshipRepository) Insert(ctx context.Context, f *domain.Friendship) error {
	return m.Called(ctx, f).Error(0)
}
func (m *MockFriendshipRepository) FindByCanonicalPair(ctx context.Context, u1, u2 string) (*domain.Friendship, error) {
	args := m.Called(ctx, u1, u2)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Friendship), args.Error(1)
}
func (m *MockFriendshipRepository) FindByID(ctx context.Context, id string) (*domain.Friendship, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Friendship), args.Error(1)
}
func (m *MockFriendshipRepository) UpdateStatus(ctx context.Context, id string, status domain.FriendshipStatus) error {
	return m.Called(ctx, id, status).Error(0)
}
func (m *MockFriendshipRepository) FindFriendsOf(ctx context.Context, userID string) ([]domain.Friendship, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Friendship), args.Error(1)
}

// MockExpenseTransactionRepository
type MockExpenseTransactionRepository struct{ mock.Mock }

func (m *MockExpenseTransactionRepository) Insert(ctx context.Context, t *domain.ExpenseTransaction) error {
	return m.Called(ctx, t).Error(0)
}
func (m *MockExpenseTransactionRepository) FindByID(ctx context.Context, id string) (*domain.ExpenseTransaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ExpenseTransaction), args.Error(1)
}
func (m *MockExpenseTransactionRepository) MarkSettled(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// MockExpenseSplitRepository
type MockExpenseSplitRepository struct{ mock.Mock }

func (m *MockExpenseSplitRepository) BulkInsert(ctx context.Context, splits []domain.ExpenseSplit) error {
	return m.Called(ctx, splits).Error(0)
}
func (m *MockExpenseSplitRepository) FindByID(ctx context.Context, id string) (*domain.ExpenseSplit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ExpenseSplit), args.Error(1)
}
func (m *MockExpenseSplitRepository) UpdatePaymentStatus(ctx context.Context, id string, status domain.PaymentStatus, slipURL string) error {
	return m.Called(ctx, id, status, slipURL).Error(0)
}
func (m *MockExpenseSplitRepository) UpdateApproval(ctx context.Context, id string, status domain.ApprovalStatus, settledAt *time.Time) error {
	return m.Called(ctx, id, status, settledAt).Error(0)
}
func (m *MockExpenseSplitRepository) FindOverdue(ctx context.Context) ([]domain.ExpenseSplit, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ExpenseSplit), args.Error(1)
}
func (m *MockExpenseSplitRepository) CountUnsettled(ctx context.Context, transactionID string) (int, error) {
	args := m.Called(ctx, transactionID)
	return args.Int(0), args.Error(1)
}

// MockJobPostingRepository
type MockJobPostingRepository struct{ mock.Mock }

func (m *MockJobPostingRepository) FindWithFilters(ctx context.Context, filters map[string]interface{}) ([]domain.JobPosting, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.JobPosting), args.Error(1)
}
func (m *MockJobPostingRepository) FindByID(ctx context.Context, id string) (*domain.JobPosting, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JobPosting), args.Error(1)
}
func (m *MockJobPostingRepository) Upsert(ctx context.Context, job *domain.JobPosting) error {
	return m.Called(ctx, job).Error(0)
}

// MockJobHousingRepository
type MockJobHousingRepository struct{ mock.Mock }

func (m *MockJobHousingRepository) FindByJobID(ctx context.Context, jobID string) ([]domain.JobHousing, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.JobHousing), args.Error(1)
}
func (m *MockJobHousingRepository) Upsert(ctx context.Context, housing *domain.JobHousing) error {
	return m.Called(ctx, housing).Error(0)
}

// MockJobOverallRatingRepository
type MockJobOverallRatingRepository struct{ mock.Mock }

func (m *MockJobOverallRatingRepository) FindByJobID(ctx context.Context, jobID string) (*domain.JobOverallRating, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.JobOverallRating), args.Error(1)
}
func (m *MockJobOverallRatingRepository) Recalculate(ctx context.Context, jobID string) error {
	return m.Called(ctx, jobID).Error(0)
}

// MockJobReviewRepository
type MockJobReviewRepository struct{ mock.Mock }

func (m *MockJobReviewRepository) Insert(ctx context.Context, r *domain.JobReview) error {
	return m.Called(ctx, r).Error(0)
}
func (m *MockJobReviewRepository) FindByJobID(ctx context.Context, jobID string) ([]domain.JobReview, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.JobReview), args.Error(1)
}

// MockUserCartRepository
type MockUserCartRepository struct{ mock.Mock }

func (m *MockUserCartRepository) Insert(ctx context.Context, c *domain.UserCart) error {
	return m.Called(ctx, c).Error(0)
}
func (m *MockUserCartRepository) FindByUserAndJob(ctx context.Context, userID, jobID string) (*domain.UserCart, error) {
	args := m.Called(ctx, userID, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserCart), args.Error(1)
}
func (m *MockUserCartRepository) FindByID(ctx context.Context, id string) (*domain.UserCart, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserCart), args.Error(1)
}
func (m *MockUserCartRepository) UpdateStatus(ctx context.Context, id string, status domain.CartStatus) error {
	return m.Called(ctx, id, status).Error(0)
}

// MockRadarRepository
type MockRadarRepository struct{ mock.Mock }

func (m *MockRadarRepository) FindNearby(ctx context.Context, lat, lng, radius float64, staleMinutes int) ([]domain.Profile, error) {
	args := m.Called(ctx, lat, lng, radius, staleMinutes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Profile), args.Error(1)
}

// MockNotificationRepository
type MockNotificationRepository struct{ mock.Mock }

func (m *MockNotificationRepository) Insert(ctx context.Context, n *domain.Notification) error {
	return m.Called(ctx, n).Error(0)
}
func (m *MockNotificationRepository) FindByUser(ctx context.Context, userID string) ([]domain.Notification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
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

func (m *MockLeaderboardRepository) FindByScope(ctx context.Context, scope, jobID string) ([]domain.User, error) {
	args := m.Called(ctx, scope, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
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

func (m *MockIssuer) Issue(userID string, isAdmin bool) (*port.TokenPair, error) {
	args := m.Called(userID, isAdmin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.TokenPair), args.Error(1)
}
func (m *MockIssuer) Verify(token string) (*port.Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.Claims), args.Error(1)
}
func (m *MockIssuer) Refresh(refreshToken string) (*port.TokenPair, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.TokenPair), args.Error(1)
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
