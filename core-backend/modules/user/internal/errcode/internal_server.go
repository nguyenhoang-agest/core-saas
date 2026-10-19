package errcode

import "github.com/example/core-backend/pkg/constants/response_code"

// 500 Internal Server codes
const (
	_ response_code.InternalServerCode = ModuleOffset + iota
	FailToCreateUser
	FailToGetUser
	FailToUpdateUser
	FailToDeleteUser
	// ONLY APPEND ERROR CODE
)
