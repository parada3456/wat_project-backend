package friendusecase_test

import (
	"context"
	"testing"
	"time"

	friendusecase "github.com/j1hub/backend/internal/friend/usecase"

	frienddomain "github.com/j1hub/backend/internal/friend/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/j1hub/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestManageFriendshipUseCase_SendRequest_Success(t *testing.T) {
	friendRepo := new(MockFriendshipRepository)
	userRepo := new(MockUserRepository)
	notifier := new(MockNotifierPort)

	nowTime := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	clock := &MockClock{NowTime: nowTime}

	uc := friendusecase.NewManageFriendshipUseCase(friendRepo, userRepo, notifier, clock)

	ctx := context.Background()
	senderID := "usr_aaa"
	targetID := "usr_bbb" // aaa < bbb, so u1 = aaa, u2 = bbb

	userRepo.On("FindByID", ctx, targetID).Return(&userdomain.User{UserID: targetID}, nil)
	friendRepo.On("FindByCanonicalPair", ctx, senderID, targetID).Return((*frienddomain.Friendship)(nil), nil)
	friendRepo.On("Insert", ctx, mock.AnythingOfType("*frienddomain.Friendship")).Return(nil).Run(func(args mock.Arguments) {
		f := args.Get(1).(*frienddomain.Friendship)
		assert.Equal(t, senderID, f.UserID1)
		assert.Equal(t, targetID, f.UserID2)
		assert.Equal(t, frienddomain.FriendshipPending, f.Status)
	})
	notifier.On("Send", ctx, targetID, "Friend request", "Someone wants to be your friend!").Return(nil)

	err := uc.SendRequest(ctx, senderID, targetID)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	friendRepo.AssertExpectations(t)
}

func TestManageFriendshipUseCase_SendRequest_Duplicate(t *testing.T) {
	friendRepo := new(MockFriendshipRepository)
	userRepo := new(MockUserRepository)
	notifier := new(MockNotifierPort)
	clock := &MockClock{}

	uc := friendusecase.NewManageFriendshipUseCase(friendRepo, userRepo, notifier, clock)

	ctx := context.Background()
	senderID := "usr_aaa"
	targetID := "usr_bbb"

	userRepo.On("FindByID", ctx, targetID).Return(&userdomain.User{UserID: targetID}, nil)
	friendRepo.On("FindByCanonicalPair", ctx, senderID, targetID).Return(&frienddomain.Friendship{FriendshipID: "frn_123"}, nil)

	err := uc.SendRequest(ctx, senderID, targetID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrDuplicateFriend, err)
}

func TestManageFriendshipUseCase_RespondToRequest_Accept(t *testing.T) {
	friendRepo := new(MockFriendshipRepository)
	userRepo := new(MockUserRepository)
	notifier := new(MockNotifierPort)
	clock := &MockClock{}

	uc := friendusecase.NewManageFriendshipUseCase(friendRepo, userRepo, notifier, clock)

	ctx := context.Background()
	responderID := "usr_bbb"
	friendshipID := "frn_123"

	mockFriendship := &frienddomain.Friendship{
		FriendshipID: friendshipID,
		UserID1:      "usr_aaa",
		UserID2:      "usr_bbb",
		Status:       frienddomain.FriendshipPending,
	}

	friendRepo.On("FindByID", ctx, friendshipID).Return(mockFriendship, nil)
	friendRepo.On("UpdateStatus", ctx, friendshipID, frienddomain.FriendshipAccepted).Return(nil)
	notifier.On("Send", ctx, "usr_aaa", "Friend request accepted", "You are now friends!").Return(nil)

	err := uc.RespondToRequest(ctx, responderID, friendshipID, true)

	assert.NoError(t, err)
}

func TestManageFriendshipUseCase_RespondToRequest_Forbidden(t *testing.T) {
	friendRepo := new(MockFriendshipRepository)
	userRepo := new(MockUserRepository)
	notifier := new(MockNotifierPort)
	clock := &MockClock{}

	uc := friendusecase.NewManageFriendshipUseCase(friendRepo, userRepo, notifier, clock)

	ctx := context.Background()
	responderID := "usr_ccc" // not aaa or bbb
	friendshipID := "frn_123"

	mockFriendship := &frienddomain.Friendship{
		FriendshipID: friendshipID,
		UserID1:      "usr_aaa",
		UserID2:      "usr_bbb",
	}

	friendRepo.On("FindByID", ctx, friendshipID).Return(mockFriendship, nil)

	err := uc.RespondToRequest(ctx, responderID, friendshipID, true)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrForbidden, err)
}

func TestManageFriendshipUseCase_ListFriends_Success(t *testing.T) {
	friendRepo := new(MockFriendshipRepository)
	userRepo := new(MockUserRepository)
	notifier := new(MockNotifierPort)
	clock := &MockClock{}

	uc := friendusecase.NewManageFriendshipUseCase(friendRepo, userRepo, notifier, clock)

	ctx := context.Background()
	userID := "usr_aaa"
	mockFriends := []frienddomain.Friendship{
		{FriendshipID: "frn_1", UserID1: "usr_aaa", UserID2: "usr_bbb", Status: frienddomain.FriendshipAccepted},
	}

	friendRepo.On("FindFriendsOf", ctx, userID).Return(mockFriends, nil)

	res, err := uc.ListFriends(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, mockFriends, res)
}

func TestManageFriendshipUseCase_ListPendingRequests_Stub(t *testing.T) {
	uc := friendusecase.NewManageFriendshipUseCase(nil, nil, nil, &MockClock{})
	res, err := uc.ListPendingRequests(context.Background(), "usr_aaa")
	assert.Nil(t, res)
	assert.NoError(t, err)
}

func TestManageFriendshipUseCase_RemoveFriend_Stub(t *testing.T) {
	uc := friendusecase.NewManageFriendshipUseCase(nil, nil, nil, &MockClock{})
	err := uc.RemoveFriend(context.Background(), "usr_aaa", "frn_1")
	assert.NoError(t, err)
}
