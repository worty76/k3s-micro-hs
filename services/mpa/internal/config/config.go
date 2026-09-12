package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv   string `mapstructure:"APP_ENV"`
	Database DatabaseConfig
	Server   ServerConfig
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

var Default = struct {
	SERVER_PORT int
	DB_HOST     string
	DB_PORT     int
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
}{
	SERVER_PORT: 8080,
	DB_HOST:     "localhost",
	DB_PORT:     5432,
	DB_USER:     "worty76",
	DB_PASSWORD: "highlysecurepassword",
	DB_NAME:     "k3s_micro_hs",
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
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("SERVER_PORT", Default.SERVER_PORT)
	viper.SetDefault("DB_HOST", Default.DB_HOST)
	viper.SetDefault("DB_PORT", Default.DB_PORT)
	viper.SetDefault("DB_USER", Default.DB_USER)
	viper.SetDefault("DB_NAME", Default.DB_NAME)
	viper.SetDefault("DB_PASSWORD", Default.DB_PASSWORD)
}
