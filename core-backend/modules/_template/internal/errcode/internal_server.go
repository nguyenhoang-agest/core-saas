package errcode

import "github.com/example/core-backend/pkg/constants/response_code"

// 500 Internal Server codes
const (
	_ response_code.InternalServerCode = ModuleOffset + iota
	FailToCreateTemplate
	FailToGetTemplate
	FailToUpdateTemplate
	FailToDeleteTemplate
	// ONLY APPEND ERROR CODE
)
