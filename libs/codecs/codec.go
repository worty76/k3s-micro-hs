package codecs

import "github.com/worty76/k3s-micro-hs/libs/canonical"

type Codec interface {
	Decode(payload []byte) (canonical.Message, error)
}
