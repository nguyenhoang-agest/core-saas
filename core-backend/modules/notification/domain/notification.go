package domain

import (
	"context"
	"time"
)

// Channel defines notification transmission medium.
type Channel string

const (
	ChannelEmail   Channel = "email"
	ChannelWebhook Channel = "webhook"
	ChannelSMS     Channel = "sms"
)

// Message represents a notification payload.
type Message struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Recipient string    `json:"recipient"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Channel   Channel   `json:"channel"`
	CreatedAt time.Time `json:"created_at"`
}

// Sender is the outbound port for delivering notifications.
type Sender interface {
	Send(ctx context.Context, msg *Message) error
}
