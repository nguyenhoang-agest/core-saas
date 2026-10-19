package errcode

import "github.com/example/core-backend/pkg/constants/response_code"

// Module offset for error codes
const ModuleOffset = 2000

// 400 Bad Request codes
const (
	_ response_code.BadRequestCode = ModuleOffset + iota
	UserNameRequired
	UserAlreadyExists
	// ONLY APPEND ERROR CODE
)
