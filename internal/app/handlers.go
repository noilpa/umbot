package app

import (
	"context"

	"github.com/noilpa/ctxlog"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (a App) help(ctx context.Context) func(m *tb.Message) {
	return func(m *tb.Message) {
		if _, err := a.tgBot.Bot.Send(m.Sender, "Я напомню тебе взять зонтик ⛱"); err != nil {
			ctxlog.From(ctx).WithError(err).Error("send help failed")
		}
	}
}

func (a App) weather(ctx context.Context) func(m *tb.Message) {
	return func(m *tb.Message) {
		isRainy, err := a.weatherProvider.IsRainy(ctx, a.cfg.Location, a.cfg.Threshold)
		if err != nil {
			ctxlog.From(ctx).WithError(err).Error("app can not check if it will rain today")
		}

		var message string
		if isRainy {
			message = "не замочи штанишки сынок"
		} else {
			message = "сегодня все сухо"
		}

		if _, err := a.tgBot.Bot.Send(m.Sender, message); err != nil {
			ctxlog.From(ctx).WithError(err).Error("send weather failed")
		}
	}
}
