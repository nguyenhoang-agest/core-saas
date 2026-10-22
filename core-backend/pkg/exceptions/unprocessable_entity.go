package exceptions

import (
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions/err_type"
	"github.com/gofiber/fiber/v3"
)

var _ CustomError = (*UnprocessableEntityError)(nil)
var _ error = (*UnprocessableEntityError)(nil)

type UnprocessableEntityError struct {
	ResponseError
}

func NewUnprocessableEntityError(responseCode response_code.UnprocessableEntityCode) UnprocessableEntityError {
	return UnprocessableEntityError{
		ResponseError: NewResponseError(responseCode, err_type.UnprocessableEntityError, nil),
	}
}

func (b UnprocessableEntityError) StatusCode() int {
	return fiber.StatusUnprocessableEntity
}

func (b UnprocessableEntityError) Error() string {
	return b.MessageCode
}

func (b UnprocessableEntityError) Json() map[string]any {
	return map[string]any{
		"status":       b.Status,
		"message_code": b.MessageCode,
		"error_detail": b.ErrorDetail,
	}
}

func (b UnprocessableEntityError) FormattedLog() string {
	return "UnprocessableEntityError: " + b.MessageCode
}
