package integrations

import "github.com/worty76/k3s-micro-hs/libs/canonical"

type Envelope struct {
	Protocol string `json:"protocol"`
}

type Integration interface {
	Name() string
	Decode(envelope Envelope, payload []byte) (canonical.Message, error)
}
