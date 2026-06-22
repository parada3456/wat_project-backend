package userusecase_test

import (
	"context"
	"errors"
	"testing"

	frienddomain "github.com/j1hub/backend/internal/friend/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserUseCase_GetProfile_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	friendRepo := new(MockFriendshipRepository)
	hasher := new(MockHasher)
	uc := usecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	userID := "usr_123"

	mockUser := &userdomain.User{UserID: userID, Email: "user@example.com"}
	mockProfile := &userdomain.Profile{UserID: userID, Bio: "Test bio"}
	mockCredit := &userdomain.CreditScore{UserID: userID, CurrentScore: 100}

	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	profileRepo.On("FindByUserID", ctx, userID).Return(mockProfile, nil)
	creditRepo.On("FindByUserID", ctx, userID).Return(mockCredit, nil)

	resp, err := uc.GetProfile(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, mockUser, resp.User)
	assert.Equal(t, mockProfile, resp.Profile)
	assert.Equal(t, mockCredit, resp.CreditScore)

	userRepo.AssertExpectations(t)
	profileRepo.AssertExpectations(t)
	creditRepo.AssertExpectations(t)
}

func TestUserUseCase_GetProfile_UserNotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	friendRepo := new(MockFriendshipRepository)
	hasher := new(MockHasher)
	uc := usecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	userID := "usr_nonexistent"

	userRepo.On("FindByID", ctx, userID).Return((*userdomain.User)(nil), domain.ErrNotFound)

	resp, err := uc.GetProfile(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestUserUseCase_UpdateProfile_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	friendRepo := new(MockFriendshipRepository)
	hasher := new(MockHasher)
	uc := usecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	userID := "usr_123"
	cmd := usecase.UpdateProfileCommand{
		FirstName: "John",
		LastName:  "Doe",
		Bio:       "New bio",
		AvatarURL: "https://newavatar.com",
	}

	mockUser := &userdomain.User{UserID: userID, FirstName: "Old", LastName: "Name"}
	mockProfile := &userdomain.Profile{ProfileID: "prf_123", UserID: userID, Bio: "Old bio"}

	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*userdomain.User")).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*userdomain.User)
		assert.Equal(t, "John", u.FirstName)
		assert.Equal(t, "Doe", u.LastName)
	})

	profileRepo.On("FindByUserID", ctx, userID).Return(mockProfile, nil)
	profileRepo.On("Update", ctx, mock.AnythingOfType("*userdomain.Profile")).Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(1).(*userdomain.Profile)
		assert.Equal(t, "New bio", p.Bio)
		assert.Equal(t, "https://newavatar.com", p.AvatarURL)
	})

	err := uc.UpdateProfile(ctx, userID, cmd)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	profileRepo.AssertExpectations(t)
}

func TestUserUseCase_GetPublicProfile_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	friendRepo := new(MockFriendshipRepository)
	hasher := new(MockHasher)
	uc := usecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	user1 := &userdomain.User{UserID: "usr_1", FirstName: "John"}
	profile1 := &userdomain.Profile{UserID: "usr_1", RadarVisibility: domain.VisibilityShowFriends}

	userRepo.On("FindByID", ctx, "usr_1").Return(user1, nil)
	profileRepo.On("FindByUserID", ctx, "usr_1").Return(profile1, nil)

	// Friend mock: they are friends
	friendship := &frienddomain.Friendship{
		UserID1: "usr_1",
		UserID2: "usr_2",
		Status:  domain.FriendshipAccepted,
	}
	friendRepo.On("FindByCanonicalPair", ctx, "usr_1", "usr_2").Return(friendship, nil)

	u, p, err := uc.GetPublicProfile(ctx, "usr_2", "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, user1, u)
	assert.Equal(t, profile1, p)
}
