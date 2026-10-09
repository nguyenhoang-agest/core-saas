package response_code

import "net/http"

type InternalServerCode int

const (
	_ InternalServerCode = iota
	InternalServerGeneral
	DatabaseError
	ThirdPartyServiceError
)

func (c InternalServerCode) Code() string {
	return FormatCode(50, int(c))
}

func (c InternalServerCode) Error() string {
	return c.Code()
}

func (c InternalServerCode) HTTPStatus() int {
	return http.StatusInternalServerError
}
