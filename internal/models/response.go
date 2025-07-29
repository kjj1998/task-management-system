package models

import "time"

// BaseResponse is the standard response wrapper for all API endpoints
type BaseResponse struct {
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
	Data      any        `json:"data,omitempty"`
	Error     *ErrorInfo `json:"error,omitempty"`
	Meta      *Meta      `json:"meta,omitempty"`
	Timestamp string     `json:"timestamp"`
	RequestID string     `json:"request_id,omitempty"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta contains additional response metadata
type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// SuccessResponse creates a standardized success response
func NewSuccessResponse(message string, data any) *BaseResponse {
	return &BaseResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// ErrorResponse creates a standardized error response
func NewErrorResponse(message string, errorInfo *ErrorInfo) *BaseResponse {
	return &BaseResponse{
		Success:   false,
		Message:   message,
		Error:     errorInfo,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// PaginatedResponse creates a response with pagination metadata
func NewPaginatedResponse(message string, data any, meta *Meta) *BaseResponse {
	return &BaseResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
