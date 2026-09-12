package ports

import (
	"context"

	"github.com/worty76/k3s-micro-hs/services/mpa/internal/domain/message"
)

type MessagePublisher interface {
	Publish(ctx context.Context, msg message.Message) error
}
