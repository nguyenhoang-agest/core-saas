package response_code

import (
	"net/http"

	"github.com/example/core-backend/pkg/exceptions/err_type"
)

type UnsupportedMediaCode int

// UnsupportedMedia response code: E08****
const (
	_ UnsupportedMediaCode = iota
	UnsupportedMediaTypeGeneral
	InvalidContentType
	UnsupportedFileExtension
	// ONLY APPEND ERROR CODE
)

func (c UnsupportedMediaCode) Code() string {
	return FormatCode(err_type.UnsupportedMediaError, int(c))
}

func (c UnsupportedMediaCode) Error() string {
	return c.Code()
}

func (c UnsupportedMediaCode) HTTPStatus() int {
	return http.StatusUnsupportedMediaType
}
