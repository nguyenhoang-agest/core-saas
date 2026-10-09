package response_code

import "net/http"

type LDAPCode int

const (
	_ LDAPCode = iota
	LDAPGeneralError
	LDAPBindFailed
	LDAPUserNotFound
	LDAPConnectionFailed
)

func (c LDAPCode) Code() string {
	return FormatCode(6, int(c))
}

func (c LDAPCode) Error() string {
	return c.Code()
}

func (c LDAPCode) HTTPStatus() int {
	return http.StatusUnauthorized
}
