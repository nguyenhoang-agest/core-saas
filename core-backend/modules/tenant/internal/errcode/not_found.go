package errcode

import "github.com/example/core-backend/pkg/constants/response_code"

// Tenant 404 Not Found codes: E041001+
const (
	_ response_code.NotFoundCode = response_code.ModuleTenantOffset + iota
	NotFoundTenant
	// ONLY APPEND ERROR CODE
)
