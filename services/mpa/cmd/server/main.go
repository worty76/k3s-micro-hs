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

	currentEnv := env.Environment(os.Getenv("env"))
	appLogger := logger.NewZap(currentEnv)

	// Load configs
	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Fatal("load config", logger.Field{Key: "error", Value: err})
	}
	appLogger.Info("resolved config",
		logger.Field{Key: "broker", Value: cfg.Mqtt.Broker},
		logger.Field{Key: "port", Value: cfg.Mqtt.Port},
		logger.Field{Key: "topics", Value: cfg.Mqtt.Topics},
	)

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

	mqttAdapter := mqtt.NewMQTTAdapter(
		mqttClient,
		messageIngestor,
		appLogger,
		mqttMapper,
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
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
