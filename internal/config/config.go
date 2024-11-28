package config

import (
	"fmt"
	configLoader "github.com/patyukin/mbs-pkg/pkg/config"
)

type Config struct {
	MinLogLevel string `yaml:"min_log_level" validate:"oneof=debug info warn error fatal panic"`
	HttpServer  struct {
		Port int `yaml:"port" validate:"required,numeric"`
	} `yaml:"http_server" validate:"required"`
	TelegramToken string `yaml:"telegram_token" validate:"required"`
	RabbitMQUrl   string `yaml:"rabbitmq_url" validate:"required"`
	Kafka         struct {
		Brokers []string `yaml:"brokers" validate:"required"`
		Topics  []string `yaml:"topics" validate:"required"`
	} `yaml:"kafka"`
}

func LoadConfig() (*Config, error) {
	var config Config
	err := configLoader.LoadConfig(&config)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return &config, nil
}
