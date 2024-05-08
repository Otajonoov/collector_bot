package handlers

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

const (
	PasswordPrompt       = "Kirish uchun parolni kiriting: "
	IncorrectPasswordMsg = "Parol noto'g'ri. Iltimos parolni qayta kiriting: "
	CorrectPassword      = "iman2024"
)

func (h *Handler) CheckClient(m *tb.Message) {
	h.b.Send(m.Sender, PasswordPrompt)
	// Register the password handler outside the CheckClient method
	h.b.Handle(tb.OnText, func(msg *tb.Message) {
		if msg.Text == CorrectPassword {
			h.AddInformation(msg)
		} else {
			// If the password is incorrect, notify the user and reset the password
			h.b.Send(msg.Sender, IncorrectPasswordMsg)
		}
	})
}

func (h *Handler) AddInformation(m *tb.Message) {
	log.Println("Adding information")
	h.b.Send(m.Sender, "Siz tizimga kirdingiz", &tb.SendOptions{
		ReplyMarkup: AdminButtons,
	})
}
