package handler

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/realtemirov/encryption/repo"
	"github.com/realtemirov/encryption/service"
)

type Handler interface {
	Messages(m *tg.Message)
}

type handler struct {
	serv service.Service
	bot  *tg.BotAPI
	db   repo.Db
}

func (h *handler) Messages(m *tg.Message) {
	u, err := h.db.GetUserByTelegramID(m.Chat.ID)
	if err != nil {
		if err == repo.ErrUserNotFound {
			u = repo.User{
				TelegramID: m.Chat.ID,
				FullName:   m.Chat.FirstName + " " + m.Chat.LastName,
			}
			if err := h.db.AddUser(u); err != nil {
				return
			}
		} else {
			return
		}
	}

	msg := tg.NewMessage(m.Chat.ID, "Hello, "+u.FullName+"! I'm a bot that can encrypt and decrypt text. To encrypt text, send me a message with the text you want to encrypt. To decrypt text, send me a message with the text you want to decrypt.")
	msg.ReplyMarkup = tg.NewReplyKeyboard(
		tg.NewKeyboardButtonRow(
			tg.NewKeyboardButton("AES"),
			tg.NewKeyboardButton("DES"),
		))

	switch m.Text {
	case string(repo.AES), string(repo.DES):
		u.Type = repo.Type(m.Text)
		if err := h.db.UpdateUser(u); err != nil {
			return
		}

		msg = tg.NewMessage(m.Chat.ID, "You have chosen the encryption method: "+string(u.Type))
		msg.ReplyMarkup = tg.NewReplyKeyboard(
			tg.NewKeyboardButtonRow(
				tg.NewKeyboardButton("Encrypt"),
				tg.NewKeyboardButton("Decrypt"),
			))
	case "Encrypt", "Decrypt":

		if m.Text == "Encrypt" {
			u.Encryption = true
			if err := h.db.UpdateUser(u); err != nil {
				return
			}
			msg = tg.NewMessage(m.Chat.ID, "Enter the text you want to encrypt.")
		}
		if m.Text == "Decrypt" {
			u.Decryption = true
			if err := h.db.UpdateUser(u); err != nil {
				return
			}
			msg = tg.NewMessage(m.Chat.ID, "Enter the text you want to decrypt.")
		}
		msg.ReplyMarkup = tg.NewRemoveKeyboard(true)

	default:
		var (
			text = msg.Text
			err  error
		)

		if u.Encryption {
			text, err = h.serv.Encryption(&u, m.Text)
			if err != nil {
				return
			}
			u.Encryption = false
		}

		if u.Decryption {
			text, err = h.serv.Decryption(&u, m.Text)
			if err != nil {
				return
			}
			u.Decryption = false
		}

		msg = tg.NewMessage(m.Chat.ID, text)
		msg.ReplyMarkup = tg.NewReplyKeyboard(
			tg.NewKeyboardButtonRow(
				tg.NewKeyboardButton("AES"),
				tg.NewKeyboardButton("DES"),
			))

		u.Type = ""
		if err := h.db.UpdateUser(u); err != nil {
			return
		}
	}

	msg.ParseMode = tg.ModeMarkdown
	if _, err = h.bot.Send(msg); err != nil {
		return
	}

}

func NewHandler(db repo.Db, serv service.Service, bot *tg.BotAPI) Handler {
	return &handler{
		serv: serv,
		bot:  bot,
		db:   db,
	}
}
