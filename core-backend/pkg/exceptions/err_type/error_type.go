package err_type

const (
	Success = iota
	UnauthorizedError
	InternalServer
	BadRequest
	NotFound
	ValidationError
	ForbiddenError
	TooManyRequestError
	UnsupportedMediaError
	UnprocessableEntityError
	LDAPError
	UnknownError

	// ONLY APPEND ERROR CODE
)
