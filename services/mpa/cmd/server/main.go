package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/worty76/k3s-micro-hs/libs/common/env"
	"github.com/worty76/k3s-micro-hs/libs/common/logger"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/adapters"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	e := echo.New()

	currentEnv := env.Environment(os.Getenv("env"))

	// Initialize logger
	appLogger := logger.NewZap(currentEnv)

	appLogger.Info("Starting MPA service...")

	// Initialize factory
	factory := adapters.NewFactory()
	mqttClient, err := factory.CreateAdapter(adapters.ProtocolMQTT)
	if err != nil {
		panic(err)
	}

	g, groupCtx := errgroup.WithContext(ctx)

	// Start the MQTT client
	g.Go(func() error {
		appLogger.Info("Starting MQTT client...")
		return mqttClient.Start(groupCtx)
	})

	// Start the HTTP server
	g.Go(func() error {
		appLogger.Info("Starting HTTP server on port 8080...")
		return e.Start(":8080")
	})

	// // Wait for SIGTERM/SIGINT or component failure.
	// <-ctx.Done()

	g.Go(func() error {
		<-groupCtx.Done()
		appLogger.Info("Context canceled, shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Gracefully shutdown the server
		if err := e.Shutdown(shutdownCtx); err != nil {
			appLogger.Error("Error occurred while shutting down server", logger.Field{Key: "error", Value: err})
			panic(err)
		}

		// And then shutdown the other adapters
		shutdown(shutdownCtx, e, mqttClient, appLogger)
		return nil
	})
}

func shutdown(shutdownCtx context.Context, e *echo.Echo, mqttClient adapters.Runnable, appLogger logger.Logger) {
	// Gracefully stop the MQTT client
	if err := mqttClient.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("Error occurred while shutting down MQTT client", logger.Field{Key: "error", Value: err})
		panic(err)
	}
}
