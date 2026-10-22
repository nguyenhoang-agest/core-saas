package exceptions

import (
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions/err_type"
	"github.com/gofiber/fiber/v3"
)

var _ CustomError = (*UnauthorizeError)(nil)
var _ error = (*UnauthorizeError)(nil)

type UnauthorizeError struct {
	ResponseError
}

func NewUnauthorizeError(responseCode response_code.UnauthorizedCode, errorDetail []FieldErrorDetail) UnauthorizeError {
	return UnauthorizeError{
		ResponseError: NewResponseError(responseCode, err_type.UnauthorizedError, errorDetail),
	}
}

func (b UnauthorizeError) StatusCode() int {
	return fiber.StatusUnauthorized
}

func (b UnauthorizeError) Error() string {
	return b.MessageCode
}

func (b UnauthorizeError) Json() map[string]any {
	return map[string]any{
		"status":       b.Status,
		"message_code": b.MessageCode,
		"error_detail": b.ErrorDetail,
	}
}

func (b UnauthorizeError) FormattedLog() string {
	return "UnauthorizeError: " + b.MessageCode
}
