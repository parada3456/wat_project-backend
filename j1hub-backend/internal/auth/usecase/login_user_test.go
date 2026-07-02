package authusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	authusecase "github.com/parada3456/wat_project-backend/internal/auth/usecase"

	userdomain "github.com/parada3456/wat_project-backend/internal/user/domain"

	port "github.com/parada3456/wat_project-backend/internal/auth/port"
	"github.com/parada3456/wat_project-backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestLoginUseCase_Login_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	hasher := new(MockHasher)
	tokenIssuer := new(MockIssuer)

	uc := authusecase.NewLoginUseCase(userRepo, profileRepo, hasher, tokenIssuer)

	ctx := context.Background()
	email := "user@example.com"
	password := "password123"
	hash := "hashed_password"
	userID := "usr_123"

	mockUser := &userdomain.User{
		UserID:       userID,
		Email:        email,
		PasswordHash: hash,
	}

	mockProfile := &userdomain.Profile{
		ProfileID: "prf_123",
		UserID:    userID,
	}

	mockTokens := &port.TokenPair{
		AccessToken:  "access_token_jwt",
		RefreshToken: "refresh_token_jwt",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	userRepo.On("FindByEmail", ctx, email).Return(mockUser, nil)
	profileRepo.On("FindByUserID", ctx, userID).Return(mockProfile, nil)
	hasher.On("Verify", password, hash).Return(true)
	tokenIssuer.On("Issue", userID, false).Return(mockTokens, nil)

	user, profile, tokens, err := uc.Login(ctx, authusecase.LoginCommand{
		Email:    email,
		Password: password,
	})

	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)
	assert.Equal(t, mockProfile, profile)
	assert.Equal(t, mockTokens, tokens)

	userRepo.AssertExpectations(t)
	profileRepo.AssertExpectations(t)
	hasher.AssertExpectations(t)
	tokenIssuer.AssertExpectations(t)
}

func TestLoginUseCase_Login_IssueTokenError(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	hasher := new(MockHasher)
	tokenIssuer := new(MockIssuer)

	uc := authusecase.NewLoginUseCase(userRepo, profileRepo, hasher, tokenIssuer)

	ctx := context.Background()
	email := "user@example.com"
	password := "password123"
	hash := "hashed_password"
	userID := "usr_123"

	mockUser := &userdomain.User{
		UserID:       userID,
		Email:        email,
		PasswordHash: hash,
	}

	userRepo.On("FindByEmail", ctx, email).Return(mockUser, nil)
	hasher.On("Verify", password, hash).Return(true)
	tokenIssuer.On("Issue", userID, false).Return((*port.TokenPair)(nil), errors.New("issuance error"))

	user, profile, tokens, err := uc.Login(ctx, authusecase.LoginCommand{
		Email:    email,
		Password: password,
	})

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, profile)
	assert.Nil(t, tokens)
	assert.Equal(t, "issuance error", err.Error())
}

func TestLoginUseCase_Login_UserNotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	hasher := new(MockHasher)
	tokenIssuer := new(MockIssuer)

	uc := authusecase.NewLoginUseCase(userRepo, profileRepo, hasher, tokenIssuer)

	ctx := context.Background()
	email := "nonexistent@example.com"

	userRepo.On("FindByEmail", ctx, email).Return((*userdomain.User)(nil), domain.ErrNotFound)

	user, profile, tokens, err := uc.Login(ctx, authusecase.LoginCommand{
		Email:    email,
		Password: "somepassword",
	})

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, profile)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "Invalid credentials")
}

func TestLoginUseCase_Login_WrongPassword(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	hasher := new(MockHasher)
	tokenIssuer := new(MockIssuer)

	uc := authusecase.NewLoginUseCase(userRepo, profileRepo, hasher, tokenIssuer)

	ctx := context.Background()
	email := "user@example.com"
	password := "wrongpassword"
	hash := "hashed_password"

	mockUser := &userdomain.User{
		UserID:       "usr_123",
		Email:        email,
		PasswordHash: hash,
	}

	userRepo.On("FindByEmail", ctx, email).Return(mockUser, nil)
	hasher.On("Verify", password, hash).Return(false)

	user, profile, tokens, err := uc.Login(ctx, authusecase.LoginCommand{
		Email:    email,
		Password: password,
	})

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, profile)
	assert.Nil(t, tokens)
	assert.Contains(t, err.Error(), "Invalid credentials")
}

func TestLoginUseCase_Refresh_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	hasher := new(MockHasher)
	tokenIssuer := new(MockIssuer)

	uc := authusecase.NewLoginUseCase(userRepo, profileRepo, hasher, tokenIssuer)

	ctx := context.Background()
	refreshToken := "valid_refresh_token"
	mockTokens := &port.TokenPair{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	}

	tokenIssuer.On("Refresh", refreshToken).Return(mockTokens, nil)

	tokens, err := uc.Refresh(ctx, refreshToken)

	assert.NoError(t, err)
	assert.Equal(t, mockTokens, tokens)
}

func TestLoginUseCase_Refresh_Error(t *testing.T) {
	userRepo := new(MockUserRepository)
	profileRepo := new(MockProfileRepository)
	hasher := new(MockHasher)
	tokenIssuer := new(MockIssuer)

	uc := authusecase.NewLoginUseCase(userRepo, profileRepo, hasher, tokenIssuer)

	ctx := context.Background()
	refreshToken := "invalid_refresh_token"

	tokenIssuer.On("Refresh", refreshToken).Return((*port.TokenPair)(nil), errors.New("invalid refresh token"))

	tokens, err := uc.Refresh(ctx, refreshToken)

	assert.Error(t, err)
	assert.Nil(t, tokens)
}
