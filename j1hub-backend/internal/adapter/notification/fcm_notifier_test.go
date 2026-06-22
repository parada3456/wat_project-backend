package notification

import (
	"context"
	"errors"
	"testing"

	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/j1hub/backend/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, u *userdomain.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *MockUserRepo) FindByID(ctx context.Context, id string) (*userdomain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userdomain.User), args.Error(1)
}
func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*userdomain.User), args.Error(1)
}
func (m *MockUserRepo) Update(ctx context.Context, u *userdomain.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *MockUserRepo) IncrementPoints(ctx context.Context, userID string, lifetimeDelta, phaseDelta int) error {
	return m.Called(ctx, userID, lifetimeDelta, phaseDelta).Error(0)
}
func (m *MockUserRepo) ResetStreak(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}
func (m *MockUserRepo) SetPhase(ctx context.Context, userID, phaseID string) error {
	return m.Called(ctx, userID, phaseID).Error(0)
}
func (m *MockUserRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func TestNewFCMNotifier(t *testing.T) {
	cfg := &config.Config{
		FCMCredentialsPath: "invalid.json",
	}
	userRepo := new(MockUserRepo)
	notifier := NewFCMNotifier(cfg, userRepo)
	assert.NotNil(t, notifier)

	err := notifier.Send(context.Background(), "u1", "t", "b")
	assert.NoError(t, err)
}

func TestFCMNotifier_Send(t *testing.T) {
	userRepo := new(MockUserRepo)
	notifier := &fcmNotifier{
		client:   nil,
		userRepo: userRepo,
	}

	// User not found / error
	userRepo.On("FindByID", mock.Anything, "usr_err").Return(nil, errors.New("database error")).Once()
	err := notifier.Send(context.Background(), "usr_err", "Title", "Body")
	assert.Error(t, err)

	// User found, but ID empty
	userRepo.On("FindByID", mock.Anything, "usr_empty").Return(&userdomain.User{UserID: ""}, nil).Once()
	err = notifier.Send(context.Background(), "usr_empty", "Title", "Body")
	assert.NoError(t, err)

	// User found successfully
	userRepo.On("FindByID", mock.Anything, "usr_1").Return(&userdomain.User{UserID: "usr_1"}, nil).Once()
	err = notifier.Send(context.Background(), "usr_1", "Title", "Body")
	assert.NoError(t, err)
}
