package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv   string         `mapstructure:"APP_ENV"`
	Database DatabaseConfig `mapstructure:",squash"`
	Server   ServerConfig   `mapstructure:",squash"`
	Mqtt     MQTTConfig     `mapstructure:",squash"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

type ServerConfig struct {
	Port int `mapstructure:"SERVER_PORT"`
}

type MQTTConfig struct {
	Broker   string `mapstructure:"MQTT_BROKER"`
	Port     int    `mapstructure:"MQTT_PORT"`
	ClientID string `mapstructure:"MQTT_CLIENT_ID"`
	Username string `mapstructure:"MQTT_USERNAME"`
	Password string `mapstructure:"MQTT_PASSWORD"`
	Topics   string `mapstructure:"MQTT_TOPICS"`
}

type DefaultConfig struct {
	// Server default values
	APP_ENV     string
	SERVER_PORT int
	DB_HOST     string
	DB_PORT     int
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	// MQTT default values
	MQTT_BROKER    string
	MQTT_PORT      int
	MQTT_CLIENT_ID string
	MQTT_USERNAME  string
	MQTT_PASSWORD  string
	MQTT_TOPICS    string
}

var Default = DefaultConfig{
	// Server default values
	APP_ENV:     "development",
	SERVER_PORT: 8080,
	DB_HOST:     "localhost",
	DB_PORT:     5432,
	DB_USER:     "worty76",
	DB_PASSWORD: "highlysecurepassword",
	DB_NAME:     "k3s_micro_hs",
	// MQTT default values
	MQTT_BROKER:    "localhost",
	MQTT_PORT:      1883,
	MQTT_CLIENT_ID: "mpa-client",
	MQTT_USERNAME:  "anonymous",
	MQTT_PASSWORD:  "123456",
	MQTT_TOPICS:    "device/+/message,device/+/status",
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".") // Look for the .env file in the root directory
	viper.AutomaticEnv()     // Read in environment variables that match

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Help convert . to _ in env variable names

	// Set default values for the configuration
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		// If the config file is not found, we can ignore the error and continue with environment variables (Maybe on production)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// set default values for server
	viper.SetDefault("APP_ENV", Default.APP_ENV)
	viper.SetDefault("SERVER_PORT", Default.SERVER_PORT)
	viper.SetDefault("DB_HOST", Default.DB_HOST)
	viper.SetDefault("DB_PORT", Default.DB_PORT)
	viper.SetDefault("DB_USER", Default.DB_USER)
	viper.SetDefault("DB_NAME", Default.DB_NAME)
	viper.SetDefault("DB_PASSWORD", Default.DB_PASSWORD)
	// set default values for MQTT
	viper.SetDefault("MQTT_BROKER", Default.MQTT_BROKER)
	viper.SetDefault("MQTT_PORT", Default.MQTT_PORT)
	viper.SetDefault("MQTT_CLIENT_ID", Default.MQTT_CLIENT_ID)
	viper.SetDefault("MQTT_USERNAME", Default.MQTT_USERNAME)
	viper.SetDefault("MQTT_PASSWORD", Default.MQTT_PASSWORD)
	viper.SetDefault("MQTT_TOPICS", Default.MQTT_TOPICS)
}
