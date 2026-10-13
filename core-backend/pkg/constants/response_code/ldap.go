package response_code

import (
	"net/http"

	"github.com/example/core-backend/pkg/exceptions/err_type"
)

type LDAPCode int

// LDAP response code: E10****
const (
	_ LDAPCode = iota
	LDAPGeneralError
	LDAPBindFailed
	LDAPUserNotFound
	LDAPConnectionFailed
	// ONLY APPEND ERROR CODE
)

func (c LDAPCode) Code() string {
	return FormatCode(err_type.LDAPError, int(c))
}

func (c LDAPCode) Error() string {
	return c.Code()
}

func (c LDAPCode) HTTPStatus() int {
	return http.StatusUnauthorized
}
