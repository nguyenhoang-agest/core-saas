package response_code

import "net/http"

type ForbiddenCode int

const (
	_ ForbiddenCode = iota
	ForbiddenGeneral
	AccessDenied
	InsufficientPermission
	TenantSuspended
	TenantAccessDenied
)

func (c ForbiddenCode) Code() string {
	return FormatCode(2, int(c))
}

func (c ForbiddenCode) Error() string {
	return c.Code()
}

func (c ForbiddenCode) HTTPStatus() int {
	return http.StatusForbidden
}
