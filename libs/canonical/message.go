package canonical

import "github.com/worty76/k3s-micro-hs/libs/common/entity"

type Message struct {
	entity.BaseEntity
	Type      string      `json:"type"`      // Type of the message, e.g., "sensor_data", "command", etc.
	Payload   string      `json:"payload"`   // Payload of the message, can be any data in string format (e.g., JSON, plain text, raw bytes)
	Metadata  interface{} `json:"metadata"`  // Additional metadata about the message, can be used for filtering or routing
	Direction string      `json:"direction"` // Direction of the message, e.g., "uplink" (from device to server) or "downlink" (from server to device)
	Source    Source      `json:"source"`
}
