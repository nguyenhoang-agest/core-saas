package exceptions

import (
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions/err_type"
	"github.com/gofiber/fiber/v3"
)

var _ CustomError = (*UnsupportedMediaError)(nil)
var _ error = (*UnsupportedMediaError)(nil)

type UnsupportedMediaError struct {
	ResponseError
}

func NewUnsupportedMediaError(responseCode response_code.UnsupportedMediaCode) UnsupportedMediaError {
	return UnsupportedMediaError{
		ResponseError: NewResponseError(responseCode, err_type.UnsupportedMediaError, nil),
	}
}

func (b UnsupportedMediaError) StatusCode() int {
	return fiber.StatusUnsupportedMediaType
}

func (b UnsupportedMediaError) Error() string {
	return b.MessageCode
}

func (b UnsupportedMediaError) Json() map[string]any {
	return map[string]any{
		"status":       b.Status,
		"message_code": b.MessageCode,
		"error_detail": b.ErrorDetail,
	}
}

func (b UnsupportedMediaError) FormattedLog() string {
	return "UnsupportedMediaError: " + b.MessageCode
}
