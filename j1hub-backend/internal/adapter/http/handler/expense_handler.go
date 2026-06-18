package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/j1hub/backend/internal/adapter/http/middleware"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/usecase"
	"github.com/j1hub/backend/pkg/apperror"
)

type ManageExpenseUC interface {
	ListExpenses(ctx context.Context, userID string) ([]domain.ExpenseTransaction, error)
	CreateExpense(ctx context.Context, userID string, cmd usecase.CreateExpenseCmd) error
	GetExpenseDetail(ctx context.Context, userID string, id string) (*domain.ExpenseTransaction, []domain.ExpenseSplit, error)
	DeleteExpense(ctx context.Context, userID string, id string) error
	ListPendingExpenses(ctx context.Context, userID string) ([]domain.ExpenseSplit, error)
	SubmitSlip(ctx context.Context, debtorID, splitID string, file io.Reader, contentType string) error
	ApproveSplit(ctx context.Context, userID string, id string) error
}

type ExpenseHandler struct {
	expenseUC ManageExpenseUC
	validate  *validator.Validate
}

func NewExpenseHandler(expenseUC ManageExpenseUC) *ExpenseHandler {
	return &ExpenseHandler{expenseUC: expenseUC, validate: validator.New()}
}

func (h *ExpenseHandler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	expenses, err := h.expenseUC.ListExpenses(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	json.NewEncoder(w).Encode(expenses)
}

type createExpenseReq struct {
	Title       string    `json:"title" validate:"required"`
	TotalAmount float64   `json:"total_amount" validate:"required,gt=0"`
	Currency    string    `json:"currency" validate:"required"`
	Memo        string    `json:"memo"`
	DueDate     time.Time `json:"due_date" validate:"required"`
	Splits      []struct {
		UserID    string  `json:"user_id" validate:"required"`
		OweAmount float64 `json:"owe_amount" validate:"required,gt=0"`
	} `json:"splits" validate:"required,dive"`
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	var req createExpenseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Invalid request body", Err: err})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Validation failed", Err: err})
		return
	}

	cmd := usecase.CreateExpenseCmd{
		Title:       req.Title,
		TotalAmount: req.TotalAmount,
		Currency:    req.Currency,
		Memo:        req.Memo,
		DueDate:     req.DueDate,
	}
	for _, s := range req.Splits {
		cmd.Splits = append(cmd.Splits, struct {
			UserID    string
			OweAmount float64
		}{UserID: s.UserID, OweAmount: s.OweAmount})
	}

	err := h.expenseUC.CreateExpense(r.Context(), claims.UserID, cmd)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ExpenseHandler) GetExpenseDetail(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	txn, splits, err := h.expenseUC.GetExpenseDetail(r.Context(), claims.UserID, id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"transaction": txn,
		"splits":      splits,
	})
}

func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	err := h.expenseUC.DeleteExpense(r.Context(), claims.UserID, id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpenseHandler) ListPending(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	pending, err := h.expenseUC.ListPendingExpenses(r.Context(), claims.UserID)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	json.NewEncoder(w).Encode(pending)
}

func (h *ExpenseHandler) PaySplit(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	
	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("slip")
	if err != nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusBadRequest, Message: "Slip file required", Err: err})
		return
	}
	defer file.Close()

	err = h.expenseUC.SubmitSlip(r.Context(), claims.UserID, id, file, header.Header.Get("Content-Type"))
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpenseHandler) ApproveSplit(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, &apperror.AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
		return
	}

	id := chi.URLParam(r, "id")
	err := h.expenseUC.ApproveSplit(r.Context(), claims.UserID, id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
