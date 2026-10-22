package errcode

import "github.com/example/core-backend/pkg/constants/response_code"

// Tenant 500 Internal Server codes: E021001+
const (
	_ response_code.InternalServerCode = response_code.ModuleTenantOffset + iota
	FailToCreateTenant
	FailToGetTenant
	FailToUpdateTenant
	FailToDeleteTenant
	// ONLY APPEND ERROR CODE
)
