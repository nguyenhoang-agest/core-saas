package exceptions

import (
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions/err_type"
	"github.com/gofiber/fiber/v3"
)

var _ CustomError = (*InternalServerError)(nil)
var _ error = (*InternalServerError)(nil)

type InternalServerError struct {
	ResponseError
}

func NewInternalServerError(responseCode response_code.InternalServerCode) InternalServerError {
	return InternalServerError{
		ResponseError: NewResponseError(responseCode, err_type.InternalServer, nil),
	}
}

func (b InternalServerError) StatusCode() int {
	return fiber.StatusInternalServerError
}

func (b InternalServerError) Error() string {
	return b.MessageCode
}

func (b InternalServerError) Json() map[string]any {
	return map[string]any{
		"status":       b.Status,
		"message_code": b.MessageCode,
		"error_detail": b.ErrorDetail,
	}
}

func (b InternalServerError) FormattedLog() string {
	return "InternalServerError: " + b.MessageCode
}
