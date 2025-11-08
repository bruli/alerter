package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsServerURL   string        `envconfig:"NATS_SERVER_URL" required:"true"`
	TelegramToken   string        `envconfig:"TELEGRAM_TOKEN" required:"true"`
	TelegramChatID  int64         `envconfig:"TELEGRAM_CHAT_ID" required:"true"`
	CacheTTL        time.Duration `envconfig:"CACHE_TTL" required:"true"`
	CacheTTLRefresh time.Duration `envconfig:"CACHE_TTL_REFRESH" required:"true"`
}

func New() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
