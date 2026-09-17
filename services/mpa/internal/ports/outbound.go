package ports

import (
	"context"

	"github.com/worty76/k3s-micro-hs/libs/canonical"
)

type MessagePublisher interface {
	Publish(ctx context.Context, msg canonical.Message) error
}
