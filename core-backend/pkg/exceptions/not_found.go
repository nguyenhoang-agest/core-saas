package exceptions

import (
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions/err_type"
	"github.com/gofiber/fiber/v3"
)

var _ CustomError = (*NotFoundError)(nil)
var _ error = (*NotFoundError)(nil)

type NotFoundError struct {
	ResponseError
}

func NewNotFoundError(responseCode response_code.NotFoundCode) NotFoundError {
	return NotFoundError{
		ResponseError: NewResponseError(responseCode, err_type.NotFound, nil),
	}
}

func (b NotFoundError) StatusCode() int {
	return fiber.StatusNotFound
}

func (b NotFoundError) Error() string {
	return b.MessageCode
}

func (b NotFoundError) Json() map[string]any {
	return map[string]any{
		"status":       b.Status,
		"message_code": b.MessageCode,
		"error_detail": b.ErrorDetail,
	}
}

func (b NotFoundError) FormattedLog() string {
	return "NotFoundError: " + b.MessageCode
}
