package response_code

import (
	"net/http"

	"github.com/example/core-backend/pkg/exceptions/err_type"
)

type NotFoundCode int

// NotFound response code: E04****
const (
	_ NotFoundCode = iota
	NotFoundGeneral
	NotFoundResource
	// ONLY APPEND ERROR CODE
)

func (c NotFoundCode) Code() string {
	return FormatCode(err_type.NotFound, int(c))
}

func (c NotFoundCode) Error() string {
	return c.Code()
}

func (c NotFoundCode) HTTPStatus() int {
	return http.StatusNotFound
}
