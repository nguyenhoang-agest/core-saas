package response_code

import (
	"net/http"

	"github.com/example/core-backend/pkg/exceptions/err_type"
)

type BadRequestCode int

// BadRequest response code: E03****
const (
	_ BadRequestCode = iota
	BadRequestGeneral
	BadRequestValidation
	InvalidPayload
	InvalidParameter
	MissingRequiredField
	// ONLY APPEND ERROR CODE
)

func (c BadRequestCode) Code() string {
	return FormatCode(err_type.BadRequest, int(c))
}

func (c BadRequestCode) Error() string {
	return c.Code()
}

func (c BadRequestCode) HTTPStatus() int {
	return http.StatusBadRequest
}
