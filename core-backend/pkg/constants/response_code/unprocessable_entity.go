package response_code

import (
	"net/http"

	"github.com/example/core-backend/pkg/exceptions/err_type"
)

type UnprocessableEntityCode int

// UnprocessableEntity response code: E09****
const (
	_ UnprocessableEntityCode = iota
	UnprocessableGeneral
	BusinessRuleViolation
	StateTransitionInvalid
	// ONLY APPEND ERROR CODE
)

func (c UnprocessableEntityCode) Code() string {
	return FormatCode(err_type.UnprocessableEntityError, int(c))
}

func (c UnprocessableEntityCode) Error() string {
	return c.Code()
}

func (c UnprocessableEntityCode) HTTPStatus() int {
	return http.StatusUnprocessableEntity
}
