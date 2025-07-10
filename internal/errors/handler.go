package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type ErrorResponseWrapper struct {
	Error     *AppError `json:"error"`
	RequestID string    `json:"request_id,omitempty"`
	Timestamp string    `json:"timestamp,omitempty"`
}

func HandleError(w http.ResponseWriter, err error) {
	var appErr *AppError

	if errors.As(err, &appErr) && appErr != nil {
	} else {
		appErr = NewInternalError("An unexpected error occurred", err)
	}

	response := ErrorResponseWrapper{
		Error:     appErr,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)
	json.NewEncoder(w).Encode(response)
}
