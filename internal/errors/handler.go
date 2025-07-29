package errors

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kjj1998/task-management-system/internal/models"
)

func HandleError(w http.ResponseWriter, err error, logger *slog.Logger) {
	var appErr *AppError

	if !errors.As(err, &appErr) {
		appErr = NewInternalError("An unexpected error occurred", err)
	}

	if appErr.StatusCode >= 500 {
		logger.Error("server error occurred",
			slog.String("error", appErr.Err.Error()),
			slog.String("details", appErr.Details),
			slog.Int("status_code", appErr.StatusCode),
			slog.String("error_code", appErr.Code),
		)
	} else {
		logger.Warn("client error occurred",
			slog.String("error", appErr.Err.Error()),
			slog.Int("status_code", appErr.StatusCode),
			slog.String("error_code", appErr.Code),
		)
	}

	errorInfo := &models.ErrorInfo{
		Code:    appErr.Code,
		Message: appErr.Message,
		Details: appErr.Details,
	}

	response := models.NewErrorResponse("An error occurred", errorInfo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		logger.Error("failed to encode error response", slog.String("encode_error",
			encodeErr.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
