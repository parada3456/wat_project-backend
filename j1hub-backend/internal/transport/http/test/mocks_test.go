package test

import (
	"context"
	"io"

	expensedomain "github.com/j1hub/backend/internal/expense/domain"
	expenseusecase "github.com/j1hub/backend/internal/expense/usecase"
	frienddomain "github.com/j1hub/backend/internal/friend/domain"
	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
	gamificationusecase "github.com/j1hub/backend/internal/gamification/usecase"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	missionusecase "github.com/j1hub/backend/internal/mission/usecase"
	notificationdomain "github.com/j1hub/backend/internal/notification/domain"
	"github.com/stretchr/testify/mock"
)

// MockJourneyUC
type MockJourneyUC struct {
	mock.Mock
}

func (m *MockJourneyUC) ListPhases(ctx context.Context) ([]missiondomain.JourneyPhase, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.JourneyPhase), args.Error(1)
}

func (m *MockJourneyUC) GetHistory(ctx context.Context, userID string) ([]missiondomain.UserPhaseHistory, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.UserPhaseHistory), args.Error(1)
}

func (m *MockJourneyUC) ListUserBadges(ctx context.Context, userID string) ([]gamificationdomain.UserBadge, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]gamificationdomain.UserBadge), args.Error(1)
}

func (m *MockJourneyUC) GetCreditScoreHistory(ctx context.Context, userID string) ([]gamificationdomain.PointLedger, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]gamificationdomain.PointLedger), args.Error(1)
}

func (m *MockJourneyUC) GetPointsLedger(ctx context.Context, userID string) ([]gamificationdomain.PointLedger, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]gamificationdomain.PointLedger), args.Error(1)
}

// MockAdvancePhaseUC
type MockAdvancePhaseUC struct {
	mock.Mock
}

func (m *MockAdvancePhaseUC) TryAdvancePhase(ctx context.Context, userID string) (*gamificationusecase.PhaseTransitionResponse, error) {
	args := m.Called(ctx, userID)
	var resp *gamificationusecase.PhaseTransitionResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(*gamificationusecase.PhaseTransitionResponse)
	}
	return resp, args.Error(1)
}

// MockLeaderboardUC
type MockLeaderboardUC struct {
	mock.Mock
}

func (m *MockLeaderboardUC) GetLeaderboard(ctx context.Context, scope, jobID string) ([]gamificationusecase.LeaderboardEntry, error) {
	args := m.Called(ctx, scope, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]gamificationusecase.LeaderboardEntry), args.Error(1)
}

// MockManageExpenseUC
type MockManageExpenseUC struct {
	mock.Mock
}

func (m *MockManageExpenseUC) ListExpenses(ctx context.Context, userID string, page, pageSize int) ([]expensedomain.ExpenseTransaction, int, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]expensedomain.ExpenseTransaction), args.Int(1), args.Error(2)
}

func (m *MockManageExpenseUC) CreateExpense(ctx context.Context, userID string, cmd expenseusecase.CreateExpenseCmd) error {
	return m.Called(ctx, userID, cmd).Error(0)
}

func (m *MockManageExpenseUC) GetExpenseDetail(ctx context.Context, userID, id string) (*expensedomain.ExpenseTransaction, []expensedomain.ExpenseSplit, error) {
	args := m.Called(ctx, userID, id)
	var t *expensedomain.ExpenseTransaction
	if args.Get(0) != nil {
		t = args.Get(0).(*expensedomain.ExpenseTransaction)
	}
	var s []expensedomain.ExpenseSplit
	if args.Get(1) != nil {
		s = args.Get(1).([]expensedomain.ExpenseSplit)
	}
	return t, s, args.Error(2)
}

func (m *MockManageExpenseUC) DeleteExpense(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func (m *MockManageExpenseUC) ListPendingExpenses(ctx context.Context, userID string, page, pageSize int) ([]expensedomain.ExpenseSplit, int, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]expensedomain.ExpenseSplit), args.Int(1), args.Error(2)
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

func (m *MockFriendshipUC) ListPendingRequests(ctx context.Context, userID string, page, pageSize int) ([]frienddomain.Friendship, int, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]frienddomain.Friendship), args.Int(1), args.Error(2)
}

func (m *MockFriendshipUC) RespondToRequest(ctx context.Context, userID, friendshipID string, accept bool) error {
	return m.Called(ctx, userID, friendshipID, accept).Error(0)
}

func (m *MockFriendshipUC) ListFriends(ctx context.Context, userID string, page, pageSize int) ([]frienddomain.Friendship, int, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]frienddomain.Friendship), args.Int(1), args.Error(2)
}

