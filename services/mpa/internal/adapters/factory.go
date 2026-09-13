package adapters

import (
	"fmt"

	"github.com/worty76/k3s-micro-hs/services/mpa/internal/adapters/transport/http"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/adapters/transport/mqtt"
)

type Protocol string

const (
	ProtocolHTTP Protocol = "http"
	ProtocolMQTT Protocol = "mqtt"
)

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) CreateAdapter(protocol Protocol) (Runnable, error) {
	switch protocol {
	case ProtocolHTTP:
		return http.NewHTTPAdapter(), nil
	case ProtocolMQTT:
		return mqtt.NewMQTTAdapter(), nil
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}
}
