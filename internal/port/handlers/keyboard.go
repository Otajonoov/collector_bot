package handlers

import (
	"gopkg.in/tucnak/telebot.v2"
)

var AuthButtons = &telebot.ReplyMarkup{
	ReplyKeyboard: [][]telebot.ReplyButton{
		{
			telebot.ReplyButton{Text: "Tizimga kirish"},
		},
	},
}

var AdminButtons = &telebot.ReplyMarkup{
	ReplyKeyboard: [][]telebot.ReplyButton{
		{
			telebot.ReplyButton{Text: "Malumot qo'shish"},
		},
	},
}
