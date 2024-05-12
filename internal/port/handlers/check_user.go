package handlers

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	PasswordPrompt       = "Kirish uchun parolni kiriting"
	IncorrectPasswordMsg = "Parol noto'g'ri."
	CorrectPassword      = "iman2024"
)

func (h *Handler) CheckClient(m *tb.Message) {
	h.b.Send(m.Sender, PasswordPrompt)

	h.b.Handle(tb.OnText, func(msg *tb.Message) {
		if msg.Text == CorrectPassword {
			h.b.Send(m.Sender, "Siz tizimga kirdingiz", &tb.SendOptions{
				ReplyMarkup: AdminButtons,
			})
		} else {
			h.b.Send(msg.Sender, IncorrectPasswordMsg)
			h.CheckClient(msg)
		}
	})
}
