package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errInvalidUrl   = errors.New("invalid url")
	errUnauthorized = errors.New("user is not authorized")
	errUnableToSave = errors.New("unable to save")
)

// msg.Text = "Ссылка невалидная"
// msg.Text = "Ты не авторизован! Чтобы авторизоваться введи команду /start"
// msg.Text = "Не удалось сохранить ссылку. Попробуй снова"
func (b *Bot) handleErrors(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, b.messages.Default)

	switch err {
	case errInvalidUrl:
		msg.Text = b.messages.InvalidURL
		b.bot.Send(msg)
	case errUnauthorized:
		msg.Text = b.messages.Unauthorized
		b.bot.Send(msg)
	case errUnableToSave:
		msg.Text = b.messages.UnableToSave
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)
	}

}
