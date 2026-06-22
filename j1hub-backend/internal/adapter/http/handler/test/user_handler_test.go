package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/j1hub/backend/internal/adapter/http/handler"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUC
type MockUserUC struct {
	mock.Mock
}

func (m *MockUserUC) GetProfile(ctx context.Context, userID string) (*usecase.UserProfileResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.UserProfileResponse), args.Error(1)
}

func (m *MockUserUC) GetPublicProfile(ctx context.Context, currentUserID, targetUserID string) (*domain.User, *domain.Profile, error) {
	args := m.Called(ctx, currentUserID, targetUserID)
	var u *domain.User
	if args.Get(0) != nil {
		u = args.Get(0).(*domain.User)
	}
	var p *domain.Profile
	if args.Get(1) != nil {
		p = args.Get(1).(*domain.Profile)
	}
	return u, p, args.Error(2)
}

func (m *MockUserUC) UpdateProfile(ctx context.Context, userID string, cmd usecase.UpdateProfileCommand) error {
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

func TestUserHandler_GetProfile_Success(t *testing.T) {
	userUC := new(MockUserUC)
	h := handler.NewUserHandler(userUC)

	userID := "usr_123"
	profileResp := &usecase.UserProfileResponse{
		User:    &domain.User{UserID: userID, Email: "john@example.com"},
		Profile: &domain.Profile{UserID: userID, Bio: "Hello world"},
	}

	userUC.On("GetProfile", mock.Anything, userID).Return(profileResp, nil)

	req := httptest.NewRequest("GET", "/api/v1/user/profile", nil)
	ctx := middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: userID})
	req = req.WithContext(ctx)
	
	w := httptest.NewRecorder()

	h.GetProfile(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp usecase.UserProfileResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "john@example.com", resp.User.Email)
	assert.Equal(t, "Hello world", resp.Profile.Bio)
}

func TestUserHandler_GetProfile_Unauthorized(t *testing.T) {
	userUC := new(MockUserUC)
	h := handler.NewUserHandler(userUC)

	// No claims in context
	req := httptest.NewRequest("GET", "/api/v1/user/profile", nil)
	w := httptest.NewRecorder()

	h.GetProfile(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandler_UpdateProfile_Success(t *testing.T) {
	userUC := new(MockUserUC)
	h := handler.NewUserHandler(userUC)

	userID := "usr_123"
	cmd := usecase.UpdateProfileCommand{
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
	h := handler.NewUserHandler(userUC)

	userUC.On("GetPublicProfile", mock.Anything, "usr_1", "usr_2").Return(
		&domain.User{UserID: "usr_2", FirstName: "Somchai", LastName: "Deejai"},
		&domain.Profile{UserID: "usr_2", AvatarURL: "somchai.png"},
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
