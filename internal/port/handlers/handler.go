package handlers

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"iman_tg_bot/internal/adapter"
	"log/slog"
)

type HandlerFunc func(*tb.Message)

type Handler struct {
	log  *slog.Logger
	b    *tb.Bot
	repo adapter.Repo
}

func NewHandler(b *tb.Bot, log *slog.Logger, repo adapter.Repo) *Handler {
	return &Handler{
		log:  log,
		b:    b,
		repo: repo,
	}
}

func (h *Handler) EmptyAnswer(q *tb.Query) {
	h.b.Answer(q, &tb.QueryResponse{
		QueryID: q.ID,
	})
}
