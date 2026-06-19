package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/pkg/apperror"
)

type JobUC interface {
	ListJobs(ctx context.Context, filters map[string]interface{}) ([]domain.JobPosting, error)
	GetJobDetail(ctx context.Context, id string) (*domain.JobPosting, []domain.JobHousing, *domain.JobOverallRating, error)
	AddToCart(ctx context.Context, userID, jobID string) error
	ListCart(ctx context.Context, userID string) ([]domain.UserCart, error)
	RemoveFromCart(ctx context.Context, userID, id string) error
	WriteReview(ctx context.Context, userID, jobID string, rev *domain.JobReview) error
}

type JobHandler struct {
	jobUC JobUC
}

func NewJobHandler(jobUC JobUC) *JobHandler {
	log.Println("debugprint: entering NewJobHandler")
	return &JobHandler{jobUC: jobUC}
}

func (h *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	log.
		// Simple filters from query params
		Println("debugprint: entering (*JobHandler).ListJobs")

	filters := make(map[string]interface{})
	if agency := r.URL.Query().Get("agency"); agency != "" {
		filters["agency_name"] = agency
	}

	jobs, err := h.jobUC.ListJobs(r.Context(), filters)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	json.NewEncoder(w).Encode(jobs)
}

func (h *JobHandler) GetJobDetail(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).GetJobDetail")
	id := chi.URLParam(r, "id")
	job, housing, rating, err := h.jobUC.GetJobDetail(r.Context(), id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"job":     job,
		"housing": housing,
		"rating":  rating,
	})
}

type cartReq struct {
	JobID string `json:"job_id" validate:"required"`
}

func (h *JobHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).AddToCart")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var req cartReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Invalid request body"})
		return
	}

	err := h.jobUC.AddToCart(r.Context(), claims.UserID, req.JobID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *JobHandler) ListCart(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).ListCart")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	cart, err := h.jobUC.ListCart(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	json.NewEncoder(w).Encode(cart)
}

func (h *JobHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).RemoveFromCart")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	err := h.jobUC.RemoveFromCart(r.Context(), claims.UserID, id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JobHandler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	log.
		// Simplified: list reviews for a job if job_id is provided, else all
		Println("debugprint: entering (*JobHandler).GetAllReviews")

	jobID := r.URL.Query().Get("job_id")
	if jobID == "" {
		// Mock implementation from before
		return
	}
	// ... implementation
}

func (h *JobHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).CreateReview")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var req struct {
		JobID       string  `json:"job_id"`
		RatingStars float64 `json:"rating_stars"`
		ReviewText  string  `json:"review_text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Invalid request body"})
		return
	}

	err := h.jobUC.WriteReview(r.Context(), claims.UserID, req.JobID, &domain.JobReview{
		RatingStars: req.RatingStars,
		ReviewText:  req.ReviewText,
	})
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
