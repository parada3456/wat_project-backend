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
	jobdomain "github.com/j1hub/backend/internal/job/domain"
	"github.com/j1hub/backend/internal/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJobHandler_ListJobs(t *testing.T) {
	jobUC := new(MockJobUC)
	h := handler.NewJobHandler(jobUC)

	// success
	jobUC.On("ListJobs", mock.Anything, map[string]interface{}{"agency_name": "agency_1"}).Return([]jobdomain.JobPosting{}, nil).Once()
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
	jobUC.On("GetJobDetail", mock.Anything, "job_1").Return(&jobdomain.JobPosting{}, []jobdomain.JobHousing{}, &jobdomain.JobOverallRating{}, nil).Once()
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
	jobUC.On("ListCart", mock.Anything, "usr_1").Return([]jobdomain.UserCart{}, nil).Once()
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
	jobUC := new(MockJobUC)
	h := handler.NewJobHandler(jobUC)

	// missing job_id
	req := httptest.NewRequest("GET", "/reviews", nil)
	w := httptest.NewRecorder()
	h.GetAllReviews(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// success
	jobUC.On("ListReviews", mock.Anything, "job_1").Return([]jobdomain.JobReview{}, nil).Once()
	req = httptest.NewRequest("GET", "/reviews?job_id=job_1", nil)
	w = httptest.NewRecorder()
	h.GetAllReviews(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJobHandler_UpdateCartStatus(t *testing.T) {
	jobUC := new(MockJobUC)
	h := handler.NewJobHandler(jobUC)

	// unauthorized
	req := httptest.NewRequest("PATCH", "/cart/cart_1", strings.NewReader(`{"status":"Applied"}`))
	w := httptest.NewRecorder()
	h.UpdateCartStatus(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// success
	jobUC.On("UpdateCartStatus", mock.Anything, "usr_1", "cart_1", jobdomain.CartStatus("Applied")).Return(nil).Once()
	req = httptest.NewRequest("PATCH", "/cart/cart_1", strings.NewReader(`{"status":"Applied"}`))
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &port.Claims{UserID: "usr_1"}))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cartId", "cart_1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	h.UpdateCartStatus(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestJobHandler_GetJobReviews(t *testing.T) {
	jobUC := new(MockJobUC)
	h := handler.NewJobHandler(jobUC)

	jobUC.On("ListReviews", mock.Anything, "job_1").Return([]jobdomain.JobReview{}, nil).Once()
	req := httptest.NewRequest("GET", "/jobs/job_1/reviews", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "job_1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	h.GetJobReviews(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
