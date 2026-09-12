package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/worty76/k3s-micro-hs/services/mpa/internal/adapters"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	e := echo.New()

	// Initialize factory
	factory := adapters.NewFactory()
	mqttClient, err := factory.CreateAdapter(adapters.ProtocolMQTT)
	if err != nil {
		panic(err)
	}

	// Start the MQTT client
	go func() {
		if err := mqttClient.Start(ctx); err != nil {
			panic(err)
		}
	}()

	// Start the HTTP server
	go func() {
		if err := e.Start(":8080"); err != nil {
			panic(err)
		}
	}()

	// Wait for SIGTERM/SIGINT or component failure.
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := e.Shutdown(shutdownCtx); err != nil {
		panic(err)
	}

	// Gracefully stop the MQTT client
	if err := mqttClient.Shutdown(shutdownCtx); err != nil {
		panic(err)
	}
}
