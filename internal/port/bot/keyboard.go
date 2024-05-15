package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var AuthButtons = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(CheckClientCommand),
	),
)

var MalumotQoshish = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(AddClientCommand),
	),
)

var Finish = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(FinishCommand),
	),
)

const (
	CheckClientCommand = "Tizimga kirish"
	AddClientCommand   = "Malumot qo'shish"
	FinishCommand      = "Malumotni yuborish"
)
