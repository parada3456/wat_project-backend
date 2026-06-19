package handler_test

import (
	"bytes"
	"context"
	"errors"
	"io"
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

// MockJourneyUC
type MockJourneyUC struct {
	mock.Mock
}

func (m *MockJourneyUC) ListPhases(ctx context.Context) ([]domain.JourneyPhase, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.JourneyPhase), args.Error(1)
}

func (m *MockJourneyUC) GetHistory(ctx context.Context, userID string) ([]domain.UserPhaseHistory, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserPhaseHistory), args.Error(1)
}

func (m *MockJourneyUC) ListUserBadges(ctx context.Context, userID string) ([]domain.UserBadge, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserBadge), args.Error(1)
}

func (m *MockJourneyUC) GetCreditScoreHistory(ctx context.Context, userID string) ([]domain.PointLedger, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.PointLedger), args.Error(1)
}

// MockAdvancePhaseUC
type MockAdvancePhaseUC struct {
	mock.Mock
}

func (m *MockAdvancePhaseUC) TryAdvancePhase(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

// MockLeaderboardUC
type MockLeaderboardUC struct {
	mock.Mock
}

func (m *MockLeaderboardUC) GetLeaderboard(ctx context.Context, scope, jobID string) ([]usecase.LeaderboardEntry, error) {
	args := m.Called(ctx, scope, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]usecase.LeaderboardEntry), args.Error(1)
}

// MockManageExpenseUC
type MockManageExpenseUC struct {
	mock.Mock
}

func (m *MockManageExpenseUC) ListExpenses(ctx context.Context, userID string) ([]domain.ExpenseTransaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ExpenseTransaction), args.Error(1)
}

func (m *MockManageExpenseUC) CreateExpense(ctx context.Context, userID string, cmd usecase.CreateExpenseCmd) error {
	return m.Called(ctx, userID, cmd).Error(0)
}

func (m *MockManageExpenseUC) GetExpenseDetail(ctx context.Context, userID, id string) (*domain.ExpenseTransaction, []domain.ExpenseSplit, error) {
	args := m.Called(ctx, userID, id)
	var t *domain.ExpenseTransaction
	if args.Get(0) != nil {
		t = args.Get(0).(*domain.ExpenseTransaction)
	}
	var s []domain.ExpenseSplit
	if args.Get(1) != nil {
		s = args.Get(1).([]domain.ExpenseSplit)
	}
	return t, s, args.Error(2)
}

func (m *MockManageExpenseUC) DeleteExpense(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func (m *MockManageExpenseUC) ListPendingExpenses(ctx context.Context, userID string) ([]domain.ExpenseSplit, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ExpenseSplit), args.Error(1)
}

func (m *MockManageExpenseUC) SubmitSlip(ctx context.Context, debtorID, splitID string, file io.Reader, contentType string) error {
	return m.Called(ctx, debtorID, splitID, mock.Anything, contentType).Error(0)
}

func (m *MockManageExpenseUC) ApproveSplit(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

// MockFriendshipUC
type MockFriendshipUC struct {
	mock.Mock
}

func (m *MockFriendshipUC) SendRequest(ctx context.Context, senderID, targetID string) error {
	return m.Called(ctx, senderID, targetID).Error(0)
}

func (m *MockFriendshipUC) ListPendingRequests(ctx context.Context, userID string) ([]domain.Friendship, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Friendship), args.Error(1)
}

func (m *MockFriendshipUC) RespondToRequest(ctx context.Context, userID, friendshipID string, accept bool) error {
	return m.Called(ctx, userID, friendshipID, accept).Error(0)
}

func (m *MockFriendshipUC) ListFriends(ctx context.Context, userID string) ([]domain.Friendship, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Friendship), args.Error(1)
}

func (m *MockFriendshipUC) RemoveFriend(ctx context.Context, userID, friendID string) error {
	return m.Called(ctx, userID, friendID).Error(0)
}

// MockRadarUC
type MockRadarUC struct {
	mock.Mock
}

func (m *MockRadarUC) GetRadar(ctx context.Context, requesterID string) ([]usecase.RadarEntry, error) {
	args := m.Called(ctx, requesterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]usecase.RadarEntry), args.Error(1)
}

// MockJobUC
type MockJobUC struct {
	mock.Mock
}

func (m *MockJobUC) ListJobs(ctx context.Context, filters map[string]interface{}) ([]domain.JobPosting, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.JobPosting), args.Error(1)
}

