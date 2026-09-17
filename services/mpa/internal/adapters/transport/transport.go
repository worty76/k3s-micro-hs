package transport

import "github.com/worty76/k3s-micro-hs/libs/canonical"

type MessageMapper interface {
	Map(payload []byte) (canonical.Message, error)
}
