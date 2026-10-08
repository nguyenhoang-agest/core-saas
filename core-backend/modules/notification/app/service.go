package app

import (
	"context"

	"github.com/example/core-backend/modules/notification/domain"
	"github.com/example/core-backend/pkg/errors"
	"github.com/example/core-backend/pkg/eventbus"
)

// Service handles notification dispatching (STUB).
type Service struct {
	sender domain.Sender
	bus    *eventbus.Bus
}

func NewService(sender domain.Sender, bus *eventbus.Bus) *Service {
	return &Service{
		sender: sender,
		bus:    bus,
	}
}

func (s *Service) Dispatch(ctx context.Context, msg *domain.Message) error {
	if msg == nil || msg.Recipient == "" {
		return errors.InvalidArgument("invalid notification message", nil)
	}
	return s.sender.Send(ctx, msg)
}
