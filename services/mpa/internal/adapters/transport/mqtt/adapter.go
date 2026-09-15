package mqtt

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/worty76/k3s-micro-hs/libs/common/logger"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/adapters/transport"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/ports"
)

type MQTTAdapter struct {
	client   *Client
	ingestor ports.MessageIngestor
	logger   logger.Logger
	mapper   transport.MessageMapper

	jobs      chan job
	workers   int
	wg        sync.WaitGroup
	queueSize int

	dropped atomic.Uint64
}

type job struct {
	topic   string
	payload []byte
}

func NewMQTTAdapter(client *Client, ingestor ports.MessageIngestor, logger logger.Logger, mapper transport.MessageMapper, workers int, queueSize int) *MQTTAdapter {
	return &MQTTAdapter{
		client:    client,
		ingestor:  ingestor,
		logger:    logger,
		mapper:    mapper,
		workers:   workers,
		queueSize: queueSize,
	}
}

func (a *MQTTAdapter) Start(ctx context.Context) error {
	// Validate the EMQX client configuration
	if err := a.client.Validate(); err != nil {
		return fmt.Errorf("failed to validate EMQX client: %w", err)
	}
	a.logger.Info("Validated EMQX client configuration, ready to connect")

	// Validate the adapter configuration
	if err := a.Validate(); err != nil {
		return fmt.Errorf("failed to validate MQTT adapter: %w", err)
	}
	a.logger.Info("Validated MQTT adapter configuration, ready to start")

	// Connect to EMQX
	if err := a.client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to EMQX: %w", err)
	}
	a.logger.Info("Connected to EMQX")

	// Intialize the job channel and worker pool
	a.jobs = make(chan job, a.queueSize)

	// Start worker goroutines
	for i := 0; i < a.workers; i++ {
		a.wg.Add(1)
		go func() {
			a.worker(ctx)
		}()
	}

	// Subscribe to topics
	return a.client.Subscribe(a.handleMessage)
}

func (a *MQTTAdapter) Shutdown(ctx context.Context) error {
	// Close the jobs channel to signal workers to stop
	_ = a.client.Disconnect()

	// Wait for all workers to finish processing
	a.wg.Wait()
	a.logger.Info("MQTT adapter shutdown complete")

	return nil
}

func (a *MQTTAdapter) Validate() error {
	if a.client == nil {
		return fmt.Errorf("MQTT client is not initialized")
	}
	if a.ingestor == nil {
		return fmt.Errorf("message ingestor is not initialized")
	}
	if a.logger == nil {
		return fmt.Errorf("logger is not initialized")
	}
	if a.mapper == nil {
		return fmt.Errorf("message mapper is not initialized")
	}
	if a.workers <= 0 {
		return fmt.Errorf("number of workers must be greater than zero")
	}
	if a.queueSize <= 0 {
		return fmt.Errorf("queue size must be greater than zero")
	}
	return nil
}

func (a *MQTTAdapter) handleMessage(topic string, payload []byte) {
	job := job{topic: topic, payload: payload}

	select {
	case a.jobs <- job:
		// Job successfully sent to the channel
	default:
		// Channel is full, drop the job and log the event
		dropped := a.dropped.Add(1)
		a.logger.Warn("Job dropped due to full channel", logger.Field{Key: "dropped_count", Value: dropped})
	}
}

func (a *MQTTAdapter) worker(ctx context.Context) {
	defer a.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-a.jobs:
			if !ok {
				return
			}
			if err := a.process(ctx, j); err != nil {
				a.logger.Error("Failed to process job", logger.Field{Key: "error", Value: err})
			}
		}
	}
}

func (a *MQTTAdapter) process(ctx context.Context, j job) error {
	msg, err := a.mapper.Map(j.payload)
	if err != nil {
		a.logger.Error("Failed to map message", logger.Field{Key: "error", Value: err})
		return nil
	}

	return a.ingestor.Ingest(ctx, msg)
}
