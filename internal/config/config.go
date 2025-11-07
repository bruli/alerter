package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsServerURL  string `envconfig:"NATS_SERVER_URL" required:"true"`
	TelegramToken  string `envconfig:"TELEGRAM_TOKEN" required:"true"`
	TelegramChatID int64  `envconfig:"TELEGRAM_CHAT_ID" required:"true"`
}

func New() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
