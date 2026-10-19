package errcode

import "github.com/example/core-backend/pkg/constants/response_code"

// Tenant 400 Bad Request codes: E031001+
const (
	_ response_code.BadRequestCode = response_code.ModuleTenantOffset + iota
	TenantNameRequired
	TenantAlreadyExists
	InvalidIsolateMode
	MissingDedicatedDBConfig
	InvalidDedicatedDBConfig
	// ONLY APPEND ERROR CODE
)
