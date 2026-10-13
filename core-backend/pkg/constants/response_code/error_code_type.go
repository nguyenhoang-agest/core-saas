package response_code

import "fmt"

type IResponseCode interface {
	error
	Code() string
	HTTPStatus() int
}

// ResponseCode is the generic union constraint of all category response code types.
type ResponseCode interface {
	UnauthorizedCode | InternalServerCode | BadRequestCode | NotFoundCode | ValidationCode | ForbiddenCode | TooManyRequestCode | UnsupportedMediaCode | UnprocessableEntityCode | LDAPCode | SuccessCode
}

const (
	ModuleCommonOffset       = 0
	ModuleTenantOffset       = 1000
	ModuleIAMOffset          = 2000
	ModuleAuditOffset        = 3000
	ModuleNotificationOffset = 4000
	ModuleExampleOffset      = 5000
)

func FormatCode(group int, code int) string {
	return fmt.Sprintf("E%02d%04d", group, code)
}
