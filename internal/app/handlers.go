package app

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/noilpa/ctxlog"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/noilpa/umbot/internal/storage/entity"
)

func (a App) helpCommand(ctx context.Context) func(m *tb.Message) {
	return func(m *tb.Message) {
		if _, err := a.tgBot.Bot.Send(m.Chat, "Я напомню тебе взять зонтик ⛱"); err != nil {
			ctxlog.From(ctx).WithError(err).Error("send help failed")
		}
	}
}

func (a App) weatherCommand(ctx context.Context) func(m *tb.Message) {
	return func(m *tb.Message) {
		if _, err := a.tgBot.Bot.Send(m.Chat, a.weatherMessage(ctx, m)); err != nil {
			ctxlog.From(ctx).WithError(err).Error("send weather failed")
		}
	}
}

func (a App) weatherMessage(ctx context.Context, m *tb.Message) string {
	log := ctxlog.From(ctx)
	u, err := a.dataStorage.Get(ctx, m.Chat.ID)
	if err != nil {
		log.WithError(err).Error("failed to get user from storage")
		return "create user with /start command."
	}

	if u.Latitude == 0 && u.Longitude == 0 {
		return "empty latitude and longitude. setup them with /location command. example service for calculate latitude and longitude - https://u-karty.ru/opredelenie-koordinat-na-karte-yandex"
	}

	isRainy, err := a.weatherProvider.IsRainy(ctx, fmt.Sprintf("%.2f,%.2f", u.Latitude, u.Longitude), a.cfg.Threshold)
	if err != nil {
		ctxlog.From(ctx).WithError(err).Error("app can not check if it will rain today")
		return "the weather forecast service is unavailable. please try later."
	}

	if isRainy {
		return "take an umbrella with you today."
	}

	return "chance of rain is extremely low."
}

func (a App) startCommand(ctx context.Context) func(m *tb.Message) {
	return func(m *tb.Message) {
		if _, err := a.tgBot.Bot.Send(m.Chat, a.startMessage(ctx, m)); err != nil {
			ctxlog.From(ctx).WithError(err).Error("send start failed")
		}
	}
}

func (a App) startMessage(ctx context.Context, m *tb.Message) string {
	userData := entity.ChatData{
		ID:       m.Chat.ID,
		IsActive: true,
	}
	if err := a.dataStorage.Set(ctx, m.Chat.ID, userData); err != nil {
		ctxlog.From(ctx).WithError(err).Error("failed to create user in storage")
		return "failed to create user. please try later."
	}

	return "user configured successfully. add weather check location with /location command."
}

func (a App) locationCommand(ctx context.Context) func(m *tb.Message) {
	return func(m *tb.Message) {
		if _, err := a.tgBot.Bot.Send(m.Chat, a.locationMessage(ctx, m)); err != nil {
			ctxlog.From(ctx).WithError(err).Error("send location failed")
		}
	}
}

func (a App) locationMessage(ctx context.Context, m *tb.Message) string {
	splitted := strings.Split(m.Payload, ",")
	if len(splitted) != 2 {
		return "wrong request format. example: /location latitude,longitude. latitude and longitude must be floating number with point delimiter."
	}
	for i := range splitted{
		splitted[i] = strings.Trim(splitted[i], " ")
	}
	log := ctxlog.From(ctx)
	lat, err := strconv.ParseFloat(splitted[0], 32)
	if err != nil {
		log.WithError(err).Error("failed to parse latitude")
		return "failed to parse latitude. latitude and longitude must be floating number with point delimiter."
	}
	long, err := strconv.ParseFloat(splitted[1], 32)
	if err != nil {
		log.WithError(err).Error("failed to parse longitude")
		return "failed to parse longitude. latitude and longitude must be floating number with point delimiter."
	}

	userData, err := a.dataStorage.Get(ctx, m.Chat.ID)
	if err != nil {
		log.WithError(err).Error("failed to get user")
		return "failed to get user for update. please try later."
	}

	userData.Latitude = float32(lat)
	userData.Longitude = float32(long)

	if err := a.dataStorage.Set(ctx, m.Chat.ID, userData); err != nil {
		log.WithError(err).Error("failed to update user's location")
		return "failed to update location. please try later."
	}

	return "location configured successfully."
}
