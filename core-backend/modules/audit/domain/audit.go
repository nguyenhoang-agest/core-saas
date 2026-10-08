package domain

import (
	"context"
	"time"
)

// LogEntry represents an immutable record of an actor performing an action.
type LogEntry struct {
	ID        string         `json:"id"`
	TenantID  string         `json:"tenant_id"`
	ActorID   string         `json:"actor_id"`
	Action    string         `json:"action"`
	Resource  string         `json:"resource"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

// Repository defines data access for audit log persistence.
type Repository interface {
	Record(ctx context.Context, entry *LogEntry) error
	ListByTenant(ctx context.Context, tenantID string) ([]LogEntry, error)
}
