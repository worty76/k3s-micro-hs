package mqtt

import (
	"context"
	"fmt"

	"github.com/worty76/k3s-micro-hs/libs/common/logger"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/adapters"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/ports"
)

type MQTTAdapter struct {
	client   *Client
	ingestor ports.MessageIngestor
	logger   logger.Logger
	mapper   adapters.MessageMapper
}

func NewMQTTAdapter(client *Client, ingestor ports.MessageIngestor, logger logger.Logger, mapper adapters.MessageMapper) *MQTTAdapter {
	return &MQTTAdapter{
		client:   client,
		ingestor: ingestor,
		logger:   logger,
		mapper:   mapper,
	}
}

func (a *MQTTAdapter) Start(ctx context.Context) error {
	if err := a.client.Validate(); err != nil {
		return fmt.Errorf("failed to validate EMQX client: %w", err)
	}
	a.logger.Info("Validated EMQX client configuration, ready to connect")

	// Connect to EMQX
	if err := a.client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to EMQX: %w", err)
	}
	a.logger.Info("Connected to EMQX")

	// Subscribe to topics
	if err := a.client.Subscribe(func(topic string, payload []byte) {
		// Route: msg -> mapper.Map -> ingestor.Ingest
		msg, err := a.mapper.Map(payload)
		if err != nil {
			a.logger.Error("Failed to map message", logger.Field{Key: "error", Value: err})
			// return here to avoid ingesting a malformed message
			return
		}

		// Ingest the mapped message
		if err := a.ingestor.Ingest(ctx, msg); err != nil {
			a.logger.Error("Failed to ingest message", logger.Field{Key: "error", Value: err})
		}
	}); err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	return nil
}

func (a *MQTTAdapter) Shutdown(ctx context.Context) error {
	// Disconnect from EMQX
	return a.client.Disconnect()
}
