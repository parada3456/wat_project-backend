package handler_test
import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/go-chi/chi/v5"
	"github.com/j1hub/backend/internal/adapter/http/handler"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotificationHandler_ListNotifications(t *testing.T) {
	notifUC := new(MockNotificationUC)
	h := handler.NewNotificationHandler(notifUC)

	// unauthorized
	req := httptest.NewRequest("GET", "/notifications", nil)
	w := httptest.NewRecorder()
	h.ListNotifications(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	notifUC.On("ListNotifications", mock.Anything, "usr_1").Return([]domain.Notification{}, nil).Once()
	req = httptest.NewRequest("GET", "/notifications", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListNotifications(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	notifUC.On("ListNotifications", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/notifications", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListNotifications(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestNotificationHandler_MarkRead(t *testing.T) {
	notifUC := new(MockNotificationUC)
	h := handler.NewNotificationHandler(notifUC)

	// success
	notifUC.On("MarkRead", mock.Anything, "n1").Return(nil).Once()
	req := httptest.NewRequest("PATCH", "/notifications/n1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "n1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	h.MarkRead(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	notifUC.On("MarkRead", mock.Anything, "n1").Return(errors.New("err")).Once()
	req = httptest.NewRequest("PATCH", "/notifications/n1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.MarkRead(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestNotificationHandler_MarkAllRead(t *testing.T) {
	notifUC := new(MockNotificationUC)
	h := handler.NewNotificationHandler(notifUC)

	// unauthorized
	req := httptest.NewRequest("PATCH", "/notifications", nil)
	w := httptest.NewRecorder()
	h.MarkAllRead(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	notifUC.On("MarkAllRead", mock.Anything, "usr_1").Return(nil).Once()
	req = httptest.NewRequest("PATCH", "/notifications", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.MarkAllRead(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	notifUC.On("MarkAllRead", mock.Anything, "usr_1").Return(errors.New("err")).Once()
	req = httptest.NewRequest("PATCH", "/notifications", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.MarkAllRead(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestNotificationHandler_DeleteNotification(t *testing.T) {
	notifUC := new(MockNotificationUC)
	h := handler.NewNotificationHandler(notifUC)

	// success
	notifUC.On("Delete", mock.Anything, "n1").Return(nil).Once()
	req := httptest.NewRequest("DELETE", "/notifications/n1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "n1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	h.DeleteNotification(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	notifUC.On("Delete", mock.Anything, "n1").Return(errors.New("err")).Once()
	req = httptest.NewRequest("DELETE", "/notifications/n1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.DeleteNotification(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

type MockTokenIssuer struct {
	mock.Mock
}

func (m *MockTokenIssuer) Issue(userID string, isAdmin bool) (*port.TokenPair, error) {
	args := m.Called(userID, isAdmin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*port.TokenPair), args.Error(1)
}
