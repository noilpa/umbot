package main

import (
	"context"

	"github.com/noilpa/ctxlog"
	"github.com/sirupsen/logrus"

	"github.com/noilpa/umbot/internal/app"
	"github.com/noilpa/umbot/internal/config"
	"github.com/noilpa/umbot/internal/weather/wttrin"
)

// umbot - umbrella bot for telegram.
// umbot remind for you to take umbrella
// in case of rainy weather.

// api
// --------------------
// set location
// set remind time
// get predict now
// check location exist

func main() {
	cfg := config.New()

	log := logrus.New()
	lvl, err := logrus.ParseLevel(cfg.LogLeveL)
	if err != nil {
		panic(err)
	}
	log.SetLevel(lvl)
	ctx := ctxlog.With(context.Background(), logrus.NewEntry(log))

	wttrin.New(cfg.Wttrin.Host, cfg.Wttrin.Format)

	a := app.New(
		cfg,
		wttrin.New(cfg.Wttrin.Host, cfg.Wttrin.Format),
	)
	a.Run(ctx)
}
