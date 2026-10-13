package exceptions

import (
	"github.com/example/core-backend/pkg/constants/response_code"
	"github.com/example/core-backend/pkg/exceptions/err_type"
	"github.com/gofiber/fiber/v3"
)

var _ CustomError = (*BadRequestError)(nil)
var _ error = (*BadRequestError)(nil)

type BadRequestError struct {
	ResponseError
}

func NewValidationError(responseCode response_code.ValidationCode, errorDetail []FieldErrorDetail) BadRequestError {
	return BadRequestError{
		ResponseError: NewResponseError(responseCode, err_type.ValidationError, errorDetail),
	}
}

func NewBadRequestError(responseCode response_code.BadRequestCode, errorDetail []FieldErrorDetail) BadRequestError {
	return BadRequestError{
		ResponseError: NewResponseError(responseCode, err_type.BadRequest, errorDetail),
	}
}

func (b BadRequestError) StatusCode() int {
	return fiber.StatusBadRequest
}

func (b BadRequestError) Error() string {
	return b.MessageCode
}

func (b BadRequestError) Json() map[string]any {
	return map[string]any{
		"status":       b.Status,
		"message_code": b.MessageCode,
		"error_detail": b.ErrorDetail,
	}
}

func (b BadRequestError) FormattedLog() string {
	return "BadRequestError: " + b.MessageCode
}
