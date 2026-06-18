package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserUseCase_GetProfile_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)

	uc := usecase.NewUserUseCase(userRepo, profileRepo, creditRepo)

	ctx := context.Background()
	userID := "usr_123"

	mockUser := &domain.User{UserID: userID, Email: "user@example.com"}
	mockProfile := &domain.Profile{UserID: userID, Bio: "Test bio"}
	mockCredit := &domain.CreditScore{UserID: userID, CurrentScore: 100}

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

	uc := usecase.NewUserUseCase(userRepo, profileRepo, creditRepo)

	ctx := context.Background()
	userID := "usr_nonexistent"

	userRepo.On("FindByID", ctx, userID).Return((*domain.User)(nil), domain.ErrNotFound)

	resp, err := uc.GetProfile(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestUserUseCase_UpdateProfile_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	creditRepo := new(MockCreditScoreRepository)

	uc := usecase.NewUserUseCase(userRepo, profileRepo, creditRepo)

	ctx := context.Background()
	userID := "usr_123"
	cmd := usecase.UpdateProfileCommand{
		FirstName: "John",
		LastName:  "Doe",
		Bio:       "New bio",
		AvatarURL: "https://newavatar.com",
	}

	mockUser := &domain.User{UserID: userID, FirstName: "Old", LastName: "Name"}
	mockProfile := &domain.Profile{ProfileID: "prf_123", UserID: userID, Bio: "Old bio"}

	userRepo.On("FindByID", ctx, userID).Return(mockUser, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*domain.User")).Return(nil).Run(func(args mock.Arguments) {
		u := args.Get(1).(*domain.User)
		assert.Equal(t, "John", u.FirstName)
		assert.Equal(t, "Doe", u.LastName)
	})

	profileRepo.On("FindByUserID", ctx, userID).Return(mockProfile, nil)
	profileRepo.On("Update", ctx, mock.AnythingOfType("*domain.Profile")).Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(1).(*domain.Profile)
		assert.Equal(t, "New bio", p.Bio)
		assert.Equal(t, "https://newavatar.com", p.AvatarURL)
	})

	err := uc.UpdateProfile(ctx, userID, cmd)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	profileRepo.AssertExpectations(t)
}
