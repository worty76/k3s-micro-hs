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
