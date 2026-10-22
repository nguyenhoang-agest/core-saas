package response_code

import (
	"net/http"

	"github.com/example/core-backend/pkg/exceptions/err_type"
)

type ForbiddenCode int

// Forbidden response code: E06****
const (
	_ ForbiddenCode = iota
	ForbiddenGeneral
	AccessDenied
	InsufficientPermission
	// ONLY APPEND ERROR CODE
)

func (c ForbiddenCode) Code() string {
	return FormatCode(err_type.ForbiddenError, int(c))
}

func (c ForbiddenCode) Error() string {
	return c.Code()
}

func (c ForbiddenCode) HTTPStatus() int {
	return http.StatusForbidden
}
