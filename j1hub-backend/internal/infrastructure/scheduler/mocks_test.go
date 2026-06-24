package scheduler

import (
	"context"

	expensedomain "github.com/j1hub/backend/internal/expense/domain"
	expenseport "github.com/j1hub/backend/internal/expense/port"
	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"
	gamificationport "github.com/j1hub/backend/internal/gamification/port"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	missionport "github.com/j1hub/backend/internal/mission/port"
	notificationport "github.com/j1hub/backend/internal/notification/port"
	userport "github.com/j1hub/backend/internal/user/port"
	"github.com/stretchr/testify/mock"
)

type MockExpenseSplitRepository struct {
	mock.Mock
	expenseport.ExpenseSplitRepository
}

func (m *MockExpenseSplitRepository) FindOverdue(ctx context.Context) ([]expensedomain.ExpenseSplit, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]expensedomain.ExpenseSplit), args.Error(1)
}

func (m *MockExpenseSplitRepository) UpdatePaymentStatus(ctx context.Context, id string, status expensedomain.PaymentStatus, slipURL string) error {
	args := m.Called(ctx, id, status, slipURL)
	return args.Error(0)
}

type MockCreditScoreRepository struct {
	mock.Mock
	gamificationport.CreditScoreRepository
}

func (m *MockCreditScoreRepository) Decrement(ctx context.Context, userID string, amount int) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

type MockPointLedgerRepository struct {
	mock.Mock
	gamificationport.PointLedgerRepository
}

func (m *MockPointLedgerRepository) Insert(ctx context.Context, ledger *gamificationdomain.PointLedger) error {
	args := m.Called(ctx, ledger)
	return args.Error(0)
}

type MockNotifierPort struct {
	mock.Mock
	notificationport.NotifierPort
}

func (m *MockNotifierPort) Send(ctx context.Context, userID string, title string, body string) error {
	args := m.Called(ctx, userID, title, body)
	return args.Error(0)
}

type MockUserMissionRepository struct {
	mock.Mock
	missionport.UserMissionRepository
}

func (m *MockUserMissionRepository) FindOverdue(ctx context.Context) ([]missiondomain.UserMission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]missiondomain.UserMission), args.Error(1)
}

func (m *MockUserMissionRepository) UpdateStatus(ctx context.Context, id string, status missiondomain.UserMissionStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

type MockMissionRepository struct {
	mock.Mock
	missionport.MissionRepository
}

func (m *MockMissionRepository) FindByID(ctx context.Context, id string) (*missiondomain.Mission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*missiondomain.Mission), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
	userport.UserRepository
}

func (m *MockUserRepository) ResetStreak(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
