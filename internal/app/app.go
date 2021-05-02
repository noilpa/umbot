package app

import (
	"context"

	"github.com/noilpa/umbot/internal/config"
	"github.com/noilpa/umbot/internal/messenger/telegram"
	"github.com/noilpa/umbot/internal/storage"
	"github.com/noilpa/umbot/internal/weather"
)

type App struct {
	cfg             config.Config
	weatherProvider weather.IWeather
	tgBot           telegram.Telegram
	dataStorage     storage.Storage
}

func New(cfg config.Config, weatherCli weather.IWeather, tg telegram.Telegram, s storage.Storage) App {
	return App{
		cfg:             cfg,
		weatherProvider: weatherCli,
		tgBot:           tg,
		dataStorage:     s,
	}
}

func (a *App) Run(ctx context.Context) {
	a.tgBot.InitHandlers(map[string]interface{}{
		"/help": a.help(ctx),
		//"/start":    nil,            // start sending everyday notification
		//"/stop":     nil,            // stop sending everyday notifications
		//"/location": nil,            // set location for current user
		//"/remind":   nil,            // set remind time for current user
		"/weather": a.weather(ctx), // get current weather
	})

	a.tgBot.Bot.Start()
}
