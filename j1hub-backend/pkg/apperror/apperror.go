package apperror

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/j1hub/backend/internal/domain"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

type PaginationMetadata struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type PagedResponse struct {
	Data       interface{}        `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func (p *ProblemDetails) Error() string {
	if p.Detail != "" {
		return p.Title + ": " + p.Detail
	}
	return p.Title
}

func (e *AppError) Error() string {
	log.Println("debugprint: entering (*AppError).Error")
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

type errorDetails struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

type errorResponse struct {
	Error errorDetails `json:"error"`
}

func RespondError(w http.ResponseWriter, err error) {
	log.Println("debugprint: entering RespondError")

	var probDetails *ProblemDetails
	if errors.As(err, &probDetails) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(probDetails.Status)
		json.NewEncoder(w).Encode(probDetails)
		return
	}

	var appErr *AppError
	if !errors.As(err, &appErr) {
		if errors.Is(err, domain.ErrNotFound) {
			appErr = &AppError{Code: http.StatusNotFound, Message: err.Error(), Err: err}
		} else if errors.Is(err, domain.ErrUnauthorized) || errors.Is(err, domain.ErrInvalidCredentials) {
			appErr = &AppError{Code: http.StatusUnauthorized, Message: err.Error(), Err: err}
		} else if errors.Is(err, domain.ErrForbidden) {
			appErr = &AppError{Code: http.StatusForbidden, Message: err.Error(), Err: err}
		} else if errors.Is(err, domain.ErrConflict) ||
			errors.Is(err, domain.ErrAlreadyCompleted) ||
			errors.Is(err, domain.ErrDuplicateFriend) ||
			errors.Is(err, domain.ErrProofAlreadySubmitted) {
			appErr = &AppError{Code: http.StatusConflict, Message: err.Error(), Err: err}
		} else if errors.Is(err, domain.ErrInvalidInput) ||
			errors.Is(err, domain.ErrSelfSplit) ||
			errors.Is(err, domain.ErrPhaseNotComplete) {
			appErr = &AppError{Code: http.StatusBadRequest, Message: err.Error(), Err: err}
		} else {
			appErr = &AppError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
				Err:     err,
			}
		}
	}

	if appErr.Code >= 500 {
		log.Printf("Internal Error: %v", appErr.Err)
	}

	errorCode := "INTERNAL_SERVER_ERROR"
	switch appErr.Code {
	case http.StatusBadRequest:
		errorCode = "BAD_REQUEST"
	case http.StatusUnauthorized:
		errorCode = "UNAUTHORIZED"
	case http.StatusForbidden:
		errorCode = "FORBIDDEN"
	case http.StatusNotFound:
		errorCode = "RESOURCE_NOT_FOUND"
	case http.StatusConflict:
		errorCode = "CONFLICT"
	}

	resp := errorResponse{
		Error: errorDetails{
			Code:    errorCode,
			Message: appErr.Message,
			Details: nil,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	json.NewEncoder(w).Encode(resp)
}

func RespondList(w http.ResponseWriter, items interface{}, page, pageSize, totalItems int) {
	// 1. Safeguard items so it never encodes as JSON "null"
	if items == nil {
		items = []interface{}{}
	} else {
		v := reflect.ValueOf(items)
		if (v.Kind() == reflect.Slice || v.Kind() == reflect.Array) && v.IsNil() {
			items = []interface{}{}
		}
	}

	// 2. Calculate the slice length for fallback logic
	var sliceLen int
	v := reflect.ValueOf(items)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		sliceLen = v.Len()
	}

	// 3. Fallbacks and defaults
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = sliceLen
		if pageSize <= 0 {
			pageSize = 1
		}
	}
	if totalItems <= 0 {
		totalItems = sliceLen
	}

	totalPages := (totalItems + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	// 4. Construct the structured response
	response := PagedResponse{
		Data: items,
		Pagination: PaginationMetadata{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}

	// 5. Send JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response) // Look ma, no dynamic maps!
}

// PaginationRequest binds and validates incoming query parameters
type PaginationRequest struct {
	Page     int
	PageSize int
}

// Offset calculates the SQL skip value: (page - 1) * pageSize
func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the page size, perfect for SQL LIMIT clauses
func (p PaginationRequest) Limit() int {
	return p.PageSize
}

func ParsePagination(r *http.Request) PaginationRequest {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	// Apply safe defaults
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return PaginationRequest{
		Page:     page,
		PageSize: pageSize,
	}
}
