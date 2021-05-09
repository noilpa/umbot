package app

import (
	"context"
	"net/http"

	"github.com/noilpa/ctxlog"
	"github.com/sirupsen/logrus"

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
		"/help":  a.helpCommand(ctx),
		"/start": a.startCommand(ctx), // start sending everyday notification
		//"/stop":     nil,            // stop sending everyday notifications
		"/location": a.locationCommand(ctx), // set location for current user
		//"/remind":   nil,            // set remind time for current user
		"/weather": a.weatherCommand(ctx), // get current weather
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/check", withLog(ctxlog.From(ctx),
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			header := http.StatusOK
			if _, err := w.Write([]byte("OK")); err != nil {
				ctxlog.From(r.Context()).WithError(err).Error("health check failed")
				header = http.StatusInternalServerError
			}
			w.WriteHeader(header)
		}))

	go func(ctx context.Context) {
		if err := http.ListenAndServe(":8888", mux); err != nil {
			ctxlog.From(ctx).WithError(err).Error("http server failed")
		}
	}(ctx)

	a.tgBot.Bot.Start()
}

func withLog(log *logrus.Entry, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.WithContext(ctxlog.With(r.Context(), log))
		log.Info(r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	}
}
