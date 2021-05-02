package telegram

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type Telegram struct {
	Bot *tb.Bot
}

func New(c Config) Telegram {
	b, err := tb.NewBot(tb.Settings{
		URL:   c.URL,
		Token: c.Token,
		Poller: &tb.LongPoller{
			Timeout: c.PollingTimeout,
		},
		Verbose: true,
	})
	if err != nil {
		panic(err)
	}

	return Telegram{
		Bot: b,
	}
}

func (tg *Telegram) InitHandlers(handlers map[string]interface{}) {
	for endpoint, handler := range handlers {
		tg.Bot.Handle(endpoint, handler)
	}
}
