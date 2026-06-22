package handler_test
import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
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

func TestMissionHandler_ListUserMissions(t *testing.T) {
	missionUC := new(MockMissionUC)
	h := handler.NewMissionHandler(missionUC, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/user-missions", nil)
	w := httptest.NewRecorder()
	h.ListUserMissions(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	missionUC.On("ListAvailableMissions", mock.Anything, "usr_1").Return([]domain.UserMission{}, nil).Once()
	req = httptest.NewRequest("GET", "/user-missions", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListUserMissions(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	missionUC.On("ListAvailableMissions", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/user-missions", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListUserMissions(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestMissionHandler_ListMissions(t *testing.T) {
	missionUC := new(MockMissionUC)
	h := handler.NewMissionHandler(missionUC, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/missions", nil)
	w := httptest.NewRecorder()
	h.ListMissions(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	missionUC.On("ListStaticMissions", mock.Anything, "usr_1").Return([]domain.Mission{}, nil).Once()
	req = httptest.NewRequest("GET", "/missions", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListMissions(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	missionUC.On("ListStaticMissions", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/missions", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListMissions(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestMissionHandler_GetMissionDetail(t *testing.T) {
	missionUC := new(MockMissionUC)
	h := handler.NewMissionHandler(missionUC, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/missions/m1", nil)
	w := httptest.NewRecorder()
	h.GetMissionDetail(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	missionUC.On("GetMissionDetail", mock.Anything, "usr_1", "m1").Return(&usecase.MissionDetailResponse{}, nil).Once()
	req = httptest.NewRequest("GET", "/missions/m1", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "m1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.GetMissionDetail(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	missionUC.On("GetMissionDetail", mock.Anything, "usr_1", "m1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/missions/m1", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.GetMissionDetail(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestMissionHandler_SubmitProof(t *testing.T) {
	completeUC := new(MockCompleteMissionUC)
	h := handler.NewMissionHandler(nil, completeUC)

	// unauthorized
	req := httptest.NewRequest("POST", "/missions/m1/verify", nil)
	w := httptest.NewRecorder()
	h.SubmitProof(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// no file
	req = httptest.NewRequest("POST", "/missions/m1/verify", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.SubmitProof(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// success
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("proof", "proof.jpg")
	part.Write([]byte("proof_data"))
	writer.Close()

	completeUC.On("SubmitProof", mock.Anything, "usr_1", "m1", mock.Anything, mock.Anything).Return(nil).Once()
	req = httptest.NewRequest("POST", "/missions/m1/verify", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "m1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.SubmitProof(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	part, _ = writer.CreateFormFile("proof", "proof.jpg")
	part.Write([]byte("proof_data"))
	writer.Close()

	completeUC.On("SubmitProof", mock.Anything, "usr_1", "m1", mock.Anything, mock.Anything).Return(errors.New("err")).Once()
	req = httptest.NewRequest("POST", "/missions/m1/verify", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.SubmitProof(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}


func TestMissionHandler_ToggleTask(t *testing.T) {
	missionUC := new(MockMissionUC)
	h := handler.NewMissionHandler(missionUC, nil)

	// unauthorized
	req := httptest.NewRequest("PATCH", "/tasks/t1", nil)
	w := httptest.NewRecorder()
	h.ToggleTask(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// bad body
	req = httptest.NewRequest("PATCH", "/tasks/t1", strings.NewReader("bad_json"))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ToggleTask(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// success
	missionUC.On("ToggleTask", mock.Anything, "usr_1", "t1", true).Return(nil).Once()
	req = httptest.NewRequest("PATCH", "/tasks/t1", strings.NewReader(`{"completed":true}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "t1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.ToggleTask(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	missionUC.On("ToggleTask", mock.Anything, "usr_1", "t1", true).Return(errors.New("err")).Once()
	req = httptest.NewRequest("PATCH", "/tasks/t1", strings.NewReader(`{"completed":true}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.ToggleTask(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
