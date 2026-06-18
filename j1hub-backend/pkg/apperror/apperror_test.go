package apperror_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/j1hub/backend/pkg/apperror"
	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	err := &apperror.AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid input",
		Err:     errors.New("raw error"),
	}
	assert.Equal(t, "Invalid input: raw error", err.Error())

	errNoRaw := &apperror.AppError{
		Code:    http.StatusNotFound,
		Message: "Not found",
	}
	assert.Equal(t, "Not found", errNoRaw.Error())
}

func TestRespondError_AppError(t *testing.T) {
	err := &apperror.AppError{
		Code:    http.StatusNotFound,
		Message: "Resource not found",
	}

	w := httptest.NewRecorder()
	apperror.RespondError(w, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, "Resource not found", body["message"])
}

func TestRespondError_InternalError(t *testing.T) {
	rawErr := errors.New("database connection lost")

	w := httptest.NewRecorder()
	apperror.RespondError(w, rawErr)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, "Internal server error", body["message"])
}
