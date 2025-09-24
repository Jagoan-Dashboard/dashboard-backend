package errors

import (
	"fmt"
	"net/http"
)

// AppError represents application-specific errors
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Error codes
const (
	// Authentication errors
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeInvalidToken     = "INVALID_TOKEN"
	ErrCodeTokenExpired     = "TOKEN_EXPIRED"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"

	// Authorization errors
	ErrCodeForbidden        = "FORBIDDEN"
	ErrCodeInsufficientRole = "INSUFFICIENT_ROLE"

	// Validation errors
	ErrCodeValidationFailed = "VALIDATION_FAILED"
	ErrCodeInvalidInput     = "INVALID_INPUT"
	ErrCodeMissingField     = "MISSING_FIELD"

	// Resource errors
	ErrCodeResourceNotFound = "RESOURCE_NOT_FOUND"
	ErrCodeResourceExists   = "RESOURCE_ALREADY_EXISTS"
	ErrCodeResourceDeleted  = "RESOURCE_DELETED"

	// Database errors
	ErrCodeDatabaseError    = "DATABASE_ERROR"
	ErrCodeConnectionFailed = "CONNECTION_FAILED"
	ErrCodeConstraintViolation = "CONSTRAINT_VIOLATION"

	// File/Upload errors
	ErrCodeFileUploadFailed = "FILE_UPLOAD_FAILED"
	ErrCodeInvalidFileType  = "INVALID_FILE_TYPE"
	ErrCodeFileTooLarge     = "FILE_TOO_LARGE"

	// Cache errors
	ErrCodeCacheError       = "CACHE_ERROR"

	// External service errors
	ErrCodeExternalService  = "EXTERNAL_SERVICE_ERROR"

	// Internal errors
	ErrCodeInternalError    = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

// Predefined errors
var (
	ErrUnauthorized = &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    "Authentication required",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrForbidden = &AppError{
		Code:       ErrCodeForbidden,
		Message:    "Access forbidden",
		HTTPStatus: http.StatusForbidden,
	}

	ErrInvalidInput = &AppError{
		Code:       ErrCodeInvalidInput,
		Message:    "Invalid input provided",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrResourceNotFound = &AppError{
		Code:       ErrCodeResourceNotFound,
		Message:    "Resource not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrInternalError = &AppError{
		Code:       ErrCodeInternalError,
		Message:    "Internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}
)

// New creates a new AppError
func New(code, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// NewWithDetails creates a new AppError with details
func NewWithDetails(code, message, details string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		HTTPStatus: httpStatus,
	}
}

// Wrap wraps an existing error with application context
func Wrap(err error, code, message string, httpStatus int) *AppError {
	details := ""
	if err != nil {
		details = err.Error()
	}

	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		HTTPStatus: httpStatus,
	}
}

// Validation errors
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeValidationFailed,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewValidationErrorWithDetails(field, message string) *AppError {
	return &AppError{
		Code:       ErrCodeValidationFailed,
		Message:    "Validation failed",
		Details:    fmt.Sprintf("%s: %s", field, message),
		HTTPStatus: http.StatusBadRequest,
	}
}

// Resource errors
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeResourceNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		HTTPStatus: http.StatusNotFound,
	}
}

func NewAlreadyExistsError(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeResourceExists,
		Message:    fmt.Sprintf("%s already exists", resource),
		HTTPStatus: http.StatusConflict,
	}
}

// Database errors
func NewDatabaseError(operation string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeDatabaseError,
		Message:    fmt.Sprintf("Database %s failed", operation),
		Details:    err.Error(),
		HTTPStatus: http.StatusInternalServerError,
	}
}

// IsAppError checks if error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError safely converts error to AppError
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	// Return generic internal error for non-AppError types
	return &AppError{
		Code:       ErrCodeInternalError,
		Message:    "Internal server error",
		Details:    err.Error(),
		HTTPStatus: http.StatusInternalServerError,
	}
}