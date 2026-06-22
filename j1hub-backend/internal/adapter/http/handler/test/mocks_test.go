package handler_test
import (
	"context"
	"io"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// MockJourneyUC
type MockJourneyUC struct {
	mock.Mock
}

func (m *MockJourneyUC) ListPhases(ctx context.Context) ([]domain.JourneyPhase, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.JourneyPhase), args.Error(1)
}

func (m *MockJourneyUC) GetHistory(ctx context.Context, userID string) ([]domain.UserPhaseHistory, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserPhaseHistory), args.Error(1)
}

func (m *MockJourneyUC) ListUserBadges(ctx context.Context, userID string) ([]domain.UserBadge, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserBadge), args.Error(1)
}

func (m *MockJourneyUC) GetCreditScoreHistory(ctx context.Context, userID string) ([]domain.PointLedger, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PointLedger), args.Error(1)
}

func (m *MockJourneyUC) GetPointsLedger(ctx context.Context, userID string) ([]domain.PointLedger, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PointLedger), args.Error(1)
}

// MockAdvancePhaseUC
type MockAdvancePhaseUC struct {
	mock.Mock
}

func (m *MockAdvancePhaseUC) TryAdvancePhase(ctx context.Context, userID string) (*usecase.PhaseTransitionResponse, error) {
	args := m.Called(ctx, userID)
	var resp *usecase.PhaseTransitionResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(*usecase.PhaseTransitionResponse)
	}
	return resp, args.Error(1)
}

// MockLeaderboardUC
type MockLeaderboardUC struct {
	mock.Mock
}

func (m *MockLeaderboardUC) GetLeaderboard(ctx context.Context, scope, jobID string) ([]usecase.LeaderboardEntry, error) {
	args := m.Called(ctx, scope, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]usecase.LeaderboardEntry), args.Error(1)
}

// MockManageExpenseUC
type MockManageExpenseUC struct {
	mock.Mock
}

func (m *MockManageExpenseUC) ListExpenses(ctx context.Context, userID string) ([]domain.ExpenseTransaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ExpenseTransaction), args.Error(1)
}

func (m *MockManageExpenseUC) CreateExpense(ctx context.Context, userID string, cmd usecase.CreateExpenseCmd) error {
	return m.Called(ctx, userID, cmd).Error(0)
}

func (m *MockManageExpenseUC) GetExpenseDetail(ctx context.Context, userID, id string) (*domain.ExpenseTransaction, []domain.ExpenseSplit, error) {
	args := m.Called(ctx, userID, id)
	var t *domain.ExpenseTransaction
	if args.Get(0) != nil {
		t = args.Get(0).(*domain.ExpenseTransaction)
	}
	var s []domain.ExpenseSplit
	if args.Get(1) != nil {
		s = args.Get(1).([]domain.ExpenseSplit)
	}
	return t, s, args.Error(2)
}

func (m *MockManageExpenseUC) DeleteExpense(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func (m *MockManageExpenseUC) ListPendingExpenses(ctx context.Context, userID string) ([]domain.ExpenseSplit, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ExpenseSplit), args.Error(1)
}

func (m *MockManageExpenseUC) SubmitSlip(ctx context.Context, debtorID, splitID string, file io.Reader, contentType string) error {
	return m.Called(ctx, debtorID, splitID, mock.Anything, contentType).Error(0)
}

func (m *MockManageExpenseUC) ApproveSplit(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

// MockFriendshipUC
type MockFriendshipUC struct {
	mock.Mock
}

func (m *MockFriendshipUC) SendRequest(ctx context.Context, senderID, targetID string) error {
	return m.Called(ctx, senderID, targetID).Error(0)
}

func (m *MockFriendshipUC) ListPendingRequests(ctx context.Context, userID string) ([]domain.Friendship, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Friendship), args.Error(1)
}

