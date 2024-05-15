package bot

import (
	"iman_tg_bot/internal/model"
)

func (h *BotHandler) startHandler(user *model.Client) error {
	err := h.repo.ClientUser().ChangeStep(user.ChatId, model.StartCommandStep)
	if err != nil {
		return err
	}

	h.SendKeyboardButton(user, model.MenuText, AuthButtons)
	return nil
}
