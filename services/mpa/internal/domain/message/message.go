package message

import "github.com/worty76/k3s-micro-hs/libs/common/entity"

type MessageType string

const (
	MessageTypeTelemetry MessageType = "telemetry"
	MessageTypeEvent     MessageType = "event"
	MessageTypeCommand   MessageType = "command"
	MessageTypeAck       MessageType = "ack"
)

type Message struct {
	entity.BaseEntity
	Type      string `json:"type"`
	Payload   string `json:"payload"`
	Metadata  string `json:"metadata"`
	Direction string `json:"direction"`
}
