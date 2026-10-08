package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Code represents a machine-readable error category.
type Code string

const (
	CodeNotFound       Code = "NOT_FOUND"
	CodeInvalidArgument Code = "INVALID_ARGUMENT"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeForbidden       Code = "FORBIDDEN"
	CodeConflict        Code = "CONFLICT"
	CodeInternal        Code = "INTERNAL"
	CodeUnprocessable   Code = "UNPROCESSABLE"
)

// AppError is the standard application error carrying an operational code,
// user-friendly message, optional details, and underlying error cause.
type AppError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatusCode maps an AppError code to a standard HTTP status code.
func (e *AppError) HTTPStatusCode() int {
	switch e.Code {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeInvalidArgument:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeConflict:
		return http.StatusConflict
	case CodeUnprocessable:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

// Helper constructors

func New(code Code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NotFound(message string, err error) *AppError {
	return New(CodeNotFound, message, err)
}

func InvalidArgument(message string, err error) *AppError {
	return New(CodeInvalidArgument, message, err)
}

func Unauthorized(message string, err error) *AppError {
	return New(CodeUnauthorized, message, err)
}

func Forbidden(message string, err error) *AppError {
	return New(CodeForbidden, message, err)
}

func Conflict(message string, err error) *AppError {
	return New(CodeConflict, message, err)
}

func Internal(message string, err error) *AppError {
	return New(CodeInternal, message, err)
}

// AsAppError extracts an *AppError from an error chain if present.
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
