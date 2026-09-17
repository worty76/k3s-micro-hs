package canonical

type MessageType string

const (
	MessageTypeTelemetry MessageType = "telemetry"
	MessageTypeEvent     MessageType = "event"
	MessageTypeCommand   MessageType = "command"
	MessageTypeAck       MessageType = "ack"

	DirectionUplink   = "uplink"
	DirectionDownlink = "downlink"
)

// IsValidType reports whether t is a known message type.
func IsValidType(t string) bool {
	switch MessageType(t) {
	case MessageTypeTelemetry, MessageTypeEvent, MessageTypeCommand, MessageTypeAck:
		return true
	default:
		return false
	}
}
