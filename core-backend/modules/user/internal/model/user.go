package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/example/core-backend/pkg/rbac"
)

type Status string

const (
	StatusActive    Status = "active"
	StatusSuspended Status = "suspended"
	StatusInvited   Status = "invited"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusSuspended, StatusInvited:
		return true
	default:
		return false
	}
}

// StringArray implements database/sql driver.Valuer and Scanner for PostgreSQL text[].
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	parts := make([]string, len(a))
	for i, s := range a {
		escaped := strings.ReplaceAll(s, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		parts[i] = "\"" + escaped + "\""
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

func (a *StringArray) Scan(src any) error {
	if src == nil {
		*a = []string{}
		return nil
	}
	var str string
	switch v := src.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("cannot scan %T into StringArray", src)
	}

	str = strings.TrimSpace(str)
	if str == "" || str == "{}" {
		*a = []string{}
		return nil
	}

	if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
		str = str[1 : len(str)-1]
	}

	var items []string
	var current strings.Builder
	inQuote := false
	escaped := false

	for _, ch := range str {
		if escaped {
			current.WriteRune(ch)
			escaped = false
			continue
		}
		if ch == '\\' {
			escaped = true
			continue
		}
		if ch == '"' {
			inQuote = !inQuote
			continue
		}
		if ch == ',' && !inQuote {
			items = append(items, current.String())
			current.Reset()
			continue
		}
		current.WriteRune(ch)
	}
	if current.Len() > 0 || inQuote {
		items = append(items, current.String())
	}

	*a = items
	return nil
}

type User struct {
	BaseModel   `gorm:"embedded"`
	TenantID    string      `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_users_tenant" json:"tenant_id"`
	Email       string      `gorm:"column:email;type:varchar(255);not null;index:idx_users_email" json:"email"`
	Name        string      `gorm:"column:name;type:varchar(255);not null;index:idx_users_name" json:"name"`
	Role        string      `gorm:"column:role;type:varchar(50);not null;default:'member'" json:"role"`
	Permissions StringArray `gorm:"column:permissions;type:text[];not null;default:'{}'" json:"permissions"`
	Status      Status      `gorm:"column:status;type:varchar(50);not null;default:'active'" json:"status"`
	Phone       string      `gorm:"column:phone;type:varchar(50)" json:"phone,omitempty"`
	AvatarURL   string      `gorm:"column:avatar_url;type:varchar(512)" json:"avatar_url,omitempty"`
	LastLoginAt *time.Time  `gorm:"column:last_login_at;type:timestamptz" json:"last_login_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (t *User) IsActive() bool {
	return t.Status == StatusActive
}

// EffectivePermissions calculates the combined permissions of the role and direct ACL grants.
func (t *User) EffectivePermissions() []string {
	return rbac.GetEffectivePermissions(t.Role, t.Permissions)
}
