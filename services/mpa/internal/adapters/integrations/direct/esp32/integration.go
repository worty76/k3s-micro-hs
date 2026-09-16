package esp32

import (
	"encoding/json"
	"fmt"

	"github.com/worty76/k3s-micro-hs/libs/canonical"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/adapters/integrations"
)

const (
	Esp32IntegrationName = "esp32"
)

type ESP32Integration struct {
	topics []string
}

func NewESP32Integration(topics []string) *ESP32Integration {
	return &ESP32Integration{
		topics: topics,
	}
}

func (e *ESP32Integration) Name() string {
	return Esp32IntegrationName
}

type devicePayload struct {
	DeviceID string `json:"device_id"`
	Payload  string `json:"payload"`
}

func (e *ESP32Integration) Decode(envelope integrations.Envelope, payload []byte) (canonical.Message, error) {
	// Implement the decoding logic specific to ESP32 here
	// For example, you might want to parse the payload and create a canonical.Message

	fmt.Printf("Decoding ESP32 message with protocol: %s\n", envelope.Protocol)

	var dp devicePayload
	if err := json.Unmarshal(payload, &dp); err != nil {
		return canonical.Message{}, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Logging raw payload from the device
	fmt.Printf("Raw payload: %s\n", dp.Payload)

	// This is a placeholder implementation. You should replace it with actual decoding logic.
	message := canonical.Message{
		Type:      "telemetry",
		Payload:   dp.Payload,
		Metadata:  nil, // Replace with actual metadata if needed
		Direction: canonical.DirectionUplink,
		Source:    canonical.Source{Integration: e.Name(), Protocol: envelope.Protocol, Network: "wifi", DeviceID: "esp32_device_id"},
	}

	return message, nil
}
