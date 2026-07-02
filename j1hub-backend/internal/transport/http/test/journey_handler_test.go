package test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	gamificationusecase "github.com/parada3456/wat_project-backend/internal/gamification/usecase"

	missiondomain "github.com/parada3456/wat_project-backend/internal/mission/domain"

	gamificationdomain "github.com/parada3456/wat_project-backend/internal/gamification/domain"

	port "github.com/parada3456/wat_project-backend/internal/auth/port"
	journeyhandler "github.com/parada3456/wat_project-backend/internal/gamification/adapter/http"
	"github.com/parada3456/wat_project-backend/internal/transport/http/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJourneyHandler_ListPhases(t *testing.T) {
	journeyUC := new(MockJourneyUC)
	h := journeyhandler.NewJourneyHandler(journeyUC, nil, nil)

	// success path
	journeyUC.On("ListPhases", mock.Anything).Return([]missiondomain.JourneyPhase{{PhaseID: "p1"}}, nil).Once()
	req := httptest.NewRequest("GET", "/journey/phases", nil)
	w := httptest.NewRecorder()
	h.ListPhases(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error path
	journeyUC.On("ListPhases", mock.Anything).Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/journey/phases", nil)
	w = httptest.NewRecorder()
	h.ListPhases(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJourneyHandler_AdvancePhase(t *testing.T) {
	advanceUC := new(MockAdvancePhaseUC)
	h := journeyhandler.NewJourneyHandler(nil, advanceUC, nil)

	// unauthorized
	req := httptest.NewRequest("POST", "/journey/phase/transition", nil)
	w := httptest.NewRecorder()
	h.AdvancePhase(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	advanceUC.On("TryAdvancePhase", mock.Anything, "usr_1").Return(&gamificationusecase.PhaseTransitionResponse{Transitioned: true}, nil).Once()
	req = httptest.NewRequest("POST", "/journey/phase/transition", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.AdvancePhase(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	advanceUC.On("TryAdvancePhase", mock.Anything, "usr_1").Return((*gamificationusecase.PhaseTransitionResponse)(nil), errors.New("err")).Once()
	req = httptest.NewRequest("POST", "/journey/phase/transition", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.AdvancePhase(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJourneyHandler_GetHistory(t *testing.T) {
	journeyUC := new(MockJourneyUC)
	h := journeyhandler.NewJourneyHandler(journeyUC, nil, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/journey/history", nil)
	w := httptest.NewRecorder()
	h.GetHistory(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	journeyUC.On("GetHistory", mock.Anything, "usr_1").Return([]missiondomain.UserPhaseHistory{}, nil).Once()
	req = httptest.NewRequest("GET", "/journey/history", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.GetHistory(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	journeyUC.On("GetHistory", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/journey/history", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.GetHistory(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJourneyHandler_GetLeaderboard(t *testing.T) {
	leaderboardUC := new(MockLeaderboardUC)
	h := journeyhandler.NewJourneyHandler(nil, nil, leaderboardUC)

	// success
	leaderboardUC.On("GetLeaderboard", mock.Anything, "global", "").Return([]gamificationusecase.LeaderboardEntry{}, nil).Once()
	req := httptest.NewRequest("GET", "/leaderboard?scope=global", nil)
	w := httptest.NewRecorder()
	h.GetLeaderboard(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	leaderboardUC.On("GetLeaderboard", mock.Anything, "global", "").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/leaderboard?scope=global", nil)
	w = httptest.NewRecorder()
	h.GetLeaderboard(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJourneyHandler_ListBadges(t *testing.T) {
	journeyUC := new(MockJourneyUC)
	h := journeyhandler.NewJourneyHandler(journeyUC, nil, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/user/badges", nil)
	w := httptest.NewRecorder()
	h.ListBadges(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	journeyUC.On("ListUserBadges", mock.Anything, "usr_1").Return([]gamificationdomain.UserBadge{}, nil).Once()
	req = httptest.NewRequest("GET", "/user/badges", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListBadges(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	journeyUC.On("ListUserBadges", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/user/badges", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListBadges(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJourneyHandler_GetCreditHistory(t *testing.T) {
	journeyUC := new(MockJourneyUC)
	h := journeyhandler.NewJourneyHandler(journeyUC, nil, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/user/credit-score/history", nil)
	w := httptest.NewRecorder()
	h.GetCreditHistory(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	journeyUC.On("GetCreditScoreHistory", mock.Anything, "usr_1").Return([]gamificationdomain.PointLedger{}, nil).Once()
	req = httptest.NewRequest("GET", "/user/credit-score/history", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.GetCreditHistory(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	journeyUC.On("GetCreditScoreHistory", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/user/credit-score/history", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.GetCreditHistory(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJourneyHandler_GetPointsLedger(t *testing.T) {
	journeyUC := new(MockJourneyUC)
	h := journeyhandler.NewJourneyHandler(journeyUC, nil, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/user/points/ledger", nil)
	w := httptest.NewRecorder()
	h.GetPointsLedger(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	journeyUC.On("GetPointsLedger", mock.Anything, "usr_1").Return([]gamificationdomain.PointLedger{}, nil).Once()
	req = httptest.NewRequest("GET", "/user/points/ledger", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.GetPointsLedger(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	journeyUC.On("GetPointsLedger", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/user/points/ledger", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.GetPointsLedger(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
