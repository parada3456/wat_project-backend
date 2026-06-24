package scheduler

import (
	"context"
	"testing"

	expensedomain "github.com/j1hub/backend/internal/expense/domain"
	missiondomain "github.com/j1hub/backend/internal/mission/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOverdueExpenseJob_Run_Success(t *testing.T) {
	splitRepo := new(MockExpenseSplitRepository)
	creditRepo := new(MockCreditScoreRepository)
	ledgerRepo := new(MockPointLedgerRepository)
	notifier := new(MockNotifierPort)

	job := NewOverdueExpenseJob(splitRepo, creditRepo, ledgerRepo, notifier)

	ctx := context.Background()
	mockOverdueSplits := []expensedomain.ExpenseSplit{
		{SplitID: "spl_1", UserID: "usr_123", OweAmount: 50.0, PayslipURL: "http://payslip.jpg"},
	}

	splitRepo.On("FindOverdue", ctx).Return(mockOverdueSplits, nil)
	splitRepo.On("UpdatePaymentStatus", ctx, "spl_1", expensedomain.PaymentOverdue, "http://payslip.jpg").Return(nil)
	creditRepo.On("Decrement", ctx, "usr_123", 10).Return(nil)
	ledgerRepo.On("Insert", ctx, mock.AnythingOfType("*gamificationdomain.PointLedger")).Return(nil)
	notifier.On("Send", ctx, "usr_123", "Overdue payment", "Your payment is overdue. Credit score -10.").Return(nil)

	err := job.Run(ctx)

	assert.NoError(t, err)
	splitRepo.AssertExpectations(t)
	creditRepo.AssertExpectations(t)
	ledgerRepo.AssertExpectations(t)
	notifier.AssertExpectations(t)
}

func TestOverdueMissionJob_Run_Success(t *testing.T) {
	umRepo := new(MockUserMissionRepository)
	missionRepo := new(MockMissionRepository)
	userRepo := new(MockUserRepository)
	notifier := new(MockNotifierPort)

	job := NewOverdueMissionJob(umRepo, missionRepo, userRepo, notifier)

	ctx := context.Background()
	mockUMs := []missiondomain.UserMission{
		{UserMissionID: "ums_1", UserID: "usr_123", MissionID: "m_1"},
	}

	umRepo.On("FindOverdue", ctx).Return(mockUMs, nil)
	umRepo.On("UpdateStatus", ctx, "ums_1", missiondomain.StatusOverdue).Return(nil)

	mockMission := &missiondomain.Mission{MissionID: "m_1", IsMandatory: true}
	missionRepo.On("FindByID", ctx, "m_1").Return(mockMission, nil)
	userRepo.On("ResetStreak", ctx, "usr_123").Return(nil)
	notifier.On("Send", ctx, "usr_123", "Mission overdue", "A mission is past its due date!").Return(nil)

	err := job.Run(ctx)

	assert.NoError(t, err)
	umRepo.AssertExpectations(t)
	missionRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	notifier.AssertExpectations(t)
}