func (m *MockFriendshipUC) RespondToRequest(ctx context.Context, userID, friendshipID string, accept bool) error {
	return m.Called(ctx, userID, friendshipID, accept).Error(0)
}

func (m *MockFriendshipUC) ListFriends(ctx context.Context, userID string) ([]domain.Friendship, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Friendship), args.Error(1)
}

func (m *MockFriendshipUC) RemoveFriend(ctx context.Context, userID, friendID string) error {
	return m.Called(ctx, userID, friendID).Error(0)
}

// MockRadarUC
type MockRadarUC struct {
	mock.Mock
}

func (m *MockRadarUC) GetRadar(ctx context.Context, requesterID string) ([]usecase.RadarEntry, error) {
	args := m.Called(ctx, requesterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]usecase.RadarEntry), args.Error(1)
}

// MockJobUC
type MockJobUC struct {
	mock.Mock
}

func (m *MockJobUC) ListJobs(ctx context.Context, filters map[string]interface{}) ([]domain.JobPosting, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.JobPosting), args.Error(1)
}

func (m *MockJobUC) GetJobDetail(ctx context.Context, id string) (*domain.JobPosting, []domain.JobHousing, *domain.JobOverallRating, error) {
	args := m.Called(ctx, id)
	var j *domain.JobPosting
	if args.Get(0) != nil {
		j = args.Get(0).(*domain.JobPosting)
	}
	var h []domain.JobHousing
	if args.Get(1) != nil {
		h = args.Get(1).([]domain.JobHousing)
	}
	var r *domain.JobOverallRating
	if args.Get(2) != nil {
		r = args.Get(2).(*domain.JobOverallRating)
	}
	return j, h, r, args.Error(3)
}

func (m *MockJobUC) AddToCart(ctx context.Context, userID, jobID string) error {
	return m.Called(ctx, userID, jobID).Error(0)
}

func (m *MockJobUC) ListCart(ctx context.Context, userID string) ([]domain.UserCart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserCart), args.Error(1)
}

func (m *MockJobUC) RemoveFromCart(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func (m *MockJobUC) WriteReview(ctx context.Context, userID, jobID string, rev *domain.JobReview) error {
	return m.Called(ctx, userID, jobID, rev).Error(0)
}

func (m *MockJobUC) ListReviews(ctx context.Context, jobID string) ([]domain.JobReview, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.JobReview), args.Error(1)
}

func (m *MockJobUC) UpdateCartStatus(ctx context.Context, userID, cartID string, status domain.CartStatus) error {
	return m.Called(ctx, userID, cartID, status).Error(0)
}

// MockMissionUC
type MockMissionUC struct {
	mock.Mock
}

func (m *MockMissionUC) ListAvailableMissions(ctx context.Context, userID string) ([]domain.UserMission, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserMission), args.Error(1)
}

func (m *MockMissionUC) ListStaticMissions(ctx context.Context, userID string) ([]domain.Mission, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Mission), args.Error(1)
}

func (m *MockMissionUC) GetMissionDetail(ctx context.Context, userID, userMissionID string) (*usecase.MissionDetailResponse, error) {
	args := m.Called(ctx, userID, userMissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.MissionDetailResponse), args.Error(1)
}

func (m *MockMissionUC) ToggleTask(ctx context.Context, userID, userTaskID string, completed bool) error {
	return m.Called(ctx, userID, userTaskID, completed).Error(0)
}

// MockCompleteMissionUC
type MockCompleteMissionUC struct {
	mock.Mock
}

func (m *MockCompleteMissionUC) SubmitProof(ctx context.Context, userID, userMissionID string, file io.Reader, contentType string) error {
	return m.Called(ctx, userID, userMissionID, mock.Anything, contentType).Error(0)
}

// MockNotificationUC
type MockNotificationUC struct {
	mock.Mock
}

func (m *MockNotificationUC) ListNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
}

func (m *MockNotificationUC) MarkRead(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockNotificationUC) MarkAllRead(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *MockNotificationUC) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// Tests
