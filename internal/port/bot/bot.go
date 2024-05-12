package bot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"iman_tg_bot/internal/adapter"
	handler "iman_tg_bot/internal/port/handlers"
	"log/slog"
	"time"
)

type Bot struct {
	b    *tb.Bot
	log  *slog.Logger
	repo adapter.Repo
}

func NewBot(log *slog.Logger, clientUser adapter.Repo) (*Bot, error) {
	settings := tb.Settings{
		Token:  "7104945895:AAErhkkfRLLN28PGeU48HRAqVO_1U_yf-X4",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tb.NewBot(settings)
	if err != nil {
		return nil, err
	}

	return &Bot{
		b:    b,
		log:  log,
		repo: clientUser,
	}, nil
}

func (b *Bot) registerHandlers() {
	b.log.Info("registering handlers")

	handlers := handler.NewHandler(b.b, b.log, b.repo)

	b.b.Handle(handler.StartCommand, handlers.Start)
	b.b.Handle(handler.CheckClientCommand, handlers.CheckClient)
	b.b.Handle(handler.AddClientCommand, handlers.AddClient)
}

func (b *Bot) Start() {
	b.log.Info("starting bot")
	b.registerHandlers()
	b.b.Start()
}
