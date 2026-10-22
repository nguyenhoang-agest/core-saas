package exceptions

import (
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions/err_type"
	"github.com/gofiber/fiber/v3"
)

var _ CustomError = (*ForbiddenError)(nil)
var _ error = (*ForbiddenError)(nil)

type ForbiddenError struct {
	ResponseError
}

func NewForbiddenError(responseCode response_code.ForbiddenCode) ForbiddenError {
	return ForbiddenError{
		ResponseError: NewResponseError(responseCode, err_type.ForbiddenError, nil),
	}
}

func (b ForbiddenError) StatusCode() int {
	return fiber.StatusForbidden
}

func (b ForbiddenError) Error() string {
	return b.MessageCode
}

func (b ForbiddenError) Json() map[string]any {
	return map[string]any{
		"status":       b.Status,
		"message_code": b.MessageCode,
		"error_detail": b.ErrorDetail,
	}
}

func (b ForbiddenError) FormattedLog() string {
	return "ForbiddenError: " + b.MessageCode
}
