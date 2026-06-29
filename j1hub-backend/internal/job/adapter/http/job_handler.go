package http

import (
	"github.com/j1hub/backend/internal/job/adapter/http/dto"

	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/transport/http/middleware"
	jobdomain "github.com/j1hub/backend/internal/job/domain"
	"github.com/j1hub/backend/pkg/apperror"
)

type JobUC interface {
	ListJobs(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]jobdomain.JobPosting, int, error)
	GetJobDetail(ctx context.Context, id string) (*jobdomain.JobPosting, []jobdomain.JobHousing, *jobdomain.JobOverallRating, error)
	AddToCart(ctx context.Context, userID, jobID string) error
	ListCart(ctx context.Context, userID string) ([]jobdomain.UserCart, error)
	RemoveFromCart(ctx context.Context, userID, id string) error
	WriteReview(ctx context.Context, userID, jobID string, rev *jobdomain.JobReview) error
	ListReviews(ctx context.Context, jobID string) ([]jobdomain.JobReview, error)
	UpdateCartStatus(ctx context.Context, userID, cartID string, status jobdomain.CartStatus) error
	CreateJob(ctx context.Context, job *jobdomain.JobPosting) error
	UpdateJob(ctx context.Context, job *jobdomain.JobPosting) error
	PatchJob(ctx context.Context, id string, updates map[string]interface{}) error
	DeleteJob(ctx context.Context, id string) error
}

type JobHandler struct {
	jobUC JobUC
}

func NewJobHandler(jobUC JobUC) *JobHandler {
	log.Println("debugprint: entering NewJobHandler")
	return &JobHandler{jobUC: jobUC}
}

func (h *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).ListJobs")

	filters := make(map[string]interface{})
	if agency := r.URL.Query().Get("agency"); agency != "" {
		filters["agency_name"] = agency
	}

	pago := apperror.ParsePagination(r)
	jobs, totalCount, err := h.jobUC.ListJobs(r.Context(), filters, pago.Page, pago.PageSize)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	apperror.RespondList(w, jobs, pago.Page, pago.PageSize, totalCount)
}

func (h *JobHandler) GetJobDetail(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).GetJobDetail")
	id := chi.URLParam(r, "id")
	job, housing, rating, err := h.jobUC.GetJobDetail(r.Context(), id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respDTO := dto.NewJobDetailResponse(job, housing, rating)
	json.NewEncoder(w).Encode(respDTO)
}

func (h *JobHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).AddToCart")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	var req dto.CartReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
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
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	cart, err := h.jobUC.ListCart(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	pago := apperror.ParsePagination(r)
	apperror.RespondList(w, cart, pago.Page, pago.PageSize, len(cart))
}

func (h *JobHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).RemoveFromCart")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "cartId")
	if id == "" {
		id = chi.URLParam(r, "id")
	}
	err := h.jobUC.RemoveFromCart(r.Context(), claims.UserID, id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JobHandler) UpdateCartStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).UpdateCartStatus")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	cartID := chi.URLParam(r, "cartId")
	if cartID == "" {
		cartID = chi.URLParam(r, "id")
	}
	var req struct {
		Status jobdomain.CartStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	err := h.jobUC.UpdateCartStatus(r.Context(), claims.UserID, cartID, req.Status)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JobHandler) GetJobReviews(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).GetJobReviews")
	id := chi.URLParam(r, "id")
	reviews, err := h.jobUC.ListReviews(r.Context(), id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	pago := apperror.ParsePagination(r)
	apperror.RespondList(w, reviews, pago.Page, pago.PageSize, len(reviews))
}

func (h *JobHandler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).GetAllReviews")

	jobID := r.URL.Query().Get("job_id")
	if jobID == "" {
		apperror.RespondError(w, fmt.Errorf("Missing job_id query parameter: %w", domain.ErrInvalidInput))
		return
	}

	reviews, err := h.jobUC.ListReviews(r.Context(), jobID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	pago := apperror.ParsePagination(r)
	apperror.RespondList(w, reviews, pago.Page, pago.PageSize, len(reviews))
}

func (h *JobHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).CreateReview")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	var req struct {
		JobID       string  `json:"job_id"`
		RatingStars float64 `json:"rating_stars"`
		ReviewText  string  `json:"review_text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	jobID := chi.URLParam(r, "id")
	if jobID == "" {
		jobID = req.JobID
	}

	err := h.jobUC.WriteReview(r.Context(), claims.UserID, jobID, &jobdomain.JobReview{
		RatingStars: req.RatingStars,
		ReviewText:  req.ReviewText,
	})
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}


func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).CreateJob")
	var req jobdomain.JobPosting
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	err := h.jobUC.CreateJob(r.Context(), &req)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).UpdateJob")
	id := chi.URLParam(r, "id")
	var req jobdomain.JobPosting
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}
	req.JobID = id

	err := h.jobUC.UpdateJob(r.Context(), &req)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

func (h *JobHandler) PatchJob(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).PatchJob")
	id := chi.URLParam(r, "id")
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	err := h.jobUC.PatchJob(r.Context(), id, updates)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*JobHandler).DeleteJob")
	id := chi.URLParam(r, "id")
	err := h.jobUC.DeleteJob(r.Context(), id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
