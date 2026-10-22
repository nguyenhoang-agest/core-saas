package exceptions

import (
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions/err_type"
	"github.com/gofiber/fiber/v3"
)

var _ CustomError = (*TooManyRequestError)(nil)
var _ error = (*TooManyRequestError)(nil)

type TooManyRequestError struct {
	ResponseError
}

func NewTooManyRequestError(responseCode response_code.TooManyRequestCode) TooManyRequestError {
	return TooManyRequestError{
		ResponseError: NewResponseError(responseCode, err_type.TooManyRequestError, nil),
	}
}

func (b TooManyRequestError) StatusCode() int {
	return fiber.StatusTooManyRequests
}

func (b TooManyRequestError) Error() string {
	return b.MessageCode
}

func (b TooManyRequestError) Json() map[string]any {
	return map[string]any{
		"status":       b.Status,
		"message_code": b.MessageCode,
		"error_detail": b.ErrorDetail,
	}
}

func (b TooManyRequestError) FormattedLog() string {
	return "TooManyRequestError: " + b.MessageCode
}
