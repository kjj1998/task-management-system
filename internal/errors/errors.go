package errors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeDatabase ErrorType = "DATABASE_ERROR"
	ErrorTypeInternal ErrorType = "INTERNAL_ERROR"
	ErrorTypeNotFound ErrorType = "NOT_FOUND"
)

type AppError struct {
	Type       ErrorType `json:"type"`
	Message    string    `json:"message"`
	Code       string    `json:"code,omitempty"`
	Details    string    `json:"details,omitempty"`
	StatusCode int       `json:"-"`
	Err        error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s, %v", e.Message, e.Err)
	}
	return e.Message
}

func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Type:       ErrorTypeDatabase,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Type:       ErrorTypeNotFound,
		Message:    message,
		StatusCode: http.StatusNotFound,
		Err:        err,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:       ErrorTypeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}
