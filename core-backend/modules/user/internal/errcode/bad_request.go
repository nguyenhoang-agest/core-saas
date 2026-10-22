package errcode

import "github.com/example/core-backend/pkg/constants/response_code"

// Module offset for error codes
const ModuleOffset = 2000

// 400 Bad Request codes
const (
	_ response_code.BadRequestCode = ModuleOffset + iota
	UserEmailRequired              // 2001
	UserAlreadyExists              // 2002
	UserInvalidRole                // 2003
	UserInvalidStatus              // 2004
	UserNameRequired               // 2005
	RoleInUse                      // 2006 (E032006)
	CannotModifySystemRole         // 2007 (E032007)
	RoleSlugAlreadyExists          // 2008 (E032008)
	PrivilegeEscalation            // 2009 (E032009 - Subset Rule violation)
	OutOfScopePermission           // 2010 (E032010 - Platform scope in tenant role)
	LastOwnerViolation             // 2011 (E032011 - Cannot demote last owner)
	UnknownPermissionKey           // 2012 (E032012 - Unknown or wildcard key)
	ConcurrencyConflict            // 2013 (E032013 - Optimistic locking conflict)
	RoleNameRequired               // 2014
	RoleSlugRequired               // 2015
	OrgNameRequired                // 2016 (E032016)
	OrgAlreadyExists               // 2017 (E032017)
	PasswordTooShort               // 2018 (E032018)
	// ONLY APPEND ERROR CODE
)
