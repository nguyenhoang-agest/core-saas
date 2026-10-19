package errcode

import "github.com/example/core-backend/pkg/constants/response_code"

// 404 Not Found codes
const (
	_ response_code.NotFoundCode = ModuleOffset + iota
	NotFoundUser
	// ONLY APPEND ERROR CODE
)
