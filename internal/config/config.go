package config

import (
	"time"

	"github.com/noilpa/umbot/internal/messenger/telegram"
	"github.com/noilpa/umbot/internal/storage/redis"
	"github.com/noilpa/umbot/internal/weather/wttrin"
)

type Config struct {
	Threshold int
	LogLeveL  string
	Wttrin    wttrin.Config
	Telegram  telegram.Config
	Storage   redis.Config
}

func New() Config {
	return Config{
		Threshold: 70,
		LogLeveL:  "info",
		Wttrin: wttrin.Config{
			Host:   "http://wttr.in",
			Format: "j1",
		},
		Telegram: telegram.Config{
			URL:            "https://api.telegram.org",
			Token:          "secret_token",
			PollingTimeout: 10 * time.Second,
		},
		Storage: redis.Config{
			Address:  "localhost:6379",
			Login:    "",
			Password: "",
		},
	}
}
