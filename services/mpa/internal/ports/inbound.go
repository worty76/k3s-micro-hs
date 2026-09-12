package ports

import "context"

type InboundAdapter interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
