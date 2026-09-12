package mqtt

import "context"

type MQTTAdapter struct {
	// mqtt client, config, etc.
}

func NewMQTTAdapter() *MQTTAdapter {
	return &MQTTAdapter{}
}

func (a *MQTTAdapter) Start(ctx context.Context) error {
	// Connect to EMQX
	// Subscribe to topics
	// Receive messages
	return nil
}

func (a *MQTTAdapter) Shutdown(ctx context.Context) error {
	// Disconnect from EMQX
	return nil
}
