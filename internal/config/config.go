package config

import (
	"fmt"
	configLoader "github.com/patyukin/mbs-pkg/pkg/config"
)

type Config struct {
	MinLogLevel   string `yaml:"min_log_level" validate:"oneof=debug info warn error fatal panic"`
	TelegramToken string `yaml:"telegram_token" validate:"required"`
	RabbitMQUrl   string `yaml:"rabbitmq_url" validate:"required"`
}

func LoadConfig() (*Config, error) {
	var config Config
	err := configLoader.LoadConfig(&config)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return &config, nil
}
