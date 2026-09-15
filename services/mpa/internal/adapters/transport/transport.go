package transport

import "github.com/worty76/k3s-micro-hs/services/mpa/internal/domain/message"

type MessageMapper interface {
	Map(payload []byte) (message.Message, error)
}
