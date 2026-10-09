package response_code

import "net/http"

type ValidationCode int

const (
	_ ValidationCode = iota
	ValidationGeneral
	ValidationFailed
	InvalidEmailFormat
	InvalidPhoneFormat
	StringTooShort
	StringTooLong
)

func (c ValidationCode) Code() string {
	return FormatCode(5, int(c))
}

func (c ValidationCode) Error() string {
	return c.Code()
}

func (c ValidationCode) HTTPStatus() int {
	return http.StatusUnprocessableEntity
}
