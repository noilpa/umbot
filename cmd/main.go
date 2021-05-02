package main

import (
	"context"

	"github.com/noilpa/ctxlog"
	"github.com/sirupsen/logrus"

	"github.com/noilpa/umbot/internal/app"
	"github.com/noilpa/umbot/internal/config"
	"github.com/noilpa/umbot/internal/messenger/telegram"
	"github.com/noilpa/umbot/internal/storage/redis"
	"github.com/noilpa/umbot/internal/weather/wttrin"
)

func main() {
	cfg := config.New()

	log := logrus.New()
	lvl, err := logrus.ParseLevel(cfg.LogLeveL)
	if err != nil {
		panic(err)
	}
	log.SetLevel(lvl)
	ctx := ctxlog.With(context.Background(), logrus.NewEntry(log))

	a := app.New(
		cfg,
		wttrin.New(cfg.Wttrin.Host, cfg.Wttrin.Format),
		telegram.New(cfg.Telegram),
		redis.New(),
	)
	a.Run(ctx)
}
