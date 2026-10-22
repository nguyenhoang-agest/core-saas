package response_code

import (
	"net/http"

	"github.com/example/core-backend/pkg/exceptions/err_type"
)

type InternalServerCode int

// InternalServer response code: E02****
const (
	_ InternalServerCode = iota
	InternalServerGeneral
	DatabaseError
	ThirdPartyServiceError
	// ONLY APPEND ERROR CODE
)

func (c InternalServerCode) Code() string {
	return FormatCode(err_type.InternalServer, int(c))
}

func (c InternalServerCode) Error() string {
	return c.Code()
}

func (c InternalServerCode) HTTPStatus() int {
	return http.StatusInternalServerError
}
