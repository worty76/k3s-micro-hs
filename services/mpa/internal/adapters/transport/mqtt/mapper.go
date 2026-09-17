package mqtt

import (
	"encoding/json"
	"fmt"

	"github.com/worty76/k3s-micro-hs/libs/canonical"
)

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) Map(payload []byte) (canonical.Message, error) {
	var msg canonical.Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		return canonical.Message{}, fmt.Errorf("unmarshal mqtt payload: %w", err)
	}
	msg.Direction = canonical.DirectionUplink
	return msg, nil
}
