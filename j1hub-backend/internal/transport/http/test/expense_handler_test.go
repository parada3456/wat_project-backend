package test

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
	expensehandler "github.com/j1hub/backend/internal/expense/adapter/http"
	"github.com/j1hub/backend/internal/transport/http/middleware"
	expensedomain "github.com/j1hub/backend/internal/expense/domain"
	port "github.com/j1hub/backend/internal/auth/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExpenseHandler_ListExpenses(t *testing.T) {
	expenseUC := new(MockManageExpenseUC)
	h := expensehandler.NewExpenseHandler(expenseUC)

	// unauthorized
	req := httptest.NewRequest("GET", "/expenses", nil)
	w := httptest.NewRecorder()
	h.ListExpenses(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	expenseUC.On("ListExpenses", mock.Anything, "usr_1", 1, 10).Return([]expensedomain.ExpenseTransaction{}, 0, nil).Once()
	req = httptest.NewRequest("GET", "/expenses", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListExpenses(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	expenseUC.On("ListExpenses", mock.Anything, "usr_1", 1, 10).Return(nil, 0, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/expenses", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListExpenses(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestExpenseHandler_CreateExpense(t *testing.T) {
	expenseUC := new(MockManageExpenseUC)
	h := expensehandler.NewExpenseHandler(expenseUC)

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
	h := expensehandler.NewExpenseHandler(expenseUC)

	// unauthorized
	req := httptest.NewRequest("GET", "/expenses/exp_1", nil)
	w := httptest.NewRecorder()
	h.GetExpenseDetail(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	expenseUC.On("GetExpenseDetail", mock.Anything, "usr_1", "exp_1").Return(&expensedomain.ExpenseTransaction{}, []expensedomain.ExpenseSplit{}, nil).Once()
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
	h := expensehandler.NewExpenseHandler(expenseUC)

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
	h := expensehandler.NewExpenseHandler(expenseUC)

	// unauthorized
	req := httptest.NewRequest("GET", "/expenses/pending", nil)
	w := httptest.NewRecorder()
	h.ListPending(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	expenseUC.On("ListPendingExpenses", mock.Anything, "usr_1", 1, 10).Return([]expensedomain.ExpenseSplit{}, 0, nil).Once()
	req = httptest.NewRequest("GET", "/expenses/pending", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListPending(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// error
	expenseUC.On("ListPendingExpenses", mock.Anything, "usr_1", 1, 10).Return(nil, 0, errors.New("err")).Once()
	req = httptest.NewRequest("GET", "/expenses/pending", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	w = httptest.NewRecorder()
	h.ListPending(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestExpenseHandler_PaySplit(t *testing.T) {
	expenseUC := new(MockManageExpenseUC)
	h := expensehandler.NewExpenseHandler(expenseUC)

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
	h := expensehandler.NewExpenseHandler(expenseUC)

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
