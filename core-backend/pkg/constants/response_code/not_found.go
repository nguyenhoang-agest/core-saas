package response_code

import "net/http"

type NotFoundCode int

const (
	_ NotFoundCode = iota
	NotFoundGeneral
	NotFoundTenant
	NotFoundUser
	NotFoundResource
)

func (c NotFoundCode) Code() string {
	return FormatCode(3, int(c))
}

func (c NotFoundCode) Error() string {
	return c.Code()
}

func (c NotFoundCode) HTTPStatus() int {
	return http.StatusNotFound
}
