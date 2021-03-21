package config

import "github.com/noilpa/umbot/internal/weather/wttrin"

type Config struct {
	Threshold int
	Location  string
	Wttrin    wttrin.Config
	LogLeveL  string
}

func New() Config {
	return Config{
		Threshold: 70,
		Location:  "48.834,2.394", // most precision method - Latitude and longitude in 48.834,2.394 format
		LogLeveL:  "info",
		Wttrin: wttrin.Config{
			Host:   "http://wttr.in",
			Format: "j1",
		},
	}
}
