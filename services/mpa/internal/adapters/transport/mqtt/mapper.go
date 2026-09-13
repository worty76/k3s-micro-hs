package mqtt

import (
	"encoding/json"
	"fmt"

	"github.com/worty76/k3s-micro-hs/services/mpa/internal/domain/message"
)

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) Map(payload []byte) (message.Message, error) {
	var msg message.Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		return message.Message{}, fmt.Errorf("unmarshal mqtt payload: %w", err)
	}
	msg.Direction = message.DirectionUplink
	return msg, nil
}
