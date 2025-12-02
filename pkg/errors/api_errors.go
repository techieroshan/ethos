package errors

import (
	"fmt"
	"net/http"
)

// APIError represents a standardized API error response
type APIError struct {
	Message    string `json:"error"`
	Code       string `json:"code"`
	HTTPStatus int    `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}

// Predefined error types
var (
	ErrInvalidCredentials = &APIError{
		Message:    "Invalid credentials",
		Code:       "AUTH_INVALID_CREDENTIALS",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrTokenExpired = &APIError{
		Message:    "Token has expired",
		Code:       "AUTH_TOKEN_EXPIRED",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrTokenInvalid = &APIError{
		Message:    "Invalid token",
		Code:       "AUTH_TOKEN_INVALID",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrEmailUnverified = &APIError{
		Message:    "Email not verified",
		Code:       "AUTH_EMAIL_UNVERIFIED",
		HTTPStatus: http.StatusForbidden,
	}

	ErrValidationFailed = &APIError{
		Message:    "Validation failed",
		Code:       "VALIDATION_FAILED",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrUserNotFound = &APIError{
		Message:    "User not found",
		Code:       "USER_NOT_FOUND",
		HTTPStatus: http.StatusNotFound,
	}

	ErrProfileNotFound = &APIError{
		Message:    "Profile not found",
		Code:       "PROFILE_NOT_FOUND",
		HTTPStatus: http.StatusNotFound,
	}

	ErrEmailAlreadyExists = &APIError{
		Message:    "Email already exists",
		Code:       "EMAIL_ALREADY_EXISTS",
		HTTPStatus: http.StatusConflict,
	}

	ErrServerError = &APIError{
		Message:    "Internal server error",
		Code:       "SERVER_ERROR",
		HTTPStatus: http.StatusInternalServerError,
	}
)

// NewValidationError creates a validation error with a custom message
func NewValidationError(message string) *APIError {
	return &APIError{
		Message:    message,
		Code:       "VALIDATION_FAILED",
		HTTPStatus: http.StatusBadRequest,
	}
}

// WrapError wraps an error with context
func WrapError(err error, context string) error {
	return fmt.Errorf("%s: %w", context, err)
}

