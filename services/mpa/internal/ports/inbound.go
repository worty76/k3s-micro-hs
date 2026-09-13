package ports

import (
	"context"

	"github.com/worty76/k3s-micro-hs/services/mpa/internal/domain/message"
)

type MessageIngestor interface {
	Ingest(ctx context.Context, msg message.Message) error
}
