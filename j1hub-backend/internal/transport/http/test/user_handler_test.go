package test

import (
	userusecase "github.com/j1hub/backend/internal/user/usecase"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/j1hub/backend/internal/user/adapter/http/dto"
	userdomain "github.com/j1hub/backend/internal/user/domain"

	"github.com/go-chi/chi/v5"
	userhandler "github.com/j1hub/backend/internal/user/adapter/http"
	"github.com/j1hub/backend/internal/transport/http/middleware"
	port "github.com/j1hub/backend/internal/auth/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUC
type MockUserUC struct {
	mock.Mock
}

func (m *MockUserUC) GetProfile(ctx context.Context, userID string) (*userusecase.UserProfileResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userusecase.UserProfileResponse), args.Error(1)
}

func (m *MockUserUC) GetPublicProfile(ctx context.Context, currentUserID, targetUserID string) (*userdomain.User, *userdomain.Profile, error) {
	args := m.Called(ctx, currentUserID, targetUserID)
	var u *userdomain.User
	if args.Get(0) != nil {
		u = args.Get(0).(*userdomain.User)
	}
	var p *userdomain.Profile
	if args.Get(1) != nil {
		p = args.Get(1).(*userdomain.Profile)
	}
	return u, p, args.Error(2)
}

func (m *MockUserUC) UpdateProfile(ctx context.Context, userID string, cmd userusecase.UpdateProfileCommand) error {
	return m.Called(ctx, userID, cmd).Error(0)
}

func (m *MockUserUC) UpdateLocation(ctx context.Context, userID string, lat, lng float64) error {
	return m.Called(ctx, userID, lat, lng).Error(0)
}

func (m *MockUserUC) UpdateSettings(ctx context.Context, userID string, settings map[string]interface{}) error {
	return m.Called(ctx, userID, settings).Error(0)
}

func (m *MockUserUC) DeleteAccount(ctx context.Context, userID string, password string) error {
	return m.Called(ctx, userID, password).Error(0)
}

func (m *MockUserUC) AssignJob(ctx context.Context, userID, jobID string, isMain bool, startDate, endDate *time.Time) error {
	return m.Called(ctx, userID, jobID, isMain, startDate, endDate).Error(0)
}

func (m *MockUserUC) UpdatePassword(ctx context.Context, userID string, currentPassword, newPassword string) error {
	return m.Called(ctx, userID, currentPassword, newPassword).Error(0)
}

func TestUserHandler_GetProfile_Success(t *testing.T) {
	userUC := new(MockUserUC)
	h := userhandler.NewUserHandler(userUC)

	userID := "usr_123"
	profileResp := &userusecase.UserProfileResponse{
		User:    &userdomain.User{UserID: userID, Email: "john@example.com"},
		Profile: &userdomain.Profile{UserID: userID, Bio: "Hello world"},
	}

	userUC.On("GetProfile", mock.Anything, userID).Return(profileResp, nil)

	req := httptest.NewRequest("GET", "/api/v1/user/profile", nil)
	ctx := middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: userID})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.GetProfile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.GetProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "john@example.com", resp.Email)
	assert.Equal(t, "Hello world", resp.Bio)
}

func TestUserHandler_GetProfile_Unauthorized(t *testing.T) {
	userUC := new(MockUserUC)
	h := userhandler.NewUserHandler(userUC)

	// No claims in context
	req := httptest.NewRequest("GET", "/api/v1/user/profile", nil)
	w := httptest.NewRecorder()

	h.GetProfile(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_UpdateProfile_Success(t *testing.T) {
	userUC := new(MockUserUC)
	h := userhandler.NewUserHandler(userUC)

	userID := "usr_123"
	cmd := userusecase.UpdateProfileCommand{
		FirstName: "John",
		LastName:  "Doe",
		Bio:       "New bio",
		AvatarURL: "https://avatar.com",
	}

	userUC.On("UpdateProfile", mock.Anything, userID, cmd).Return(nil)

	body := `{"first_name":"John","last_name":"Doe","bio":"New bio","avatar_url":"https://avatar.com"}`
	req := httptest.NewRequest("PATCH", "/api/v1/user/profile", bytes.NewBufferString(body))
	ctx := middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: userID})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	h.UpdateProfile(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestUserHandler_GetPublicProfile_Success(t *testing.T) {
	userUC := new(MockUserUC)
	h := userhandler.NewUserHandler(userUC)

	userUC.On("GetPublicProfile", mock.Anything, "usr_1", "usr_2").Return(
		&userdomain.User{UserID: "usr_2", FirstName: "Somchai", LastName: "Deejai"},
		&userdomain.Profile{UserID: "usr_2", AvatarURL: "somchai.png"},
		nil,
	)

	req := httptest.NewRequest("GET", "/api/v1/users/usr_2", nil)
	ctx := middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"})

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "usr_2")
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)

	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.GetPublicProfile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "usr_2", resp["user_id"])
	assert.Equal(t, "Somchai", resp["first_name"])
	assert.Equal(t, "Deejai", resp["last_name"])
	assert.Equal(t, "somchai.png", resp["avatar_url"])
}

func TestUserHandler_AssignJob_Success(t *testing.T) {
	userUC := new(MockUserUC)
	h := userhandler.NewUserHandler(userUC)

	userID := "usr_123"
	jobID := "job_456"

	userUC.On("AssignJob", mock.Anything, userID, jobID, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	body := `{"job_id":"job_456"}`
	req := httptest.NewRequest("POST", "/api/v1/users/me/job", bytes.NewBufferString(body))
	ctx := middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: userID})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.AssignJob(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	userUC.AssertExpectations(t)
}

func TestUserHandler_UpdatePassword_Success(t *testing.T) {
	userUC := new(MockUserUC)
	h := userhandler.NewUserHandler(userUC)

	userID := "usr_123"
	currentPass := "oldpassword"
	newPass := "newpassword"

	userUC.On("UpdatePassword", mock.Anything, userID, currentPass, newPass).Return(nil)

	body := `{"current_password":"oldpassword","new_password":"newpassword"}`
	req := httptest.NewRequest("PUT", "/api/v1/users/me/password", bytes.NewBufferString(body))
	ctx := middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: userID})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.UpdatePassword(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	userUC.AssertExpectations(t)
}

func TestUserHandler_UpdatePassword_Unauthorized(t *testing.T) {
	userUC := new(MockUserUC)
	h := userhandler.NewUserHandler(userUC)

	// No claims in context
	body := `{"current_password":"oldpassword","new_password":"newpassword"}`
	req := httptest.NewRequest("PUT", "/api/v1/users/me/password", bytes.NewBufferString(body))

	w := httptest.NewRecorder()
	h.UpdatePassword(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_UpdatePassword_ValidationFailed(t *testing.T) {
	userUC := new(MockUserUC)
	h := userhandler.NewUserHandler(userUC)

	userID := "usr_123"

	// Missing new_password or too short
	body := `{"current_password":"oldpassword","new_password":"short"}`
	req := httptest.NewRequest("PUT", "/api/v1/users/me/password", bytes.NewBufferString(body))
	ctx := middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: userID})
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	h.UpdatePassword(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

