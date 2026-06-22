package handler_test
import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestFriendHandler_SendRequest(t *testing.T) {
	friendshipUC := new(MockFriendshipUC)
	h := handler.NewFriendHandler(friendshipUC, nil)

	// unauthorized
	req := httptest.NewRequest("POST", "/friends/request", nil)
	w := httptest.NewRecorder()
	h.SendRequest(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// bad body
	req = httptest.NewRequest("POST", "/friends/request", strings.NewReader("bad_json"))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.SendRequest(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// success
	friendshipUC.On("SendRequest", mock.Anything, "usr_1", "usr_2").Return(nil).Once()
	req = httptest.NewRequest("POST", "/friends/request", strings.NewReader(`{"target_user_id":"usr_2"}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.SendRequest(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// error
	friendshipUC.On("SendRequest", mock.Anything, "usr_1", "usr_2").Return(errors.New("err")).Once()
	req = httptest.NewRequest("POST", "/friends/request", strings.NewReader(`{"target_user_id":"usr_2"}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.SendRequest(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestFriendHandler_ListPendingRequests(t *testing.T) {
	friendshipUC := new(MockFriendshipUC)
	h := handler.NewFriendHandler(friendshipUC, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/friends/requests/pending", nil)
	w := httptest.NewRecorder()
	h.ListPendingRequests(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	friendshipUC.On("ListPendingRequests", mock.Anything, "usr_1").Return([]domain.Friendship{}, nil).Once()
	req = httptest.NewRequest("GET", "/friends/requests/pending", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListPendingRequests(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	friendshipUC.On("ListPendingRequests", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/friends/requests/pending", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListPendingRequests(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestFriendHandler_RespondToRequest(t *testing.T) {
	friendshipUC := new(MockFriendshipUC)
	h := handler.NewFriendHandler(friendshipUC, nil)

	// unauthorized
	req := httptest.NewRequest("PATCH", "/friends/respond", nil)
	w := httptest.NewRecorder()
	h.RespondToRequest(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// bad body
	req = httptest.NewRequest("PATCH", "/friends/respond", strings.NewReader("bad_json"))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.RespondToRequest(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// success
	friendshipUC.On("RespondToRequest", mock.Anything, "usr_1", "fr_1", true).Return(nil).Once()
	req = httptest.NewRequest("PATCH", "/friends/respond", strings.NewReader(`{"friendship_id":"fr_1","accept":true}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "fr_1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.RespondToRequest(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	friendshipUC.On("RespondToRequest", mock.Anything, "usr_1", "fr_1", true).Return(errors.New("err")).Once()
	req = httptest.NewRequest("PATCH", "/friends/respond", strings.NewReader(`{"friendship_id":"fr_1","accept":true}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.RespondToRequest(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestFriendHandler_ListFriends(t *testing.T) {
	friendshipUC := new(MockFriendshipUC)
	h := handler.NewFriendHandler(friendshipUC, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/friends", nil)
	w := httptest.NewRecorder()
	h.ListFriends(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	friendshipUC.On("ListFriends", mock.Anything, "usr_1").Return([]domain.Friendship{}, nil).Once()
	req = httptest.NewRequest("GET", "/friends", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListFriends(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	friendshipUC.On("ListFriends", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/friends", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListFriends(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestFriendHandler_RemoveFriend(t *testing.T) {
	friendshipUC := new(MockFriendshipUC)
	h := handler.NewFriendHandler(friendshipUC, nil)

	// unauthorized
	req := httptest.NewRequest("DELETE", "/friends/usr_2", nil)
	w := httptest.NewRecorder()
	h.RemoveFriend(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	friendshipUC.On("RemoveFriend", mock.Anything, "usr_1", "usr_2").Return(nil).Once()
	req = httptest.NewRequest("DELETE", "/friends/usr_2", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "usr_2")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.RemoveFriend(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	friendshipUC.On("RemoveFriend", mock.Anything, "usr_1", "usr_2").Return(errors.New("err")).Once()
	req = httptest.NewRequest("DELETE", "/friends/usr_2", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.RemoveFriend(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestFriendHandler_GetRadar(t *testing.T) {
	radarUC := new(MockRadarUC)
	h := handler.NewFriendHandler(nil, radarUC)

	// unauthorized
	req := httptest.NewRequest("GET", "/radar", nil)
	w := httptest.NewRecorder()
	h.GetRadar(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	radarUC.On("GetRadar", mock.Anything, "usr_1").Return([]usecase.RadarEntry{}, nil).Once()
	req = httptest.NewRequest("GET", "/radar", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.GetRadar(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	radarUC.On("GetRadar", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/radar", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.GetRadar(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
