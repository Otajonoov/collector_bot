package bot

import (
	"fmt"
	"iman_tg_bot/internal/adapter"
	"iman_tg_bot/internal/model"
	"iman_tg_bot/internal/pkg/config"
	"log"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	bot  *tgbotapi.BotAPI
	log  *slog.Logger
	repo adapter.Repo
	cfg  config.Config
}

func NewBot(cfg config.Config, repo adapter.Repo, log *slog.Logger) BotHandler {
	bot, err := tgbotapi.NewBotAPI("7104945895:AAErhkkfRLLN28PGeU48HRAqVO_1U_yf-X4")
	if err != nil {
		log.Error("", err)
	}
	bot.Debug = true

	return BotHandler{
		cfg:  cfg,
		log:  log,
		bot:  bot,
		repo: repo,
	}
}

func (h *BotHandler) Start() {
	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)
	for update := range updates {
		go h.HandleBot(update)
	}
}

func (h *BotHandler) HandleBot(update tgbotapi.Update) {
	if update.Message == nil {
		h.log.Error("Error", "update.Message is nil")
		return
	}
	user, err := h.repo.ClientUser().GetOrCreate(update.Message.Chat.ID, update.Message.From.FirstName+"_"+update.Message.From.LastName)
	if err != nil {
		h.SendTextMessage(user, "xatolik yuz berdi")
	}
	if update.Message.Command() == "start" {
		err = h.startHandler(user)
		if err != nil {
			log.Println(err)
		}
	} else if update.Message.Text != "" {
		msg := update.Message.Text
		switch user.Step {
		case model.StartCommandStep:
			switch msg {
			case CheckClientCommand:
				h.SendTextMessage(user, "Tizimga kirish uchun parolni kiriting")
				// Set the next step after receiving the password
				user.Step = model.CheckUserPassword
				err := h.repo.ClientUser().ChangeStep(user.ChatId, model.CheckUserPassword)
				if err != nil {
					log.Println(err)
				}
			}
		case model.CheckUserPassword:
			// Verify password
			if msg == "iman2024" {
				h.SendKeyboardButton(user, "Siz tizimga kirdingiz", MalumotQoshish)

				user.Step = model.AddData
				err := h.repo.ClientUser().ChangeStep(user.ChatId, model.AddData)
				if err != nil {
					log.Println(err)
				}
			} else {
				h.SendTextMessage(user, "Noto'g'ri parol, qaytadan kiriting")
			}
		case model.AddData:
			//h.SendTextMessage(user, "Contract ID")
			h.RemoveKeyboard(user, "Contract ID")

			user.Step = model.ContractId
			err = h.repo.ClientUser().ChangeStep(user.ChatId, model.ContractId)
			if err != nil {
				log.Println(err)
			}
		case model.ContractId:
			h.SendTextMessage(user, "Phone Number")
			err := h.repo.ClientUser().UpdateOneFild(user.ChatId, "contract_id", msg)
			if err != nil {
				log.Println(err)
			}

			user.Step = model.PhoneNumber
			err = h.repo.ClientUser().ChangeStep(user.ChatId, model.PhoneNumber)
			if err != nil {
				log.Println(err)
			}
		case model.PhoneNumber:
			h.SendTextMessage(user, "Address")
			err := h.repo.ClientUser().UpdateOneFild(user.ChatId, "phone_number", msg)
			if err != nil {
				log.Println(err)
			}

			user.Step = model.Address
			err = h.repo.ClientUser().ChangeStep(user.ChatId, model.Address)
			if err != nil {
				log.Println(err)
			}
		case model.Address:
			h.SendTextMessage(user, "Payment sum")
			err := h.repo.ClientUser().UpdateOneFild(user.ChatId, "address", msg)
			if err != nil {
				log.Println(err)
			}

			user.Step = model.PaymentSum
			err = h.repo.ClientUser().ChangeStep(user.ChatId, model.PaymentSum)
			if err != nil {
				log.Println(err)
			}
		case model.PaymentSum:
			h.SendTextMessage(user, "Comment")
			err := h.repo.ClientUser().UpdateOneFild(user.ChatId, "payment_sum", msg)
			if err != nil {
				log.Println(err)
			}

			user.Step = model.Comment
			err = h.repo.ClientUser().ChangeStep(user.ChatId, model.Comment)
			if err != nil {
				log.Println(err)
			}
		case model.Comment:
			h.SendTextMessage(user, "Location")
			err := h.repo.ClientUser().UpdateOneFild(user.ChatId, "comment", msg)
			if err != nil {
				log.Println(err)
			}

			user.Step = model.Location
			err = h.repo.ClientUser().ChangeStep(user.ChatId, model.Location)
			if err != nil {
				log.Println(err)
			}
		case model.Finish:
			if msg == FinishCommand {
				h.finalizeAddClient(user.ChatId)
				h.SendKeyboardButton(user, "Qayta malumot qo'shishingiz mumkin", MalumotQoshish)
				err := h.repo.ClientUser().ChangeStep(user.ChatId, model.AddData)
				if err != nil {
					log.Println(err)
				}
			}

		}
	} else if update.Message.Location != nil {
		loc := update.Message.Location
		switch user.Step {
		case model.Location:
			h.SendTextMessage(user, "Address Foto")
			location := fmt.Sprintf("%f,%f", loc.Latitude, loc.Longitude)
			err := h.repo.ClientUser().UpdateOneFild(user.ChatId, "location", location)
			if err != nil {
				log.Println(err)
			}

			user.Step = model.AddressFotoPath
			err = h.repo.ClientUser().ChangeStep(user.ChatId, model.AddressFotoPath)
			if err != nil {
				log.Println(err)
			}
		}
	} else if update.Message.Photo != nil {
		switch user.Step {
		case model.AddressFotoPath:
			// Check if the update has a photo
			if len(update.Message.Photo) == 0 {
				h.SendTextMessage(user, "Rasimni qaytadan yuboring")
				return
			}
			h.SendTextMessage(user, "Payment Photo")

			// Get the FileID of the last photo in the array
			fileId := update.Message.Photo[len(update.Message.Photo)-1].FileID

			// Save the photo to the images folder
			imagePath, err := h.savePhotoToFolder(h.bot, fileId)
			if err != nil {
				log.Println("Error saving address photo:", err)
				return
			}
			// Update the database with the address photo path
			err = h.repo.ClientUser().UpdateOneFild(user.ChatId, "address_foto_path", imagePath)
			if err != nil {
				log.Println("Error updating address photo path:", err)
				return
			}
			// Move to the next step
			user.Step = model.PaymentFotoPath
			err = h.repo.ClientUser().ChangeStep(user.ChatId, model.PaymentFotoPath)
			if err != nil {
				log.Println("Error changing step to PaymentFotoPath:", err)
				return
			}
		case model.PaymentFotoPath:
			h.SendKeyboardButton(user, "Malumot qo'shish yakunlandi.", Finish)

			// Check if the update has a photo
			if len(update.Message.Photo) == 0 {
				h.SendTextMessage(user, "Rasimni qaytadan yuboring")
				return
			}
			// Get the FileID of the last photo in the array
			fileId := update.Message.Photo[len(update.Message.Photo)-1].FileID
			log.Println("FileID:", fileId)
			// Save the photo to the images folder
			imagePath, err := h.savePhotoToFolder(h.bot, fileId)
			if err != nil {
				h.SendTextMessage(user, "Rasimni qaytadan yuboring")
			}

			err = h.repo.ClientUser().UpdateOneFild(user.ChatId, "payment_foto_path", imagePath)
			if err != nil {
				h.SendTextMessage(user, "Rasimni qaytadan yuboring")
			}

			user.Step = model.Finish
			err = h.repo.ClientUser().ChangeStep(user.ChatId, model.Finish)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (h *BotHandler) SendTextMessage(user *model.Client, text string) {
	msg := tgbotapi.NewMessage(user.ChatId, text)
	msg.ParseMode = "html"
	if _, err := h.bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func (h *BotHandler) SendKeyboardButton(user *model.Client, text string, keywd tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(user.ChatId, text)
	msg.ParseMode = "html"
	msg.ReplyMarkup = keywd
	if _, err := h.bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func (h *BotHandler) RemoveKeyboard(user *model.Client, text string) {
	msg := tgbotapi.NewMessage(user.ChatId, text)
	msg.ParseMode = "html"
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	if _, err := h.bot.Send(msg); err != nil {
		log.Println(err)
	}
}
