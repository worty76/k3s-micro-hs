package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/worty76/k3s-micro-hs/libs/common/env"
	"github.com/worty76/k3s-micro-hs/libs/common/logger"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/adapters"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/adapters/transport/mqtt"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/application"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load configs
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Set environment
	currentEnv := env.Environment(cfg.AppEnv)
	if !currentEnv.IsValid() {
		fmt.Printf("Invalid environment: %s\n", currentEnv)
		os.Exit(1)
	}

	// Initialize logger
	appLogger := logger.NewZap(currentEnv)

	appLogger.Info("Starting MPA service...")

	// Initialize message ingestor
	messageIngestor := application.NewIngestMessage(appLogger)

	// Initialize MQTT client
	mqttCfg := mqtt.Config{
		BrokerURL: fmt.Sprintf("tcp://%s:%d", cfg.Mqtt.Broker, cfg.Mqtt.Port),
		ClientID:  cfg.Mqtt.ClientID,
		Username:  cfg.Mqtt.Username,
		Password:  cfg.Mqtt.Password,
		Topics:    strings.Split(cfg.Mqtt.Topics, ","),
	}

	mqttClient := mqtt.NewMQTTClient(
		mqttCfg,
		appLogger,
	)

	mqttMapper := mqtt.NewMapper()

	if cfg.Mqtt.Workers <= 0 {
		appLogger.Fatal("number of workers must be greater than zero", logger.Field{Key: "workers", Value: cfg.Mqtt.Workers})
	}

	if cfg.Mqtt.QueueSize <= 0 {
		appLogger.Fatal("queue size must be greater than zero", logger.Field{Key: "queueSize", Value: cfg.Mqtt.QueueSize})
	}

	mqttAdapter := mqtt.NewMQTTAdapter(
		mqttClient,
		messageIngestor,
		appLogger,
		mqttMapper,
		cfg.Mqtt.Workers,
		cfg.Mqtt.QueueSize,
	)

	// Initialize transport adapters
	registry := adapters.NewRegistry(
		map[adapters.Protocol]adapters.Runnable{
			adapters.ProtocolMQTT: mqttAdapter,
		},
	)

	e := echo.New()

	g, groupCtx := errgroup.WithContext(ctx)

	// Start transport adapters
	for _, adapter := range registry.All() {
		adapter := adapter

		g.Go(func() error {
			return adapter.Start(groupCtx)
		})
	}

	// Start HTTP server
	g.Go(func() error {
		if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	// Graceful shutdown
	g.Go(func() error {
		<-groupCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = e.Shutdown(shutdownCtx)
		for _, adapter := range registry.All() {
			_ = adapter.Shutdown(shutdownCtx)
		}
		return nil
	})

	// Wait until shutdown/failure
	if err := g.Wait(); err != nil {
		appLogger.Error(
			"MPA stopped",
			logger.Field{
				Key:   "error",
				Value: err,
			},
		)
	}
}
