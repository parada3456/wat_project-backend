package http

import (
	"github.com/j1hub/backend/internal/expense/adapter/http/dto"

	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/j1hub/backend/internal/domain"
	"github.com/j1hub/backend/internal/transport/http/middleware"
	expensedomain "github.com/j1hub/backend/internal/expense/domain"
	expenseusecase "github.com/j1hub/backend/internal/expense/usecase"
	"github.com/j1hub/backend/pkg/apperror"
)

type ManageExpenseUC interface {
	ListExpenses(ctx context.Context, userID string, page, pageSize int) ([]expensedomain.ExpenseTransaction, int, error)
	CreateExpense(ctx context.Context, userID string, cmd expenseusecase.CreateExpenseCmd) error
	GetExpenseDetail(ctx context.Context, userID string, id string) (*expensedomain.ExpenseTransaction, []expensedomain.ExpenseSplit, error)
	DeleteExpense(ctx context.Context, userID string, id string) error
	ListPendingExpenses(ctx context.Context, userID string, page, pageSize int) ([]expensedomain.ExpenseSplit, int, error)
	SubmitSlip(ctx context.Context, debtorID, splitID string, file io.Reader, contentType string) error
	ApproveSplit(ctx context.Context, userID string, id string) error
}

type ExpenseHandler struct {
	expenseUC ManageExpenseUC
	validate  *validator.Validate
}

func NewExpenseHandler(expenseUC ManageExpenseUC) *ExpenseHandler {
	log.Println("debugprint: entering NewExpenseHandler")
	return &ExpenseHandler{expenseUC: expenseUC, validate: validator.New()}
}

func (h *ExpenseHandler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*ExpenseHandler).ListExpenses")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	pago := apperror.ParsePagination(r)
	expenses, totalCount, err := h.expenseUC.ListExpenses(r.Context(), claims.UserID, pago.Page, pago.PageSize)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	apperror.RespondList(w, expenses, pago.Page, pago.PageSize, totalCount)
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*ExpenseHandler).CreateExpense")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	var req dto.CreateExpenseReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Invalid request body: %w", domain.ErrInvalidInput))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		apperror.RespondError(w, fmt.Errorf("Validation failed: %w", domain.ErrInvalidInput))
		return
	}

	cmd := expenseusecase.CreateExpenseCmd{
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
	log.Println("debugprint: entering (*ExpenseHandler).GetExpenseDetail")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	txn, splits, err := h.expenseUC.GetExpenseDetail(r.Context(), claims.UserID, id)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respDTO := dto.NewExpenseDetailResponse(txn, splits)
	json.NewEncoder(w).Encode(respDTO)
}

func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*ExpenseHandler).DeleteExpense")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
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
	log.Println("debugprint: entering (*ExpenseHandler).ListPending")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	pago := apperror.ParsePagination(r)
	pending, totalCount, err := h.expenseUC.ListPendingExpenses(r.Context(), claims.UserID, pago.Page, pago.PageSize)
	if err != nil {
		apperror.RespondError(w, err)
		return
	}
	apperror.RespondList(w, pending, pago.Page, pago.PageSize, totalCount)
}

func (h *ExpenseHandler) UpdateSplit(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*ExpenseHandler).UpdateSplit")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	splitID := chi.URLParam(r, "splitId")
	if splitID == "" {
		splitID = chi.URLParam(r, "id")
	}

	contentType := r.Header.Get("Content-Type")
	isMultipart := false
	if len(contentType) >= 19 && contentType[:19] == "multipart/form-data" {
		isMultipart = true
	}

	if isMultipart {
		r.ParseMultipartForm(10 << 20)
		file, header, err := r.FormFile("slip")
		if err != nil {
			apperror.RespondError(w, fmt.Errorf("Slip file required: %w", domain.ErrInvalidInput))
			return
		}
		defer file.Close()

		err = h.expenseUC.SubmitSlip(r.Context(), claims.UserID, splitID, file, header.Header.Get("Content-Type"))
		if err != nil {
			apperror.RespondError(w, err)
			return
		}
	} else {
		err := h.expenseUC.ApproveSplit(r.Context(), claims.UserID, splitID)
		if err != nil {
			apperror.RespondError(w, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ExpenseHandler) PaySplit(w http.ResponseWriter, r *http.Request) {
	log.Println("debugprint: entering (*ExpenseHandler).PaySplit")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("slip")
	if err != nil {
		apperror.RespondError(w, fmt.Errorf("Slip file required: %w", domain.ErrInvalidInput))
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
	log.Println("debugprint: entering (*ExpenseHandler).ApproveSplit")
	claims := middleware.GetClaims(r.Context())
	if claims == nil {
		apperror.RespondError(w, domain.ErrUnauthorized)
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
