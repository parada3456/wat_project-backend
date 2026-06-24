package userusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	userusecase "github.com/j1hub/backend/internal/user/usecase"

	gamificationdomain "github.com/j1hub/backend/internal/gamification/domain"

	frienddomain "github.com/j1hub/backend/internal/friend/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/j1hub/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserUseCase_GetProfile_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	friendRepo := new(MockFriendshipRepository)
	hasher := new(MockHasher)
	uc := userusecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	userID := "usr_123"

	mockUser := &userdomain.User{UserID: userID, Email: "user@example.com"}
	mockProfile := &userdomain.Profile{UserID: userID, Bio: "Test bio"}
	mockCredit := &gamificationdomain.CreditScore{UserID: userID, CurrentScore: 100}

	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	profileRepo.On("FindByUserID", ctx, userID).Return(mockProfile, nil)
	creditRepo.On("FindByUserID", ctx, userID).Return(mockCredit, nil)
	userRepo.On("FindUserJob", ctx, userID).Return((*userdomain.UserJob)(nil), domain.ErrNotFound)
	userRepo.On("FindUserJobs", ctx, userID).Return(([]userdomain.UserJob)(nil), nil)

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
	uc := userusecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

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
	uc := userusecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	userID := "usr_123"
	cmd := userusecase.UpdateProfileCommand{
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
	uc := userusecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	user1 := &userdomain.User{UserID: "usr_1", FirstName: "John"}
	profile1 := &userdomain.Profile{UserID: "usr_1", RadarVisibility: userdomain.VisibilityShowFriends}

	userRepo.On("FindByID", ctx, "usr_1").Return(user1, nil)
	profileRepo.On("FindByUserID", ctx, "usr_1").Return(profile1, nil)

	// Friend mock: they are friends
	friendship := &frienddomain.Friendship{
		UserID1: "usr_1",
		UserID2: "usr_2",
		Status:  frienddomain.FriendshipAccepted,
	}
	friendRepo.On("FindByCanonicalPair", ctx, "usr_1", "usr_2").Return(friendship, nil)

	u, p, err := uc.GetPublicProfile(ctx, "usr_2", "usr_1")
	assert.NoError(t, err)
	assert.Equal(t, user1, u)
	assert.Equal(t, profile1, p)
}

func TestUserUseCase_AssignJob_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	friendRepo := new(MockFriendshipRepository)
	hasher := new(MockHasher)
	uc := userusecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	userID := "usr_123"
	jobID := "job_456"

	userRepo.On("AssignJob", ctx, userID, jobID, true, (*time.Time)(nil), (*time.Time)(nil)).Return(nil)

	err := uc.AssignJob(ctx, userID, jobID, true, nil, nil)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestUserUseCase_UpdatePassword_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	friendRepo := new(MockFriendshipRepository)
	hasher := new(MockHasher)
	uc := userusecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	userID := "usr_123"
	currentPassword := "old_pass"
	newPassword := "new_pass"
	hashedOldPass := "hashed_old_pass"
	hashedNewPass := "hashed_new_pass"

	mockUser := &userdomain.User{
		UserID:       userID,
		PasswordHash: hashedOldPass,
	}

	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	hasher.On("Verify", currentPassword, hashedOldPass).Return(true)
	hasher.On("Hash", newPassword).Return(hashedNewPass, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*userdomain.User")).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*userdomain.User)
		assert.Equal(t, hashedNewPass, u.PasswordHash)
	})

	err := uc.UpdatePassword(ctx, userID, currentPassword, newPassword)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
}

func TestUserUseCase_UpdatePassword_UserNotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	friendRepo := new(MockFriendshipRepository)
	hasher := new(MockHasher)
	uc := userusecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	userID := "usr_nonexistent"

	userRepo.On("FindByID", ctx, userID).Return((*userdomain.User)(nil), domain.ErrNotFound)

	err := uc.UpdatePassword(ctx, userID, "current", "new")

	assert.ErrorIs(t, err, domain.ErrNotFound)
	userRepo.AssertExpectations(t)
}

func TestUserUseCase_UpdatePassword_InvalidCurrentPassword(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	friendRepo := new(MockFriendshipRepository)
	hasher := new(MockHasher)
	uc := userusecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	userID := "usr_123"
	currentPassword := "wrong_pass"
	hashedOldPass := "hashed_old_pass"

	mockUser := &userdomain.User{
		UserID:       userID,
		PasswordHash: hashedOldPass,
	}

	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	hasher.On("Verify", currentPassword, hashedOldPass).Return(false)

	err := uc.UpdatePassword(ctx, userID, currentPassword, "new")

	assert.ErrorIs(t, err, domain.ErrUnauthorized)
	userRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
}

func TestUserUseCase_UpdatePassword_HashError(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)
	friendRepo := new(MockFriendshipRepository)
	hasher := new(MockHasher)
	uc := userusecase.NewUserUseCase(userRepo, profileRepo, creditRepo, friendRepo, hasher)

	ctx := context.Background()
	userID := "usr_123"
	currentPassword := "old_pass"
	newPassword := "new_pass"
	hashedOldPass := "hashed_old_pass"
	hashErr := errors.New("hashing error")

	mockUser := &userdomain.User{
		UserID:       userID,
		PasswordHash: hashedOldPass,
	}

	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	hasher.On("Verify", currentPassword, hashedOldPass).Return(true)
	hasher.On("Hash", newPassword).Return("", hashErr)

	err := uc.UpdatePassword(ctx, userID, currentPassword, newPassword)

	assert.ErrorIs(t, err, hashErr)
	userRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
}

