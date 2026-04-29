package errors

import (
	"fmt"
	"net/http"
)

type ErrorType string
const (
	ErrorTypeInternal     ErrorType = "INTERNAL"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeBadRequest   ErrorType = "BAD_REQUEST"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeTimeout      ErrorType = "TIMEOUT"
	ErrorTypeUnavailable  ErrorType = "UNAVAILABLE"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	Type        ErrorType         `json:"type"`
	Code        int               `json:"-"`
	Message     string            `json:"message"`
	Retryable   bool              `json:"retryable,omitempty"`
	Validations []ValidationError `json:"validations,omitempty"`
	Internal    error             `json:"-"`
}


func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Internal)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Internal
}

func (e *AppError) WithInternal(err error) *AppError {
	copy := *e
	copy.Internal = err
	return &copy
}

func (e *AppError) WithMessage(message string) *AppError {
	copy := *e
	copy.Message = message
	return &copy
}

func (e *AppError) WithValidations(validations []ValidationError) *AppError {
	copy := *e
	copy.Validations = validations
	return &copy
}

func (e *AppError) AsRetryable() *AppError {
	copy := *e
	copy.Retryable = true
	return &copy
}

func NewValidationError(validations []ValidationError) *AppError {
	return ErrValidationFailed.WithValidations(validations)
}


var (
	ErrBadRequest = &AppError{
		Type:    ErrorTypeBadRequest,
		Code:    http.StatusBadRequest,
		Message: "Bad request",
	}

	ErrValidationFailed = &AppError{
		Type:    ErrorTypeBadRequest,
		Code:    http.StatusBadRequest,
		Message: "Validation failed",
	}

	ErrUnauthorized = &AppError{
		Type:    ErrorTypeUnauthorized,
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized",
	}

	ErrForbidden = &AppError{
		Type:    ErrorTypeForbidden,
		Code:    http.StatusForbidden,
		Message: "Forbidden",
	}

	ErrNotFound = &AppError{
		Type:    ErrorTypeNotFound,
		Code:    http.StatusNotFound,
		Message: "Resource not found",
	}

	ErrConflict = &AppError{
		Type:    ErrorTypeConflict,
		Code:    http.StatusConflict,
		Message: "Resource conflict",
	}

	ErrTooManyRequests = &AppError{
		Type:      ErrorTypeBadRequest, // Or a dedicated RATE_LIMIT type
		Code:      http.StatusTooManyRequests,
		Message:   "Too many requests",
		Retryable: true,
	}

	ErrInternal = &AppError{
		Type:    ErrorTypeInternal,
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	}

	ErrServiceUnavailable = &AppError{
		Type:      ErrorTypeUnavailable,
		Code:      http.StatusServiceUnavailable,
		Message:   "Service unavailable",
		Retryable: true,
	}

	ErrTimeout = &AppError{
		Type:      ErrorTypeTimeout,
		Code:      http.StatusGatewayTimeout,
		Message:   "Request timeout",
		Retryable: true,
	}
)