func (m *MockFriendshipUC) RemoveFriend(ctx context.Context, userID, friendID string) error {
	return m.Called(ctx, userID, friendID).Error(0)
}

// MockRadarUC
type MockRadarUC struct {
	mock.Mock
}

func (m *MockRadarUC) GetRadar(ctx context.Context, requesterID string) ([]gamificationusecase.RadarEntry, error) {
	args := m.Called(ctx, requesterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]gamificationusecase.RadarEntry), args.Error(1)
}

// MockJobUC
type MockJobUC struct {
	mock.Mock
}

func (m *MockJobUC) ListJobs(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]jobdomain.JobPosting, int, error) {
	args := m.Called(ctx, filters, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]jobdomain.JobPosting), args.Int(1), args.Error(2)
}

func (m *MockJobUC) GetJobDetail(ctx context.Context, id string) (*jobdomain.JobPosting, []jobdomain.JobHousing, *jobdomain.JobOverallRating, error) {
	args := m.Called(ctx, id)
	var j *jobdomain.JobPosting
	if args.Get(0) != nil {
		j = args.Get(0).(*jobdomain.JobPosting)
	}
	var h []jobdomain.JobHousing
	if args.Get(1) != nil {
		h = args.Get(1).([]jobdomain.JobHousing)
	}
	var r *jobdomain.JobOverallRating
	if args.Get(2) != nil {
		r = args.Get(2).(*jobdomain.JobOverallRating)
	}
	return j, h, r, args.Error(3)
}

func (m *MockJobUC) AddToCart(ctx context.Context, userID, jobID string) error {
	return m.Called(ctx, userID, jobID).Error(0)
}

func (m *MockJobUC) ListCart(ctx context.Context, userID string) ([]jobdomain.UserCart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jobdomain.UserCart), args.Error(1)
}

func (m *MockJobUC) RemoveFromCart(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func (m *MockJobUC) WriteReview(ctx context.Context, userID, jobID string, rev *jobdomain.JobReview) error {
	return m.Called(ctx, userID, jobID, rev).Error(0)
}

func (m *MockJobUC) ListReviews(ctx context.Context, jobID string) ([]jobdomain.JobReview, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]jobdomain.JobReview), args.Error(1)
}

func (m *MockJobUC) UpdateCartStatus(ctx context.Context, userID, cartID string, status jobdomain.CartStatus) error {
	return m.Called(ctx, userID, cartID, status).Error(0)
}

func (m *MockJobUC) CreateJob(ctx context.Context, job *jobdomain.JobPosting) error {
	return m.Called(ctx, job).Error(0)
}

func (m *MockJobUC) UpdateJob(ctx context.Context, job *jobdomain.JobPosting) error {
	return m.Called(ctx, job).Error(0)
}

func (m *MockJobUC) PatchJob(ctx context.Context, id string, updates map[string]interface{}) error {
	return m.Called(ctx, id, updates).Error(0)
}

func (m *MockJobUC) DeleteJob(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// MockMissionUC
type MockMissionUC struct {
	mock.Mock
}

func (m *MockMissionUC) ListAvailableMissions(ctx context.Context, userID string, ids []string) ([]missiondomain.UserMission, error) {
	args := m.Called(ctx, userID, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.UserMission), args.Error(1)
}

func (m *MockMissionUC) ListStaticMissions(ctx context.Context, userID string, ids []string) ([]missiondomain.Mission, error) {
	args := m.Called(ctx, userID, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.Mission), args.Error(1)
}

func (m *MockMissionUC) GetMissionDetail(ctx context.Context, userID, userMissionID string) (*missionusecase.MissionDetailResponse, error) {
	args := m.Called(ctx, userID, userMissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*missionusecase.MissionDetailResponse), args.Error(1)
}

func (m *MockMissionUC) ToggleTask(ctx context.Context, userID, userTaskID string, completed bool) error {
	return m.Called(ctx, userID, userTaskID, completed).Error(0)
}

func (m *MockMissionUC) ListTasks(ctx context.Context, ids []string) ([]missiondomain.Task, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.Task), args.Error(1)
}

func (m *MockMissionUC) ListUserTasks(ctx context.Context, ids []string) ([]missiondomain.UserTask, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.UserTask), args.Error(1)
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

func (m *MockNotificationUC) ListNotifications(ctx context.Context, userID string, isRead *bool, page, pageSize int) ([]notificationdomain.Notification, int, error) {
	args := m.Called(ctx, userID, isRead, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]notificationdomain.Notification), args.Int(1), args.Error(2)
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
