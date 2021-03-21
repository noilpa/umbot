package app

import (
	"context"
	"fmt"

	"github.com/noilpa/ctxlog"

	"github.com/noilpa/umbot/internal/config"
	"github.com/noilpa/umbot/internal/weather"
)

type App struct {
	provider weather.IWeather
	cfg      config.Config
}

func New(cfg config.Config, weatherCli weather.IWeather) App {
	return App{
		cfg:      cfg,
		provider: weatherCli,
	}
}

func (a *App) Run(ctx context.Context) {
	// todo add retry
	isRainy, err := a.provider.IsRainy(ctx, a.cfg.Location, a.cfg.Threshold)
	if err != nil {
		ctxlog.From(ctx).WithError(err).Error("app can not check if it will rain today")
	}

	if isRainy {
		fmt.Println("Take the umbrella!")
	}
}
