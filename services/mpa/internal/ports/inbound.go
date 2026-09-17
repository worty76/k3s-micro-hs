package ports

import (
	"context"

	"github.com/worty76/k3s-micro-hs/libs/canonical"
)

type MessageIngestor interface {
	Ingest(ctx context.Context, msg canonical.Message) error
}
