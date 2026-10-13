package model

type Status string

const (
	StatusActive    Status = "active"
	StatusSuspended Status = "suspended"
)

type IsolateMode string

const (
	IsolateModeShared    IsolateMode = "shared"
	IsolateModeSchema    IsolateMode = "schema"
	IsolateModeDedicated IsolateMode = "dedicated"
)

func (m IsolateMode) IsValid() bool {
	switch m {
	case IsolateModeShared, IsolateModeSchema, IsolateModeDedicated:
		return true
	default:
		return false
	}
}

type TenantDBConfig struct {
	DSN      string `json:"dsn,omitempty"`
	MaxConns int32  `json:"max_conns,omitempty"`
	SSLMode  string `json:"ssl_mode,omitempty"`
}

type Tenant struct {
	BaseModel   `gorm:"embedded"`
	Name        string          `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string          `gorm:"column:description;type:varchar(255)" json:"description"`
	Status      Status          `gorm:"column:status;type:varchar(50);not null" json:"status"`
	IsolateMode IsolateMode     `gorm:"column:isolate_mode;type:varchar(50);not null" json:"isolate_mode"`
	DBConfig    *TenantDBConfig `gorm:"column:db_config;type:jsonb;serializer:json" json:"db_config,omitempty"`
}

func (Tenant) TableName() string {
	return "tenants"
}

func (t *Tenant) IsActive() bool {
	return t.Status == StatusActive
}

func (t *Tenant) Suspend() {
	t.Status = StatusSuspended
}

func (t *Tenant) Activate() {
	t.Status = StatusActive
}
