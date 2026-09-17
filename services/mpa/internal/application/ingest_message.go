package application

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/worty76/k3s-micro-hs/libs/canonical"
	"github.com/worty76/k3s-micro-hs/libs/common/logger"
)

type IngestMessage struct {
	logger logger.Logger
}

func NewIngestMessage(logger logger.Logger) *IngestMessage {
	return &IngestMessage{
		logger: logger,
	}
}

func (im *IngestMessage) Ingest(ctx context.Context, msg canonical.Message) error {
	// Set default values for ID and timestamps if they are not set
	if msg.ID == uuid.Nil {
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}
		msg.ID = id
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now().UTC()
		msg.UpdatedAt = msg.CreatedAt
	}

	// Validation - check if the message is valid
	if !canonical.IsValidType(msg.Type) || msg.Source.Integration == "" || msg.Source.DeviceID == "" {
		im.logger.Warn("message rejected",
			logger.Field{Key: "type", Value: msg.Type},
			logger.Field{Key: "integration", Value: msg.Source.Integration},
			logger.Field{Key: "device", Value: msg.Source.DeviceID},
		)
		return nil
	}

	// Log the accepted message
	im.logger.Info("message accepted",
		logger.Field{Key: "id", Value: msg.ID},
		logger.Field{Key: "integration", Value: msg.Source.Integration},
		logger.Field{Key: "device", Value: msg.Source.DeviceID},
		logger.Field{Key: "type", Value: msg.Type},
	)

	im.logger.Info("received message", logger.Field{Key: "payload", Value: string(msg.Payload)})
	return nil
}
