package errors

import (
	"database/sql"
	"strings"
)

type DatabaseErrorHandler struct{}

func NewDatabaseErrorHandler() *DatabaseErrorHandler {
	return &DatabaseErrorHandler{}
}

func (d *DatabaseErrorHandler) HandleDatabaseError(operation string, err error) *AppError {
	switch {
	case err == sql.ErrNoRows:
		return NewNotFoundError("Resource not found", err)

	case strings.Contains(err.Error(), "connection"):
		return NewDatabaseError("Service temporarily unavailable", nil)

	case strings.Contains(err.Error(), "timeout"):
		return NewDatabaseError("Request timeout", nil)

	default:
		return NewDatabaseError("Database operation failed", nil)
	}
}
