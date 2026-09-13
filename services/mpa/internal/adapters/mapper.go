package adapters

import (
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/domain/message"
)

type MessageMapper interface {
	Map(msg []byte) (message.Message, error)
}
