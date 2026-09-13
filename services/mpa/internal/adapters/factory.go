package adapters

import (
	"fmt"
)

type Protocol string

const (
	ProtocolHTTP Protocol = "http"
	ProtocolMQTT Protocol = "mqtt"
)

type Factory struct {
	adapters map[Protocol]Runnable
}

func NewFactory(adapters map[Protocol]Runnable) *Factory {
	return &Factory{
		adapters: adapters,
	}
}

func (f *Factory) CreateAdapter(protocol Protocol) (Runnable, error) {
	adapter, ok := f.adapters[protocol]
	if !ok {
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}

	return adapter, nil
}
