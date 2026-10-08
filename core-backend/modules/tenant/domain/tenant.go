package domain

type Status string

const (
	StatusActive    Status = "active"
	StatusSuspended Status = "suspended"
)

type Tenant struct {
	BaseModel   `gorm:"embedded"`
	Name        string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string `gorm:"column:description;type:varchar(255)" json:"description"`
	Status      Status `gorm:"column:status;type:varchar(50);not null" json:"status"`
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
