package telegram

import "time"

type Config struct {
	URL            string
	Token          string
	PollingTimeout time.Duration
}
