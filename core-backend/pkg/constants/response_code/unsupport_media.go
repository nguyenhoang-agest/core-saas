package response_code

import "net/http"

type UnsupportedMediaCode int

const (
	_ UnsupportedMediaCode = iota
	UnsupportedMediaTypeGeneral
	InvalidContentType
	UnsupportedFileExtension
)

func (c UnsupportedMediaCode) Code() string {
	return FormatCode(9, int(c))
}

func (c UnsupportedMediaCode) Error() string {
	return c.Code()
}

func (c UnsupportedMediaCode) HTTPStatus() int {
	return http.StatusUnsupportedMediaType
}
