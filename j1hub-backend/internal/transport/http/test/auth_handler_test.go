package test

import (
	authusecase "github.com/j1hub/backend/internal/auth/usecase"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/j1hub/backend/internal/domain"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	authhandler "github.com/j1hub/backend/internal/auth/adapter/http"
	port "github.com/j1hub/backend/internal/auth/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRegisterUserUC
type MockRegisterUserUC struct {
	mock.Mock
}

func (m *MockRegisterUserUC) Register(ctx context.Context, cmd authusecase.RegisterCommand) (*userdomain.User, *port.TokenPair, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*userdomain.User), args.Get(1).(*port.TokenPair), args.Error(2)
}

// MockLoginUC
type MockLoginUC struct {
	mock.Mock
}

func (m *MockLoginUC) Login(ctx context.Context, cmd authusecase.LoginCommand) (*userdomain.User, *port.TokenPair, error) {
	args := m.Called(ctx, cmd)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*userdomain.User), args.Get(1).(*port.TokenPair), args.Error(2)
}

func (m *MockLoginUC) Refresh(ctx context.Context, refreshToken string) (*port.TokenPair, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.TokenPair), args.Error(1)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	regUC := new(MockRegisterUserUC)
	logUC := new(MockLoginUC)
	h := authhandler.NewAuthHandler(regUC, logUC)

	user := &userdomain.User{
		UserID:    "usr_123",
		Email:     "john@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	tokens := &port.TokenPair{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	regUC.On("Register", mock.Anything, authusecase.RegisterCommand{
		Email:     "john@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}).Return(user, tokens, nil)

	body := `{"email":"john@example.com","password":"password123","first_name":"John","last_name":"Doe"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "usr_123", resp["user_id"])
	assert.Equal(t, "access", resp["access_token"])
	assert.Equal(t, "refresh", resp["refresh_token"])
}

func TestAuthHandler_Register_ValidationError(t *testing.T) {
	regUC := new(MockRegisterUserUC)
	logUC := new(MockLoginUC)
	h := authhandler.NewAuthHandler(regUC, logUC)

	// Missing fields
	body := `{"email":"john"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	regUC := new(MockRegisterUserUC)
	logUC := new(MockLoginUC)
	h := authhandler.NewAuthHandler(regUC, logUC)

	user := &userdomain.User{
		UserID: "usr_123",
		Email:  "john@example.com",
	}
	tokens := &port.TokenPair{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}

	logUC.On("Login", mock.Anything, authusecase.LoginCommand{
		Email:    "john@example.com",
		Password: "password123",
	}).Return(user, tokens, nil)

	body := `{"email":"john@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_Login_Failure(t *testing.T) {
	regUC := new(MockRegisterUserUC)
	logUC := new(MockLoginUC)
	h := authhandler.NewAuthHandler(regUC, logUC)

	logUC.On("Login", mock.Anything, authusecase.LoginCommand{
		Email:    "john@example.com",
		Password: "wrong_password",
	}).Return((*userdomain.User)(nil), (*port.TokenPair)(nil), domain.ErrUnauthorized)

	body := `{"email":"john@example.com","password":"wrong_password"}`
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
