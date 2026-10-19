package model

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive:
		return true
	default:
		return false
	}
}

type User struct {
	BaseModel   `gorm:"embedded"`
	Name        string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string `gorm:"column:description;type:varchar(255)" json:"description"`
	Status      Status `gorm:"column:status;type:varchar(50);not null" json:"status"`
}

func (User) TableName() string {
	return "users"
}

func (t *User) IsActive() bool {
	return t.Status == StatusActive
}
