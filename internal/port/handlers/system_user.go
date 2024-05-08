package handlers

import tb "gopkg.in/tucnak/telebot.v2"

func (h *Handler) SignUp(m *tb.Message) {
	h.b.Send(m.Sender, "Sign Up")
}

func (h *Handler) SignIn(m *tb.Message) {
	h.b.Send(m.Sender, "Sign In")
}