func (m *MockJobUC) GetJobDetail(ctx context.Context, id string) (*domain.JobPosting, []domain.JobHousing, *domain.JobOverallRating, error) {
	args := m.Called(ctx, id)
	var j *domain.JobPosting
	if args.Get(0) != nil {
		j = args.Get(0).(*domain.JobPosting)
	}
	var h []domain.JobHousing
	if args.Get(1) != nil {
		h = args.Get(1).([]domain.JobHousing)
	}
	var r *domain.JobOverallRating
	if args.Get(2) != nil {
		r = args.Get(2).(*domain.JobOverallRating)
	}
	return j, h, r, args.Error(3)
}

func (m *MockJobUC) AddToCart(ctx context.Context, userID, jobID string) error {
	return m.Called(ctx, userID, jobID).Error(0)
}

func (m *MockJobUC) ListCart(ctx context.Context, userID string) ([]domain.UserCart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserCart), args.Error(1)
}

func (m *MockJobUC) RemoveFromCart(ctx context.Context, userID, id string) error {
	return m.Called(ctx, userID, id).Error(0)
}

func (m *MockJobUC) WriteReview(ctx context.Context, userID, jobID string, rev *domain.JobReview) error {
	return m.Called(ctx, userID, jobID, rev).Error(0)
}

// MockMissionUC
type MockMissionUC struct {
	mock.Mock
}

func (m *MockMissionUC) ListAvailableMissions(ctx context.Context, userID string) ([]domain.UserMission, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserMission), args.Error(1)
}

func (m *MockMissionUC) GetMissionDetail(ctx context.Context, userID, userMissionID string) (*usecase.MissionDetailResponse, error) {
	args := m.Called(ctx, userID, userMissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.MissionDetailResponse), args.Error(1)
}

func (m *MockMissionUC) ToggleTask(ctx context.Context, userID, userTaskID string, completed bool) error {
	return m.Called(ctx, userID, userTaskID, completed).Error(0)
}

// MockCompleteMissionUC
type MockCompleteMissionUC struct {
	mock.Mock
}

func (m *MockCompleteMissionUC) SubmitProof(ctx context.Context, userID, userMissionID string, file io.Reader, contentType string) error {
	return m.Called(ctx, userID, userMissionID, mock.Anything, contentType).Error(0)
}

// MockNotificationUC
type MockNotificationUC struct {
	mock.Mock
}

func (m *MockNotificationUC) ListNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
}

func (m *MockNotificationUC) MarkRead(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockNotificationUC) MarkAllRead(ctx context.Context, userID string) error {
	return m.Called(ctx, userID).Error(0)
}

