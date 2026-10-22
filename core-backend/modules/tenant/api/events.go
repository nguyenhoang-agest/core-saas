package api

const (
	TenantCreatedTopic   = "tenant.created"
	TenantUpdatedTopic   = "tenant.updated"
	TenantSuspendedTopic = "tenant.suspended"
	TenantDeletedTopic   = "tenant.deleted"
)

type TenantCreatedEvent struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	IsolateMode   string `json:"isolate_mode"`
	IsHost        bool   `json:"is_host"`
	OwnerEmail    string `json:"owner_email,omitempty"`
	OwnerName     string `json:"owner_name,omitempty"`
	OwnerPassword string `json:"owner_password,omitempty"`
}

type TenantUpdatedEvent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	IsolateMode string `json:"isolate_mode"`
}

type TenantStatusEvent struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
