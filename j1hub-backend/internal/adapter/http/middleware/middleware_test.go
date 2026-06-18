package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTokenIssuer struct {
	mock.Mock
}

func (m *MockTokenIssuer) Issue(userID string, isAdmin bool) (*port.TokenPair, error) {
	args := m.Called(userID, isAdmin)
	return args.Get(0).(*port.TokenPair), args.Error(1)
}

func (m *MockTokenIssuer) Verify(token string) (*port.Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.Claims), args.Error(1)
}

func (m *MockTokenIssuer) Refresh(token string) (*port.TokenPair, error) {
	args := m.Called(token)
	return args.Get(0).(*port.TokenPair), args.Error(1)
}

func TestAuthenticate(t *testing.T) {
	issuer := new(MockTokenIssuer)
	authMW := middleware.Authenticate(issuer)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		assert.NotNil(t, claims)
		assert.Equal(t, "usr_1", claims.UserID)
		assert.True(t, claims.IsAdmin)
		w.WriteHeader(http.StatusOK)
	})

	// Missing header
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	authMW(next).ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Invalid token
	issuer.On("Verify", "bad_token").Return(nil, errors.New("invalid")).Once()
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer bad_token")
	w = httptest.NewRecorder()
	authMW(next).ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Valid token
	issuer.On("Verify", "good_token").Return(&port.Claims{UserID: "usr_1", IsAdmin: true}, nil).Once()
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer good_token")
	w = httptest.NewRecorder()
	authMW(next).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireAdmin(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Not admin
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	middleware.RequireAdmin(next).ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Admin
	req = httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1", IsAdmin: true}))
	w = httptest.NewRecorder()
	middleware.RequireAdmin(next).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetClaims_Nil(t *testing.T) {
	claims := middleware.GetClaims(context.Background())
	assert.Nil(t, claims)
}

func TestLogger(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	middleware.Logger(next).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
