package adapters

import (
	"context"
	"log/slog"

	"github.com/example/core-backend/modules/notification/domain"
)

// LogSender is the default development adapter that logs notifications.
type LogSender struct {
	logger *slog.Logger
}

func NewLogSender() *LogSender {
	return &LogSender{
		logger: slog.Default(),
	}
}

func (s *LogSender) Send(ctx context.Context, msg *domain.Message) error {
	s.logger.InfoContext(ctx, "notification_dispatched_stub",
		slog.String("channel", string(msg.Channel)),
		slog.String("recipient", msg.Recipient),
		slog.String("subject", msg.Subject),
	)
	return nil
}