func (m *MockNotificationUC) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// Tests
func TestJourneyHandler_ListPhases(t *testing.T) {
	journeyUC := new(MockJourneyUC)
	h := handler.NewJourneyHandler(journeyUC, nil, nil)

	// success path
	journeyUC.On("ListPhases", mock.Anything).Return([]domain.JourneyPhase{{PhaseID: "p1"}}, nil).Once()
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
	h := handler.NewJourneyHandler(nil, advanceUC, nil)

	// unauthorized
	req := httptest.NewRequest("POST", "/journey/phase/transition", nil)
	w := httptest.NewRecorder()
	h.AdvancePhase(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	advanceUC.On("TryAdvancePhase", mock.Anything, "usr_1").Return(true, nil).Once()
	req = httptest.NewRequest("POST", "/journey/phase/transition", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.AdvancePhase(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	advanceUC.On("TryAdvancePhase", mock.Anything, "usr_1").Return(false, errors.New("err")).Once()
	req = httptest.NewRequest("POST", "/journey/phase/transition", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.AdvancePhase(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJourneyHandler_GetHistory(t *testing.T) {
	journeyUC := new(MockJourneyUC)
	h := handler.NewJourneyHandler(journeyUC, nil, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/journey/history", nil)
	w := httptest.NewRecorder()
	h.GetHistory(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	journeyUC.On("GetHistory", mock.Anything, "usr_1").Return([]domain.UserPhaseHistory{}, nil).Once()
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
	h := handler.NewJourneyHandler(nil, nil, leaderboardUC)

	// success
	leaderboardUC.On("GetLeaderboard", mock.Anything, "global", "").Return([]usecase.LeaderboardEntry{}, nil).Once()
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
	h := handler.NewJourneyHandler(journeyUC, nil, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/user/badges", nil)
	w := httptest.NewRecorder()
	h.ListBadges(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	journeyUC.On("ListUserBadges", mock.Anything, "usr_1").Return([]domain.UserBadge{}, nil).Once()
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
	h := handler.NewJourneyHandler(journeyUC, nil, nil)

	// unauthorized
	req := httptest.NewRequest("GET", "/user/credit-score/history", nil)
	w := httptest.NewRecorder()
	h.GetCreditHistory(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	journeyUC.On("GetCreditScoreHistory", mock.Anything, "usr_1").Return([]domain.PointLedger{}, nil).Once()
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

func TestExpenseHandler_ListExpenses(t *testing.T) {
	expenseUC := new(MockManageExpenseUC)
	h := handler.NewExpenseHandler(expenseUC)

	// unauthorized
	req := httptest.NewRequest("GET", "/expenses", nil)
	w := httptest.NewRecorder()
	h.ListExpenses(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	expenseUC.On("ListExpenses", mock.Anything, "usr_1").Return([]domain.ExpenseTransaction{}, nil).Once()
	req = httptest.NewRequest("GET", "/expenses", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListExpenses(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	expenseUC.On("ListExpenses", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/expenses", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListExpenses(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestExpenseHandler_CreateExpense(t *testing.T) {
	expenseUC := new(MockManageExpenseUC)
	h := handler.NewExpenseHandler(expenseUC)

	// unauthorized
	req := httptest.NewRequest("POST", "/expenses", nil)
	w := httptest.NewRecorder()
	h.CreateExpense(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// bad body
	req = httptest.NewRequest("POST", "/expenses", strings.NewReader("bad_json"))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.CreateExpense(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// validation fail
	req = httptest.NewRequest("POST", "/expenses", strings.NewReader(`{"title":""}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.CreateExpense(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// success
	body := `{"title":"Dinner","total_amount":100,"currency":"USD","due_date":"2026-06-17T15:00:00Z","splits":[{"user_id":"usr_2","owe_amount":50}]}`
	expenseUC.On("CreateExpense", mock.Anything, "usr_1", mock.Anything).Return(nil).Once()
	req = httptest.NewRequest("POST", "/expenses", strings.NewReader(body))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.CreateExpense(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// error
	expenseUC.On("CreateExpense", mock.Anything, "usr_1", mock.Anything).Return(errors.New("err")).Once()
	req = httptest.NewRequest("POST", "/expenses", strings.NewReader(body))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.CreateExpense(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestExpenseHandler_GetExpenseDetail(t *testing.T) {
	expenseUC := new(MockManageExpenseUC)
	h := handler.NewExpenseHandler(expenseUC)

	// unauthorized
	req := httptest.NewRequest("GET", "/expenses/exp_1", nil)
	w := httptest.NewRecorder()
	h.GetExpenseDetail(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	expenseUC.On("GetExpenseDetail", mock.Anything, "usr_1", "exp_1").Return(&domain.ExpenseTransaction{}, []domain.ExpenseSplit{}, nil).Once()
	req = httptest.NewRequest("GET", "/expenses/exp_1", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "exp_1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w = httptest.NewRecorder()
	h.GetExpenseDetail(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	expenseUC.On("GetExpenseDetail", mock.Anything, "usr_1", "exp_1").Return(nil, nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/expenses/exp_1", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.GetExpenseDetail(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestExpenseHandler_DeleteExpense(t *testing.T) {
	expenseUC := new(MockManageExpenseUC)
	h := handler.NewExpenseHandler(expenseUC)

	// unauthorized
	req := httptest.NewRequest("DELETE", "/expenses/exp_1", nil)
	w := httptest.NewRecorder()
	h.DeleteExpense(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	expenseUC.On("DeleteExpense", mock.Anything, "usr_1", "exp_1").Return(nil).Once()
	req = httptest.NewRequest("DELETE", "/expenses/exp_1", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "exp_1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.DeleteExpense(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	expenseUC.On("DeleteExpense", mock.Anything, "usr_1", "exp_1").Return(errors.New("err")).Once()
	req = httptest.NewRequest("DELETE", "/expenses/exp_1", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.DeleteExpense(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestExpenseHandler_ListPending(t *testing.T) {
	expenseUC := new(MockManageExpenseUC)
	h := handler.NewExpenseHandler(expenseUC)

	// unauthorized
	req := httptest.NewRequest("GET", "/expenses/pending", nil)
	w := httptest.NewRecorder()
	h.ListPending(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	expenseUC.On("ListPendingExpenses", mock.Anything, "usr_1").Return([]domain.ExpenseSplit{}, nil).Once()
	req = httptest.NewRequest("GET", "/expenses/pending", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListPending(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	expenseUC.On("ListPendingExpenses", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/expenses/pending", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListPending(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestExpenseHandler_PaySplit(t *testing.T) {
	expenseUC := new(MockManageExpenseUC)
	h := handler.NewExpenseHandler(expenseUC)

	// unauthorized
	req := httptest.NewRequest("POST", "/expenses/splits/s1/pay", nil)
	w := httptest.NewRecorder()
	h.PaySplit(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// multipart error (no file)
	req = httptest.NewRequest("POST", "/expenses/splits/s1/pay", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.PaySplit(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// success
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("slip", "slip.jpg")
	part.Write([]byte("image_data"))
	writer.Close()

	expenseUC.On("SubmitSlip", mock.Anything, "usr_1", "s1", mock.Anything, mock.Anything).Return(nil).Once()
	req = httptest.NewRequest("POST", "/expenses/splits/s1/pay", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "s1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w = httptest.NewRecorder()
	h.PaySplit(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	part, _ = writer.CreateFormFile("slip", "slip.jpg")
	part.Write([]byte("image_data"))
	writer.Close()

	expenseUC.On("SubmitSlip", mock.Anything, "usr_1", "s1", mock.Anything, mock.Anything).Return(errors.New("err")).Once()
	req = httptest.NewRequest("POST", "/expenses/splits/s1/pay", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w = httptest.NewRecorder()
	h.PaySplit(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestExpenseHandler_ApproveSplit(t *testing.T) {
	expenseUC := new(MockManageExpenseUC)
	h := handler.NewExpenseHandler(expenseUC)

	// unauthorized
	req := httptest.NewRequest("PATCH", "/expenses/splits/s1/approve", nil)
	w := httptest.NewRecorder()
	h.ApproveSplit(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	expenseUC.On("ApproveSplit", mock.Anything, "usr_1", "s1").Return(nil).Once()
	req = httptest.NewRequest("PATCH", "/expenses/splits/s1/approve", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "s1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.ApproveSplit(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	expenseUC.On("ApproveSplit", mock.Anything, "usr_1", "s1").Return(errors.New("err")).Once()
	req = httptest.NewRequest("PATCH", "/expenses/splits/s1/approve", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.ApproveSplit(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

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
	w = httptest.NewRecorder()
	h.RespondToRequest(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	friendshipUC.On("RespondToRequest", mock.Anything, "usr_1", "fr_1", true).Return(errors.New("err")).Once()
	req = httptest.NewRequest("PATCH", "/friends/respond", strings.NewReader(`{"friendship_id":"fr_1","accept":true}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
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

func TestJobHandler_ListJobs(t *testing.T) {
	jobUC := new(MockJobUC)
	h := handler.NewJobHandler(jobUC)

	// success
	jobUC.On("ListJobs", mock.Anything, map[string]interface{}{"agency_name": "agency_1"}).Return([]domain.JobPosting{}, nil).Once()
	req := httptest.NewRequest("GET", "/jobs?agency=agency_1", nil)
	w := httptest.NewRecorder()
	h.ListJobs(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	jobUC.On("ListJobs", mock.Anything, map[string]interface{}{}).Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/jobs", nil)
	w = httptest.NewRecorder()
	h.ListJobs(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJobHandler_GetJobDetail(t *testing.T) {
	jobUC := new(MockJobUC)
	h := handler.NewJobHandler(jobUC)

	// success
	jobUC.On("GetJobDetail", mock.Anything, "job_1").Return(&domain.JobPosting{}, []domain.JobHousing{}, &domain.JobOverallRating{}, nil).Once()
	req := httptest.NewRequest("GET", "/jobs/job_1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "job_1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	h.GetJobDetail(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	jobUC.On("GetJobDetail", mock.Anything, "job_1").Return(nil, nil, nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/jobs/job_1", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.GetJobDetail(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJobHandler_AddToCart(t *testing.T) {
	jobUC := new(MockJobUC)
	h := handler.NewJobHandler(jobUC)

	// unauthorized
	req := httptest.NewRequest("POST", "/cart", nil)
	w := httptest.NewRecorder()
	h.AddToCart(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// bad body
	req = httptest.NewRequest("POST", "/cart", strings.NewReader("bad_json"))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.AddToCart(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// success
	jobUC.On("AddToCart", mock.Anything, "usr_1", "job_1").Return(nil).Once()
	req = httptest.NewRequest("POST", "/cart", strings.NewReader(`{"job_id":"job_1"}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.AddToCart(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// error
	jobUC.On("AddToCart", mock.Anything, "usr_1", "job_1").Return(errors.New("err")).Once()
	req = httptest.NewRequest("POST", "/cart", strings.NewReader(`{"job_id":"job_1"}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.AddToCart(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJobHandler_ListCart(t *testing.T) {
	jobUC := new(MockJobUC)
	h := handler.NewJobHandler(jobUC)

	// unauthorized
	req := httptest.NewRequest("GET", "/cart", nil)
	w := httptest.NewRecorder()
	h.ListCart(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	jobUC.On("ListCart", mock.Anything, "usr_1").Return([]domain.UserCart{}, nil).Once()
	req = httptest.NewRequest("GET", "/cart", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListCart(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	jobUC.On("ListCart", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/cart", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListCart(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJobHandler_RemoveFromCart(t *testing.T) {
	jobUC := new(MockJobUC)
	h := handler.NewJobHandler(jobUC)

	// unauthorized
	req := httptest.NewRequest("DELETE", "/cart/cart_1", nil)
	w := httptest.NewRecorder()
	h.RemoveFromCart(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	jobUC.On("RemoveFromCart", mock.Anything, "usr_1", "cart_1").Return(nil).Once()
	req = httptest.NewRequest("DELETE", "/cart/cart_1", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "cart_1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.RemoveFromCart(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	jobUC.On("RemoveFromCart", mock.Anything, "usr_1", "cart_1").Return(errors.New("err")).Once()
	req = httptest.NewRequest("DELETE", "/cart/cart_1", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.RemoveFromCart(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJobHandler_CreateReview(t *testing.T) {
	jobUC := new(MockJobUC)
	h := handler.NewJobHandler(jobUC)

	// unauthorized
	req := httptest.NewRequest("POST", "/reviews", nil)
	w := httptest.NewRecorder()
	h.CreateReview(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// bad body
	req = httptest.NewRequest("POST", "/reviews", strings.NewReader("bad_json"))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.CreateReview(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// success
	jobUC.On("WriteReview", mock.Anything, "usr_1", "job_1", mock.Anything).Return(nil).Once()
	req = httptest.NewRequest("POST", "/reviews", strings.NewReader(`{"job_id":"job_1","rating_stars":5,"review_text":"good"}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.CreateReview(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// error
	jobUC.On("WriteReview", mock.Anything, "usr_1", "job_1", mock.Anything).Return(errors.New("err")).Once()
	req = httptest.NewRequest("POST", "/reviews", strings.NewReader(`{"job_id":"job_1","rating_stars":5,"review_text":"good"}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.CreateReview(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestJobHandler_GetAllReviews(t *testing.T) {
	h := handler.NewJobHandler(nil)
	req := httptest.NewRequest("GET", "/reviews", nil)
	w := httptest.NewRecorder()
	h.GetAllReviews(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req = httptest.NewRequest("GET", "/reviews?job_id=job_1", nil)
	w = httptest.NewRecorder()
	h.GetAllReviews(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
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
	missionUC.On("ListAvailableMissions", mock.Anything, "usr_1").Return([]domain.UserMission{}, nil).Once()
	req = httptest.NewRequest("GET", "/missions", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListMissions(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	missionUC.On("ListAvailableMissions", mock.Anything, "usr_1").Return(nil, errors.New("err")).Once()
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
	req := httptest.NewRequest("PATCH", "/notifications/n1/read", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "n1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	h.MarkRead(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	notifUC.On("MarkRead", mock.Anything, "n1").Return(errors.New("err")).Once()
	req = httptest.NewRequest("PATCH", "/notifications/n1/read", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.MarkRead(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestNotificationHandler_MarkAllRead(t *testing.T) {
	notifUC := new(MockNotificationUC)
	h := handler.NewNotificationHandler(notifUC)

	// unauthorized
	req := httptest.NewRequest("PATCH", "/notifications/read-all", nil)
	w := httptest.NewRecorder()
	h.MarkAllRead(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	notifUC.On("MarkAllRead", mock.Anything, "usr_1").Return(nil).Once()
	req = httptest.NewRequest("PATCH", "/notifications/read-all", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.MarkAllRead(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// error
	notifUC.On("MarkAllRead", mock.Anything, "usr_1").Return(errors.New("err")).Once()
	req = httptest.NewRequest("PATCH", "/notifications/read-all", nil)
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

func TestRouter_NewRouter(t *testing.T) {
	issuer := new(MockTokenIssuer)
	authH := handler.NewAuthHandler(nil, nil)
	userH := handler.NewUserHandler(nil)
	missionH := handler.NewMissionHandler(nil, nil)
	journeyH := handler.NewJourneyHandler(nil, nil, nil)
	friendH := handler.NewFriendHandler(nil, nil)
	expenseH := handler.NewExpenseHandler(nil)
	notifH := handler.NewNotificationHandler(nil)
	jobH := handler.NewJobHandler(nil)

	r := handler.NewRouter(issuer, authH, userH, missionH, journeyH, friendH, expenseH, notifH, jobH, nil)
	assert.NotNil(t, r)

	// test public route health
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
