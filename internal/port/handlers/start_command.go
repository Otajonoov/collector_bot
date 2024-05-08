package handlers

import (
	"go.uber.org/zap"
	"gopkg.in/tucnak/telebot.v2"
)

func (h *Handler) Start(m *telebot.Message) {

	msg, err := h.b.Send(m.Sender, "Xush kelibsiz", &telebot.SendOptions{
		ReplyMarkup: AuthButtons,
	})

	h.log.Info("sent data", zap.Any("msg", msg), zap.Error(err))
}
