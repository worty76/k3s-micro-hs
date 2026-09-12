package application

import (
	"context"

	"github.com/worty76/k3s-micro-hs/services/mpa/internal/domain/message"
)

type IngestMessage struct {
}

func (im *IngestMessage) Execute(ctx context.Context, msg *message.Message) error {
	return nil
}
