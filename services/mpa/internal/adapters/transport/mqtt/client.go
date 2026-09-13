package mqtt

import (
	"context"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/worty76/k3s-micro-hs/libs/common/logger"
)

type Config struct {
	BrokerURL string
	ClientID  string
	Username  string
	Password  string
	Topics    []string
}

type MessageHandler func(topic string, payload []byte)

type Client struct {
	cfg     Config
	cLogger logger.Logger
	inner   mqtt.Client
}

func NewMQTTClient(cfg Config, cLogger logger.Logger) *Client {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(cfg.BrokerURL)
	opts.SetClientID(cfg.ClientID)

	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	inner := mqtt.NewClient(opts)

	return &Client{
		cfg:     cfg,
		cLogger: cLogger,
		inner:   inner,
	}
}

func NewMQTTClientConfig(brokerURL, clientID, username, password string, topics []string) *Config {
	return &Config{
		BrokerURL: brokerURL,
		ClientID:  clientID,
		Username:  username,
		Password:  password,
		Topics:    topics,
	}
}

func (c *Client) Validate() error {
	if c.cfg.BrokerURL == "" {
		c.cLogger.Error("broker URL is required", logger.Field{Key: "brokerURL", Value: c.cfg.BrokerURL})
		return fmt.Errorf("broker URL is required")
	}
	if c.cfg.ClientID == "" {
		c.cLogger.Error("client ID is required", logger.Field{Key: "clientID", Value: c.cfg.ClientID})
		return fmt.Errorf("client ID is required")
	}
	if len(c.cfg.Topics) == 0 {
		c.cLogger.Error("at least one topic is required", logger.Field{Key: "topics", Value: c.cfg.Topics})
		return fmt.Errorf("at least one topic is required")
	}
	return nil
}

func (c *Client) Connect(ctx context.Context) error {
	token := c.inner.Connect()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-token.Done():
		if err := token.Error(); err != nil {
			return fmt.Errorf("connect to MQTT broker: %w", err)
		}

		return nil
	}
}

func (c *Client) Subscribe(handler MessageHandler) error {
	for _, topic := range c.cfg.Topics {
		c.cLogger.Info("subscribing", logger.Field{Key: "topic", Value: topic})
		token := c.inner.Subscribe(topic, 1, func(
			_ mqtt.Client,
			msg mqtt.Message,
		) {
			handler(msg.Topic(), msg.Payload())
		})

		if token.Wait() && token.Error() != nil {
			return fmt.Errorf(
				"subscribe to topic %q: %w",
				topic,
				token.Error(),
			)
		}
	}

	return nil
}

func (c *Client) Disconnect() error {
	c.inner.Disconnect(1000)
	return nil
}
