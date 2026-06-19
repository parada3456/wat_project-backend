package apperror

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	log.Println("debugprint: entering (*AppError).Error")
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func RespondError(w http.ResponseWriter, err error) {
	log.Println("debugprint: entering RespondError")
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = &AppError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Err:     err,
		}
	}

	if appErr.Code >= 500 {
		log.Printf("Internal Error: %v", appErr.Err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	json.NewEncoder(w).Encode(appErr)
}
