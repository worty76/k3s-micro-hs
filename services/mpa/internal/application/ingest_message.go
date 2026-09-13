package application

import (
	"context"

	"github.com/worty76/k3s-micro-hs/libs/common/logger"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/domain/message"
)

type IngestMessage struct {
	logger logger.Logger
}

func NewIngestMessage(logger logger.Logger) *IngestMessage {
	return &IngestMessage{
		logger: logger,
	}
}

func (im *IngestMessage) Ingest(ctx context.Context, msg message.Message) error {
	im.logger.Info("received message", logger.Field{Key: "payload", Value: string(msg.Payload)})
	return nil
}
